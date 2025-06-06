package httpclient

var defaultClient = &Client{}

// Get performs a GET request using the default client
func Get(url string, result interface{}, opts ...Option) error {
	return defaultClient.Get(url, result, opts...)
}

// Post performs a POST request using the default client
func Post(url string, body interface{}, result interface{}, opts ...Option) error {
	return defaultClient.Post(url, body, result, opts...)
}

// Patch performs a PATCH request using the default client
func Patch(url string, body interface{}, result interface{}, opts ...Option) error {
	return defaultClient.Patch(url, body, result, opts...)
}

// Delete performs a DELETE request using the default client
func Delete(url string, result interface{}, opts ...Option) error {
	return defaultClient.Delete(url, result, opts...)
}
