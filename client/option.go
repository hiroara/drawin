package client

type Option func(cli *Client)

func WithRetryConfig(cfg *RetryConfig) Option {
	return func(cli *Client) {
		cli.retry = cfg
	}
}
