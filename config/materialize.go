package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

func NewDecoderConfig(applyOpts ...func(*mapstructure.DecoderConfig)) *mapstructure.DecoderConfig {
	decoderConfig := &mapstructure.DecoderConfig{
		Metadata:         nil,
		TagName:          "yaml",
		WeaklyTypedInput: true,
	}
	for _, applyOpt := range applyOpts {
		applyOpt(decoderConfig)
	}
	return decoderConfig
}

func WithFailOnUnusedFields(decoderConfig *mapstructure.DecoderConfig) {
	decoderConfig.ErrorUnused = true
}

func WithResult(result interface{}) func(*mapstructure.DecoderConfig) {
	return func(config *mapstructure.DecoderConfig) {
		config.Result = result
	}
}

func WithResolvingReferences(values MapStr) func(*mapstructure.DecoderConfig) {
	return func(config *mapstructure.DecoderConfig) {
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			resolveReferencesInStringFields(values),
			resolveReferencesInUntypedMapFields(values))
	}
}

func Materialize(decoderConfig *mapstructure.DecoderConfig, values interface{}) error {
	decoder, err := mapstructure.NewDecoder(decoderConfig)

	if err != nil {
		return errors.WithMessage(err, "Failed to create structure decoder")
	}

	err = decoder.Decode(values)

	if err != nil {
		return errors.WithMessage(err, "Failed to materialize data structure from values")
	}

	return nil
}
