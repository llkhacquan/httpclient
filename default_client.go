package httpclient

import "context"

var defaultClient = &Client{}

// Get performs a GET request using the default client
func Get(ctx context.Context, url string, result interface{}, opts ...Option) error {
	return defaultClient.Get(ctx, url, result, opts...)
}

// Post performs a POST request using the default client
func Post(ctx context.Context, url string, body interface{}, result interface{}, opts ...Option) error {
	return defaultClient.Post(ctx, url, body, result, opts...)
}

// Patch performs a PATCH request using the default client
func Patch(ctx context.Context, url string, body interface{}, result interface{}, opts ...Option) error {
	return defaultClient.Patch(ctx, url, body, result, opts...)
}

// Put performs a PUT request using the default client
func Put(ctx context.Context, url string, body interface{}, result interface{}, opts ...Option) error {
	return defaultClient.Put(ctx, url, body, result, opts...)
}

// Delete performs a DELETE request using the default client
func Delete(ctx context.Context, url string, result interface{}, opts ...Option) error {
	return defaultClient.Delete(ctx, url, result, opts...)
}
