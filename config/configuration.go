package config

import (
	"encoding/json"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"io"
)

type Configuration interface {
	GetValues() (MapStr, error)

	Materialize(data interface{}) error

	StrictMaterializeAt(path string, data interface{}) (bool, error)

	MaterializeAt(path string, data interface{}) (bool, error)

	NestedConfig(key string) (Configuration, error)
}

type WritableConfiguration interface {
	Dematerialize(data interface{}) error
}

func FromYamlReader(reader io.Reader) (Configuration, error) {
	values, err := ReadYamlToMap(reader)

	if err != nil {
		return nil, err
	}

	config := configuration{
		values: values,
	}

	return &config, nil
}

func FromJsonReader(reader io.Reader) (Configuration, error) {
	decoder := json.NewDecoder(reader)

	values := make(MapStr)
	err := decoder.Decode(&values)

	config := configuration{
		values: values.deepClone(),
	}

	if err != nil {
		return nil, err
	}

	return &config, nil
}

func NewEmpty() Configuration {
	return &configuration{
		values: MapStr{},
	}
}

func FromValues(values MapStr) Configuration {
	return &configuration{
		values: values,
	}
}

type configuration struct {
	values MapStr
}

func (c *configuration) GetValues() (MapStr, error) {
	return c.values.deepClone(), nil
}

func (c *configuration) Materialize(data interface{}) error {
	decoderConfig := NewDecoderConfig(
		WithResult(data),
		WithResolvingReferences(c.values))

	return Materialize(decoderConfig, c.values)
}

func (c *configuration) MaterializeAt(path string, data interface{}) (bool, error) {
	nestedValues := c.values.FindPath(path)

	if nestedValues == nil {
		return false, nil
	}

	decoderConfig := NewDecoderConfig(
		WithResult(data),
		WithResolvingReferences(c.values))

	err := Materialize(decoderConfig, nestedValues)

	if err != nil {
		return false, errors.WithMessage(err, "Failed to materialize data structure from values")
	}

	return true, nil
}

func (c *configuration) StrictMaterializeAt(path string, data interface{}) (bool, error) {
	nestedValues := c.values.FindPath(path)

	if nestedValues == nil {
		return false, nil
	}

	decoderConfig := NewDecoderConfig(
		WithResult(data),
		WithResolvingReferences(c.values),
		WithFailOnUnusedFields)

	err := Materialize(decoderConfig, nestedValues)

	if err != nil {
		return false, errors.WithMessage(err, "Failed to materialize data structure from values")
	}

	return true, nil
}

func (c *configuration) NestedConfig(path string) (Configuration, error) {
	nestedValues := c.values.FindPath(path)
	if nestedValues == nil {
		return c, nil
	}

	result := c.values.deepClone()

	err := mergo.Merge(&result, nestedValues, mergo.WithAppendSlice, mergo.WithOverride, mergo.WithTypeCheck)

	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create nested configuration")
	}

	result.RemovePath(path)

	return &configuration{values: result}, nil
}

func (c *configuration) Dematerialize(data interface{}) error {
	decoderConfig := NewDecoderConfig(
		WithResult(&c.values),
		WithFailOnUnusedFields)

	return Materialize(decoderConfig, data)
}

func Combine(configs []Configuration) (Configuration, error) {
	if len(configs) == 0 {
		return nil, nil
	}

	for _, config := range configs {
		if config == nil {
			return nil, errors.New("trying to combine with nil configuration")
		}
	}

	if len(configs) > 1 {
		return newCascade(configs), nil
	} else {
		return configs[0], nil
	}
}

func newCascade(configurations []Configuration) Configuration {
	config := &cascadeConfiguration{configurations: configurations}

	return config
}

type cascadeConfiguration struct {
	configurations []Configuration
}

func (c *cascadeConfiguration) GetValues() (MapStr, error) {
	return c.GetValuesAt("")
}

func (c *cascadeConfiguration) GetValuesAt(path string) (MapStr, error) {
	var result MapStr

	for _, config := range c.configurations {
		values, err := config.GetValues()

		if err != nil {
			return nil, err
		}

		nestedValues := values.FindPath(path)

		if nestedValues != nil {
			if result == nil {
				result = MapStr{}
			}

			err = mergo.Merge(&result, nestedValues, mergo.WithAppendSlice, mergo.WithOverride, mergo.WithTypeCheck)

			if err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func (c *cascadeConfiguration) Materialize(data interface{}) error {
	values, err := c.GetValues()

	if err != nil {
		return err
	}

	decoderConfig := NewDecoderConfig(
		WithResult(data),
		WithResolvingReferences(values))

	return Materialize(decoderConfig, values)
}

func (c *cascadeConfiguration) MaterializeAt(path string, data interface{}) (bool, error) {
	values, err := c.GetValues()

	if err != nil {
		return false, err
	}

	if values == nil {
		return false, nil
	}

	nestedValues := values.FindPath(path)

	if nestedValues == nil {
		return false, nil
	}

	decoderConfig := NewDecoderConfig(
		WithResult(data),
		WithResolvingReferences(values))

	err = Materialize(decoderConfig, nestedValues)

	if err != nil {
		return false, errors.WithMessage(err, "Failed to materialize data structure from values")
	}

	return true, nil
}

func (c *cascadeConfiguration) StrictMaterializeAt(path string, data interface{}) (bool, error) {
	values, err := c.GetValues()

	if err != nil {
		return false, err
	}

	if values == nil {
		return false, nil
	}

	nestedValues := values.FindPath(path)

	if nestedValues == nil {
		return false, nil
	}

	decoderConfig := NewDecoderConfig(
		WithResult(data),
		WithResolvingReferences(values),
		WithFailOnUnusedFields)

	err = Materialize(decoderConfig, nestedValues)

	if err != nil {
		return false, errors.WithMessage(err, "Failed to materialize data structure from values")
	}

	return true, nil
}

func (c *cascadeConfiguration) NestedConfig(key string) (Configuration, error) {
	result := MapStr{}

	for _, config := range c.configurations {
		nestedConfig, err := config.NestedConfig(key)

		if err != nil {
			return nil, errors.WithMessage(err, "Failed to create nested configuration")
		}

		values, err := nestedConfig.GetValues()

		if err != nil {
			return nil, errors.WithMessage(err, "Failed to get nested configuration values")
		}

		err = mergo.Merge(&result, values, mergo.WithAppendSlice, mergo.WithOverride, mergo.WithTypeCheck)

		if err != nil {
			return nil, err
		}
	}

	return &configuration{values: result}, nil
}
