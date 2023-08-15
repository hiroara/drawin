package client

type Option func(cli *Client)

func WithHandlers(hs ...Handler) Option {
	return func(cli *Client) {
		cli.handlers = hs
	}
}

func WithRetryConfig(cfg *RetryConfig) Option {
	return func(cli *Client) {
		cli.retry = cfg
	}
}
