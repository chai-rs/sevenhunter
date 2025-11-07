package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

func MustNew[T any](prefix string) *T {
	conf, err := New[T](prefix)
	if err != nil {
		panic(err)
	}
	return conf
}

func New[T any](prefix string) (*T, error) {
	filepath := flag.String("env", "", "path to .env file")
	flag.Parse()

	if filepath != nil && *filepath != "" {
		if err := exportEnvironment(*filepath); err != nil {
			return nil, fmt.Errorf("failed to load env file: %w", err)
		}
	}

	var conf T
	if err := envconfig.Process(prefix, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func exportEnvironment(filepath string) error {
	viper.SetConfigFile(filepath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	for k, v := range viper.AllSettings() {
		if err := os.Setenv(strings.ToUpper(k), fmt.Sprint(v)); err != nil {
			return err
		}
	}

	return nil
}
