package llm

import "github.com/dsreek/sreekdAgent/pkg/model"

type Config struct {
	Provider string
	Model    string
	Enabled  bool
}

type Client struct{}

func New(cfg Config) *Client { return &Client{} }

func (c *Client) Enrich(issues []model.Issue) []model.Issue {
	return issues
}
