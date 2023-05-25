package config

import (
	"bytes"
	"io"

	"gopkg.in/yaml.v2"
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

func SerializeToYaml(data interface{}) ([]byte, error) {
	cfg := NewEmpty()

	writableConfig, isWritable := cfg.(WritableConfiguration)

	if !isWritable {
		panic("Configuration isn't writable. No way to continue normal work")
	}

	err := writableConfig.Dematerialize(data)
	if err != nil {
		return nil, err
	}

	values, err := cfg.GetValues()
	if err != nil {
		return nil, err
	}

	buffer := bytes.Buffer{}

	err = WriteMapToYaml(&values, &buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
