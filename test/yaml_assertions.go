package test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/gottenheim/ariadne/config"
)

func AssertIdenticalYamlStrings(t *testing.T, first string, second string) {
	firstConf, err := config.FromYamlReader(strings.NewReader(first))

	if err != nil {
		t.Error("First string doesn't have yaml format")
	}

	firstValues, _ := firstConf.GetValues()

	secondConf, err := config.FromYamlReader(strings.NewReader(second))

	if err != nil {
		t.Error("Second string doesn't have yaml format")
	}

	secondValues, _ := secondConf.GetValues()

	if !reflect.DeepEqual(firstValues, secondValues) {
		t.Errorf("Strings have different yaml content. First: %s, second: %s", first, second)
	}
}
