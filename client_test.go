package httpclient

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// Pokemon represents a Pokemon from the PokeAPI
type Pokemon struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Types  []struct {
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

// PokemonList represents a paginated list of Pokemon
type PokemonList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func TestClient_Get(t *testing.T) {
	client := &Client{}

	t.Run("get single pokemon", func(t *testing.T) {
		var pokemon Pokemon
		err := client.Get(context.Background(), "https://pokeapi.co/api/v2/pokemon/pikachu", &pokemon)
		if err != nil {
			t.Fatalf("GET request failed: %v", err)
		}

		if pokemon.Name != "pikachu" {
			t.Errorf("expected name 'pikachu', got '%s'", pokemon.Name)
		}
		if pokemon.ID != 25 {
			t.Errorf("expected ID 25, got %d", pokemon.ID)
		}
		if len(pokemon.Types) == 0 {
			t.Error("expected pokemon to have types")
		}
	})

	t.Run("get pokemon list", func(t *testing.T) {
		var list PokemonList
		err := client.Get(context.Background(), "https://pokeapi.co/api/v2/pokemon?limit=5", &list)
		if err != nil {
			t.Fatalf("GET request failed: %v", err)
		}

		if list.Count == 0 {
			t.Error("expected count to be greater than 0")
		}
		if len(list.Results) != 5 {
			t.Errorf("expected 5 results, got %d", len(list.Results))
		}
	})

	t.Run("get with custom headers", func(t *testing.T) {
		var pokemon Pokemon
		err := client.Get(context.Background(), "https://pokeapi.co/api/v2/pokemon/1", &pokemon,
			WithHeader("User-Agent", "httpclient-test/1.0"))
		if err != nil {
			t.Fatalf("GET request failed: %v", err)
		}

		if pokemon.Name == "" {
			t.Error("expected pokemon name to be set")
		}
	})

	t.Run("get with status capture", func(t *testing.T) {
		var status int
		var pokemon Pokemon
		err := client.Get(context.Background(), "https://pokeapi.co/api/v2/pokemon/bulbasaur", &pokemon,
			WithStatus(&status))
		if err != nil {
			t.Fatalf("GET request failed: %v", err)
		}

		if status != http.StatusOK {
			t.Errorf("expected status 200, got %d", status)
		}
		if pokemon.Name != "bulbasaur" {
			t.Errorf("expected name 'bulbasaur', got '%s'", pokemon.Name)
		}
	})

	t.Run("get non-existent pokemon with status capture", func(t *testing.T) {
		var status int
		var result interface{} // Use interface{} since 404 response might not be valid Pokemon JSON
		err := client.Get(context.Background(), "https://pokeapi.co/api/v2/pokemon/nonexistent", &result,
			WithStatus(&status))
		// Should not return error when status is captured, even if JSON parsing fails
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		// The status should still be set
		if status != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", status)
		}
	})

	t.Run("get non-existent pokemon without status capture", func(t *testing.T) {
		var pokemon Pokemon
		err := client.Get(context.Background(), "https://pokeapi.co/api/v2/pokemon/nonexistent", &pokemon)
		// Should return error when status is not captured
		if err == nil {
			t.Error("expected error for 404 response")
		}
	})
}

