package delta

// config customize delta behaviour
type config struct {
	transformer    Transformer
	encodeResponse bool
}

type Options func(*config)

func WithEncodeResponse() Options {
	return func(c *config) {
		c.encodeResponse = true
	}
}
