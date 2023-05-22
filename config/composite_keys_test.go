package config

import (
	"strings"
	"testing"
)

func TestFindValueByCompositeKeyPath(t *testing.T) {
	config, err := FromYamlReader(strings.NewReader(`
                root1,root2: 
                    root3, root4: 
                        key: value`))

	if err != nil {
		t.Fatal(err)
	}

	values, err := config.GetValues()

	if err != nil {
		t.Fatal("Failed to read configuration value")
	}

	value := values.FindPath("root1/root4/key")

	if value == nil {
		t.Fatal("Failed to find key by composite path")
	}

	if value != "value" {
		t.Fatal("Value found by composite path is wrong")
	}
}
