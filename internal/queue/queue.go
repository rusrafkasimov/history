package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/go-kit/kit/sd/lb"
	"github.com/nats-io/stan.go"
	"github.com/rusrafkasimov/history/internal/config"
	"github.com/rusrafkasimov/history/internal/trace"
	"github.com/rusrafkasimov/history/pkg/dto"
	"sync"
	"time"
)

var (
	errNoConnection = errors.New("no connection to NATS system")
	errQueueClosed  = errors.New("queue already closed")
)

const historyQueueName = "history_queue"

type Event interface {
	Operation() dto.AccountHistory
	Sequence() uint64
	Ack() error
}

type event struct {
	hist dto.AccountHistory
	seq  uint64
	ack  func() error
}

func (e *event) Operation() dto.AccountHistory {
	return e.hist
}

func (e *event) Sequence() uint64 {
	return e.seq
}

func (e *event) Ack() error {
	return e.ack()
}

type Queue struct {
	writeOnly        bool
	logger           promtail.Client
	nodeID           string
	url              string
	ackWait          time.Duration
	reconnectTimeout time.Duration
	clusterID        string
	subject          string

	mu             sync.RWMutex
	conn           stan.Conn
	input          chan Event
	output         chan Event
	sequenceNumber uint64
	closed         bool

	wg     sync.WaitGroup
	doneCh chan struct{}

	now func() time.Time
}

func NewQueue(ctx context.Context, logger promtail.Client, config *config.Configuration) (*Queue, error) {

	URL, err := config.Get("EVENT_QUEUE_URL")
	if err != nil {
		fmt.Println(err)
	}

	ClusterID, err := config.Get("EVENT_QUEUE_CLUSTER_ID")
	if err != nil {
		fmt.Println(err)
	}

	Subject, err := config.Get("EVENT_QUEUE_SUBJECT")
	if err != nil {
		fmt.Println(err)
	}

	clientId := ctx.Value("Name")

	q := &Queue{
		writeOnly:        false,
		logger:           logger,
		nodeID:           clientId.(string),
		url:              URL,
		ackWait:          time.Second,
		reconnectTimeout: time.Second,
		clusterID:        ClusterID,
		subject:          Subject,
		mu:               sync.RWMutex{},
		conn:             nil,
		input:            make(chan Event),
		output:           make(chan Event),
		sequenceNumber:   0,
		closed:           false,
		wg:               sync.WaitGroup{},
		doneCh:           make(chan struct{}),
		now:              time.Now,
	}

	err = q.connect()
	if err != nil {
		q.logger.Errorf("error first connection to NATS. trying dial background...")
		trace.OnError(q.logger,nil, err)
		q.dialBackground()
	}

	if !q.writeOnly {
		q.handleEvents()
	}

	return q, nil
}

func (q *Queue) dialBackground() {
	q.wg.Add(1)
	go func() {
		ticker := time.NewTicker(q.reconnectTimeout)
		defer func() {
			ticker.Stop()
			q.wg.Done()
		}()

		for {
			select {
			case <-ticker.C:
				err := q.connect()
				if err != nil {
					q.logger.Errorf("Error background connection NATS_Streaming")
					trace.OnError(q.logger,nil, err)
					continue
				}
				return
			case <-q.doneCh:
				return
			}
		}
	}()
}

func (q *Queue) connect() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	stanOpts := []stan.Option{
		stan.SetConnectionLostHandler(q.recoverConn),
		stan.NatsURL(q.url),
	}

	conn, err := stan.Connect(q.clusterID, q.nodeID, stanOpts...)
	if err != nil {
		return err
	}

	q.logger.Infof("Established NATS_Streaming connection")

	if !q.writeOnly {
		subsOpts := []stan.SubscriptionOption{
			stan.MaxInflight(1),
			stan.AckWait(q.ackWait),
			stan.SetManualAckMode(),
		}

		if q.sequenceNumber > 0 {
			subsOpts = append(subsOpts, stan.StartAtSequence(q.sequenceNumber))
		} else {
			q.logger.Infof("No deliver all available")
			//subsOpts = append(subsOpts, stan.DeliverAllAvailable())
		}

		_, err = conn.QueueSubscribe(q.subject, historyQueueName, q.handleMessage, subsOpts...)
		if err != nil {
			_ = conn.Close()
			return fmt.Errorf("failed to subscribe to queue: %w", err)
		}
	}

	q.conn = conn
	return nil
}

func (q *Queue) recoverConn(_ stan.Conn, err error) {
	q.logger.Infof("NATS_Streaming connection lost. Trying dial background...")
	q.dialBackground()
}

func (q *Queue) getConn() (stan.Conn, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.conn == nil {
		return nil, errNoConnection
	}

	return q.conn, nil
}

func (q *Queue) Subscribe() (<-chan Event, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return nil, errQueueClosed
	}

	return q.output, nil
}

func (q *Queue) handleMessage(msg *stan.Msg) {
	history := &dto.AccountHistory{}
	err := json.Unmarshal(msg.Data, history)
	if err != nil {
		q.logger.Errorf("Can't unmarshall message")
		trace.OnError(q.logger,nil, err)
		return
	}

	q.mu.Lock()
	q.sequenceNumber = msg.Sequence
	q.mu.Unlock()

	hist := &event{
		hist: *history,
		seq:  msg.Sequence,
		ack:  msg.Ack,
	}

	select {
	case q.input <- hist:
	case <-q.doneCh:
	}
}

func (q *Queue) Publish (ah dto.AccountHistory) error {
	conn, err := q.getConn()
	if err != nil {
		return err
	}

	data, err := json.Marshal(ah)
	if err != nil {
		return fmt.Errorf("can't encode data: %w", err)
	}

	return conn.Publish(q.subject, data)
}

func (q *Queue) handleEvents() {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		for {
			var evt Event

			select {
			case evt = <-q.input:
			case <-q.doneCh:
				return
			}

			select {
			case q.output <- evt:
			case <-q.doneCh:
				return
			}
		}
	}()
}

func (q *Queue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return nil
	}

	close(q.doneCh)
	q.wg.Wait()

	close(q.output)

	finalErr := lb.RetryError{}

	if q.conn != nil {
		err := q.conn.Close()
		if err != nil {
			err = fmt.Errorf("failed to close connection: %w", err)
			finalErr.RawErrors = append(finalErr.RawErrors, err)
		}
	}

	q.closed = true

	if len(finalErr.RawErrors) > 0 {
		err := errors.New("unable to close queue")
		finalErr.RawErrors = append(finalErr.RawErrors, err)
		finalErr.Final = err
		return finalErr
	}
	return nil
}