func TestClient_Post(t *testing.T) {
	client := &Client{}

	// Note: PokeAPI is read-only, so we'll test POST against httpbin.org
	// httpbin.org is a free HTTP testing service that echoes back request data
	t.Run("post json data", func(t *testing.T) {
		postData := map[string]interface{}{
			"pokemon": "pikachu",
			"level":   25,
			"moves":   []string{"thunderbolt", "quick-attack"},
		}

		var result map[string]interface{}
		err := client.Post(context.Background(), "https://httpbin.org/post", postData, &result)
		if err != nil {
			t.Fatalf("POST request failed: %v", err)
		}

		// httpbin.org returns the posted data in the "json" field
		if result["json"] == nil {
			t.Error("expected response to contain 'json' field")
		}

		jsonData := result["json"].(map[string]interface{})
		if jsonData["pokemon"] != "pikachu" {
			t.Errorf("expected pokemon 'pikachu', got '%v'", jsonData["pokemon"])
		}
	})

	t.Run("post with custom headers", func(t *testing.T) {
		postData := map[string]string{"test": "value"}
		var result map[string]interface{}

		err := client.Post(context.Background(), "https://httpbin.org/post", postData, &result,
			WithHeader("X-Test-Header", "test-value"))
		if err != nil {
			t.Fatalf("POST request failed: %v", err)
		}

		headers := result["headers"].(map[string]interface{})
		if headers["X-Test-Header"] != "test-value" {
			t.Errorf("expected header 'test-value', got '%v'", headers["X-Test-Header"])
		}
	})

	t.Run("post with multiple headers using WithHeaders", func(t *testing.T) {
		postData := map[string]string{"test": "value"}
		var result map[string]interface{}

		customHeaders := map[string]string{
			"X-API-Key":    "secret123",
			"X-Client-ID":  "client456",
			"X-Request-ID": "req789",
		}

		err := client.Post(context.Background(), "https://httpbin.org/post", postData, &result,
			WithHeaders(customHeaders))
		if err != nil {
			t.Fatalf("POST request with multiple headers failed: %v", err)
		}

		headers := result["headers"].(map[string]interface{})

		// Verify at least one of our custom headers is present to confirm WithHeaders works
		foundCustomHeader := false
		for key, value := range headers {
			if key == "X-Api-Key" && value == "secret123" {
				foundCustomHeader = true
				break
			}
			if key == "X-Client-Id" && value == "client456" {
				foundCustomHeader = true
				break
			}
		}

		if !foundCustomHeader {
			t.Errorf("WithHeaders test failed - no custom headers found in response. Headers: %+v", headers)
		}
	})
}

func TestClient_Patch(t *testing.T) {
	client := &Client{}

	t.Run("patch json data", func(t *testing.T) {
		patchData := map[string]interface{}{
			"pokemon": "pikachu",
			"level":   30, // leveled up!
		}

		var result map[string]interface{}
		err := client.Patch(context.Background(), "https://httpbin.org/patch", patchData, &result)
		if err != nil {
			t.Fatalf("PATCH request failed: %v", err)
		}

		jsonData := result["json"].(map[string]interface{})
		if jsonData["level"].(float64) != 30 {
			t.Errorf("expected level 30, got %v", jsonData["level"])
		}
	})
}

func TestClient_Delete(t *testing.T) {
	client := &Client{}

	t.Run("delete request", func(t *testing.T) {
		var result map[string]interface{}
		err := client.Delete(context.Background(), "https://httpbin.org/delete", &result)
		if err != nil {
			t.Fatalf("DELETE request failed: %v", err)
		}

		// httpbin.org returns request info for DELETE
		if result["url"] != "https://httpbin.org/delete" {
			t.Errorf("expected URL 'https://httpbin.org/delete', got '%v'", result["url"])
		}
	})

	t.Run("delete with status capture", func(t *testing.T) {
		var status int
		var result map[string]interface{}
		err := client.Delete(context.Background(), "https://httpbin.org/delete", &result,
			WithStatus(&status))
		if err != nil {
			t.Fatalf("DELETE request failed: %v", err)
		}

		if status != http.StatusOK {
			t.Errorf("expected status 200, got %d", status)
		}
	})
}

func TestClient_CustomMarshalUnmarshal(t *testing.T) {
	// Test with custom marshal/unmarshal functions
	client := &Client{
		MarshalFunc: func(v any) ([]byte, error) {
			// Custom marshal that adds a wrapper for POST data
			wrapped := map[string]interface{}{"data": v}
			return json.Marshal(wrapped)
		},
		UnmarshalFunc: func(data []byte, v any) error {
			// For this test, just use default unmarshal
			return json.Unmarshal(data, v)
		},
	}

	t.Run("post with custom marshal", func(t *testing.T) {
		postData := map[string]string{"pokemon": "ditto"}
		var result map[string]interface{}

		err := client.Post(context.Background(), "https://httpbin.org/post", postData, &result)
		if err != nil {
			t.Fatalf("POST request failed: %v", err)
		}

		// The custom marshal should wrap the data
		jsonData := result["json"].(map[string]interface{})
		if jsonData["data"] == nil {
			t.Error("expected custom marshal to wrap data")
		}

		wrappedData := jsonData["data"].(map[string]interface{})
		if wrappedData["pokemon"] != "ditto" {
			t.Errorf("expected pokemon 'ditto', got '%v'", wrappedData["pokemon"])
		}
	})
}
