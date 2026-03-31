package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	HTTPAddress string `yaml:"http_address" env-default:":8080"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
}

// MustLoad loads the config. "Must" = panics on error.
// This is fine for application startup — no config means no reason to run.
func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

// MustLoadPath loads config from a specific path.
// Extracted into a separate function for use in tests
// (pass a test config path directly, without flags).
func MustLoadPath(configPath string) *Config {
	// Check if the file exists before parsing —
	// to give a clear error instead of "cannot read config".
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	// cleanenv parses the YAML file and maps it to the struct via tags.
	// It also validates env-required and applies env-default values.
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath determines where to get the config path from.
// Priority: command line flag > environment variable.
// This allows flexible startup:
//   go run ./cmd --config=config/local.yaml   (flag)
//   CONFIG_PATH=config/local.yaml go run ./cmd (env).
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	// If flag is not provided — try environment variable.
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
