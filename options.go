package httpclient

import "net/http"

// Options contains configuration for HTTP requests
type Options struct {
	// Headers to add to the request
	Headers map[string]string
	// Status allows non-200 status codes without returning an error
	Status *int
	// Custom HTTP client for this request only
	Client *http.Client
}

// Option is a function that modifies Options
type Option func(*Options)

// WithHeaders sets custom headers for the request
func WithHeaders(headers map[string]string) Option {
	return func(o *Options) {
		o.Headers = headers
	}
}

// WithHeader sets a single header for the request
func WithHeader(key, value string) Option {
	return func(o *Options) {
		if o.Headers == nil {
			o.Headers = make(map[string]string)
		}
		o.Headers[key] = value
	}
}

// WithStatus allows non-200 status codes without returning an error, the response status code will be stored in the provided pointer
func WithStatus(status *int) Option {
	return func(o *Options) {
		o.Status = status
	}
}

// buildOptions creates Options from Option functions
func buildOptions(opts ...Option) *Options {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
