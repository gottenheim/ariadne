package config

import "testing"

func TestFindJsonStringByPath(t *testing.T) {
	json := `{
          "auth": {
            "client_token": "some_client_token"
          }
    }`
	result, err := FindJsonStringByPath(json, "auth/client_token")

	if err != nil {
		t.Fatal(err)
	}

	if result != "some_client_token" {
		t.Fatal("Should find string value by path")
	}
}
