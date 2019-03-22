package fakesentry

type Option func(*Handler)

type Logger interface {
	Printf(format string, arg ...interface{})
}

func WithLogger(logger Logger) Option {
	return func(h *Handler) {
		if logger != nil {
			h.logger = logger
		}
	}
}
