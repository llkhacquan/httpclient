package httpclient

import (
	"context"
	"net/http"
	"testing"
)

func TestDefaultClientFunctions(t *testing.T) {
	t.Run("package Get function", func(t *testing.T) {
		var result map[string]interface{}
		err := Get(context.Background(), "https://httpbin.org/get", &result)
		if err != nil {
			t.Fatalf("package Get function failed: %v", err)
		}

		if result["url"] != "https://httpbin.org/get" {
			t.Errorf("unexpected response: %v", result)
		}
	})

	t.Run("package Post function", func(t *testing.T) {
		postData := map[string]string{"test": "value"}
		var result map[string]interface{}

		err := Post(context.Background(), "https://httpbin.org/post", postData, &result)
		if err != nil {
			t.Fatalf("package Post function failed: %v", err)
		}

		jsonData := result["json"].(map[string]interface{})
		if jsonData["test"] != "value" {
			t.Errorf("expected test value 'value', got '%v'", jsonData["test"])
		}
	})

	t.Run("package Patch function", func(t *testing.T) {
		patchData := map[string]string{"update": "patch"}
		var result map[string]interface{}

		err := Patch(context.Background(), "https://httpbin.org/patch", patchData, &result)
		if err != nil {
			t.Fatalf("package Patch function failed: %v", err)
		}

		jsonData := result["json"].(map[string]interface{})
		if jsonData["update"] != "patch" {
			t.Errorf("expected update value 'patch', got '%v'", jsonData["update"])
		}
	})

	t.Run("package Delete function", func(t *testing.T) {
		var result map[string]interface{}
		err := Delete(context.Background(), "https://httpbin.org/delete", &result)
		if err != nil {
			t.Fatalf("package Delete function failed: %v", err)
		}

		if result["url"] != "https://httpbin.org/delete" {
			t.Errorf("unexpected response: %v", result)
		}
	})

	t.Run("package functions with options", func(t *testing.T) {
		var status int
		var result map[string]interface{}

		err := Get(context.Background(), "https://httpbin.org/get", &result,
			WithHeader("X-Test", "package-function"),
			WithStatus(&status))
		if err != nil {
			t.Fatalf("package Get with options failed: %v", err)
		}

		if status != http.StatusOK {
			t.Errorf("expected status 200, got %d", status)
		}

		headers := result["headers"].(map[string]interface{})
		if headers["X-Test"] != "package-function" {
			t.Errorf("expected X-Test header 'package-function', got '%v'", headers["X-Test"])
		}
	})
}
