package config

import (
	"reflect"
	"strings"
	"testing"
)

func AssertIdenticalYamlStrings(t *testing.T, first string, second string) {
	firstConf, err := FromYamlReader(strings.NewReader(first))

	if err != nil {
		t.Error("First string doesn't have yaml format")
	}

	firstValues, _ := firstConf.GetValues()

	secondConf, err := FromYamlReader(strings.NewReader(second))

	if err != nil {
		t.Error("Second string doesn't have yaml format")
	}

	secondValues, _ := secondConf.GetValues()

	if !reflect.DeepEqual(firstValues, secondValues) {
		t.Errorf("Strings have different yaml content. First: %s, second: %s", first, second)
	}
}
