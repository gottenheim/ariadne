package config

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func ReadJsonToMap(reader io.Reader) (MapStr, error) {
	decoder := json.NewDecoder(reader)

	values := make(MapStr)
	err := decoder.Decode(&values)

	if err != nil {
		return nil, err
	}

	return values, nil
}

func FindJsonPath(json string, path string) (interface{}, error) {
	values, err := ReadJsonToMap(strings.NewReader(json))

	if err != nil {
		return nil, err
	}

	return values.FindPath(path), nil
}

func FindJsonStringByPath(json string, path string) (string, error) {
	values, err := ReadJsonToMap(strings.NewReader(json))

	if err != nil {
		return "", err
	}

	val := values.FindPath(path)

	if val == nil {
		return "", nil
	}

	str, isStr := val.(string)

	if isStr {
		return str, nil
	}

	return fmt.Sprintf("%v", val), nil
}
