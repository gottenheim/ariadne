package config

import (
	"github.com/pkg/errors"
)

func ResolveReferences(config Configuration, resolutionContext Configuration) (Configuration, error) {
	values, err := config.GetValues()

	if err != nil {
		return nil, errors.WithMessage(err, "Failed to get configuration values")
	}

	contextValues, err := resolutionContext.GetValues()

	result := map[string]interface{}{}

	decoderConfig := NewDecoderConfig(
		WithResult(&result),
		WithResolvingReferences(contextValues),
		WithFailOnUnusedFields)

	err = Materialize(decoderConfig, values)

	if err != nil {
		return nil, errors.WithMessage(err, "Failed to replace references in configuration")
	}

	return FromValues(result), nil
}
