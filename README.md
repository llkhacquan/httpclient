# HTTP Client

A simple and flexible HTTP client for Go that automatically handles JSON marshaling/unmarshaling with support for custom options.

## Features

- Simple HTTP methods: GET, POST, PUT, PATCH, DELETE
- Automatic JSON marshaling and unmarshaling
- Customizable marshal/unmarshal functions
- Request options (headers, status capture)
- Built on Go's standard `net/http` package

## Installation

```bash
go get github.com/llkhacquan/httpclient
```

## API Reference

### Client

```go
type Client struct {
    Client        *http.Client                        // HTTP client (defaults to http.DefaultClient)
    MarshalFunc   func(v any) ([]byte, error)        // JSON marshal function (defaults to json.Marshal)
    UnmarshalFunc func(data []byte, v any) error     // JSON unmarshal function (defaults to json.Unmarshal)
}
```

### Methods

- `Get(ctx context.Context, url string, result interface{}, opts ...Option) error`
- `Post(ctx context.Context, url string, body interface{}, result interface{}, opts ...Option) error`
- `Put(ctx context.Context, url string, body interface{}, result interface{}, opts ...Option) error`
- `Patch(ctx context.Context, url string, body interface{}, result interface{}, opts ...Option) error`
- `Delete(ctx context.Context, url string, result interface{}, opts ...Option) error`

### Options

- `WithHeader(key, value string) Option` - Add a custom header
- `WithHeaders(headers map[string]string) Option` - Add multiple headers
- `WithStatus(status *int) Option` - Capture HTTP status code and allow non-200 responses

## Testing

Run the tests:

```bash
go test -v
```

The tests use real APIs (Pokemon API and httpbin.org) to ensure the client works correctly with actual HTTP services.

## Credits

This package was built with assistance from [Claude Code](https://claude.ai/code) by Anthropic.

## License

MIT License
