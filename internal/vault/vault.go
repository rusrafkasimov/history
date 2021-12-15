package vault

import (
	"errors"
	"fmt"
	"github.com/hashicorp/vault/api"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Provider struct {
	path    string
	client  *api.Logical
	results map[string]map[string]string
}

// NewVaultProvider initialize new vault provider
func NewVaultProvider() *Provider {
	content, err := ioutil.ReadFile("/run/secrets/vault_dev_root_token_id")
	if err != nil {
		content = []byte(os.Getenv("VAULT_TOKEN"))
	}

	vaultPath := os.Getenv("VAULT_PATH")
	vaultToken := string(content)
	vaultAddress := os.Getenv("VAULT_ADDRESS")

	provider, err := New(vaultToken, vaultAddress, vaultPath)
	if err != nil {
		log.Fatalln("Couldn't load provider", err)
	}

	return provider
}

// New prepare and create new vault provider
func New(token, addr, path string) (*Provider, error) {
	config := &api.Config{
		Address: addr,
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	client.SetToken(token)

	return &Provider{
		path:    path,
		client:  client.Logical(),
		results: make(map[string]map[string]string),
	}, nil
}

// Get retrieves a value from vault using the KV engine. The actual key selected is determined by the value
// separated by the colon. For example "database:password" will retrieve the key "password" from the path
// "database".
func (p *Provider) Get(v string) (string, error) {
	// <path>/data/<path-secret>:key
	split := strings.Split(v, ":")
	if len(split) == 1 {
		return "", errors.New("missing key value")
	}

	pathSecret := split[0]
	key := split[1]

	res, ok := p.results[pathSecret]
	if ok {
		val, ok := res[key]
		if !ok {
			return "", errors.New("key not found in cached data")
		}

		return val, nil
	}

	secret, err := p.client.Read(fmt.Sprintf("%s/data/%s", p.path, pathSecret))
	if err != nil {
		return "", fmt.Errorf("reading: %w", err)
	}

	if secret == nil {
		return "", errors.New("secret not found")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", errors.New("invalid data in secret")
	}

	secrets := make(map[string]string)

	for k, v := range data {
		val, ok := v.(string)
		if !ok {
			return "", errors.New("secret value in data is not string")
		}

		secrets[k] = val
	}

	val, ok := secrets[key]
	if !ok {
		return "", errors.New("key not found in retrieved data")
	}

	p.results[pathSecret] = secrets

	return val, nil
}


