package logger

import (
	"fmt"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/rusrafkasimov/history/internal/config"

	"time"
)

func NewLogger(sourceName, jobName string, configuration *config.Configuration) (client promtail.Client, err error) {
	labels := "{source=\"" + sourceName + "\",job=\"" + jobName + "\"}"

	host, err := configuration.Get("LOKI_AGENT_HOST")
	if err != nil {
		return nil, err
	}

	port, err := configuration.Get("LOKI_AGENT_PORT")
	if err != nil {
		return nil, err
	}

	conf := promtail.ClientConfig{
		PushURL:            fmt.Sprintf("http://%v:%v/api/prom/push", host, port),
		Labels:             labels,
		BatchWait:          5 * time.Second,
		BatchEntriesNumber: 10000,
		SendLevel:          promtail.INFO,
		PrintLevel:         promtail.DEBUG,
	}

	loki, err := promtail.NewClientJson(conf)

	return loki, err
}
