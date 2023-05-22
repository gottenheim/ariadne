package config

import (
	"gopkg.in/yaml.v2"
	"io"
)

func ReadYamlToMap(reader io.Reader) (MapStr, error) {
	decoder := yaml.NewDecoder(reader)

	rawValues := make(MapItf)
	err := decoder.Decode(&rawValues)

	if err != nil {
		return nil, err
	}

	return cleanUpInterfaceMap(rawValues), nil
}

func MapToYamlBytes(yamlMap *MapStr) ([]byte, error) {
	byteData, err := yaml.Marshal(&yamlMap)
	if err != nil {
		return nil, err
	}

	return byteData, nil
}

func WriteMapToYaml(values *MapStr, writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)

	err := encoder.Encode(values)
	if err != nil {
		return err
	}

	return nil
}

func WriteMapSliceToYaml(values yaml.MapSlice, writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)

	err := encoder.Encode(values)
	if err != nil {
		return err
	}

	return nil
}
