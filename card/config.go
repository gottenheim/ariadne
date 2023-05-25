package card

import (
	"bufio"
	"fmt"

	"github.com/gottenheim/ariadne/details/config"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type Config struct {
	AnswerFileName string
}

func LoadConfig(fs afero.Fs, configFilePath string) (*Config, error) {
	configFileExists, err := afero.Exists(fs, configFilePath)

	if err != nil {
		return nil, err
	}

	if !configFileExists {
		return nil, errors.New(fmt.Sprintf("Configuration file %s does not exist", configFilePath))
	}

	configFile, err := fs.Open(configFilePath)

	if err != nil {
		return nil, errors.WithMessagef(err, "Failed to open config file %s", configFilePath)
	}

	defer func() {
		_ = configFile.Close()
	}()

	cfg, err := config.FromYamlReader(bufio.NewReader(configFile))

	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create envConfig configuration")
	}

	config, err := MaterializeConfig(cfg)

	if err != nil {
		return nil, errors.WithMessage(err, "Failed to materialize stored envConfig")
	}

	return config, nil
}

func MaterializeConfig(cfg config.Configuration) (*Config, error) {
	config := &Config{}

	ok, err := cfg.StrictMaterializeAt("config", &config)

	if err != nil {
		return nil, errors.WithMessage(err, "Failed to materialize configuration")
	}

	if !ok {
		return nil, errors.New("Configuration has wrong format")
	}

	return config, nil
}
