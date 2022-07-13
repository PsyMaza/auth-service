package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var cfg *Config

// New returns services config
func New() *Config {
	if cfg != nil {
		return cfg
	}

	return &Config{}
}

// Database - contains all parameters databases connection.
type Database struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Migrations string `yaml:"migrations"`
	Name       string `yaml:"name"`
	Timeout    int    `yaml:"timeout"`
}

// Http - contains parameter rest json connection.
type Http struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	DebugPort       int    `yaml:"debugPort"`
	SwaggerPort     int    `yaml:"swaggerPort"`
	ShutdownTimeout int    `yaml:"shutdownTimeout"`
	ReadTimeout     int    `yaml:"readTimeout"`
	WriteTimeout    int    `yaml:"writeTimeout"`
	IdleTimeout     int    `yaml:"idleTimeout"`
}

// App - contains all parameters project information.
type App struct {
	Debug       bool   `yaml:"debug"`
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Version     string `yaml:"version"`
}

type Jwt struct {
	SecretKey  string `yaml:"secretKey"`
	AtLifeTime int    `yaml:"atLifeTime"`
	RtLifeTime int    `yaml:"rtLifeTime"`
}

// Metrics - contains all parameters metrics information.
type Metrics struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
	Path string `yaml:"path"`
}

// Jaeger - contains all parameters metrics information.
type Jaeger struct {
	Service string `yaml:"service"`
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
}

// Grpc - contains parameter address grpc.
type Grpc struct {
	Port              int    `yaml:"port"`
	MaxConnectionIdle int64  `yaml:"maxConnectionIdle"`
	Timeout           int64  `yaml:"timeout"`
	MaxConnectionAge  int64  `yaml:"maxConnectionAge"`
	Host              string `yaml:"host"`
}

// Config - contains all configuration parameters in config package.
type Config struct {
	App      App      `yaml:"app"`
	Jwt      Jwt      `yaml:"jwt"`
	Http     Http     `yaml:"http"`
	Database Database `yaml:"database"`
	Metrics  Metrics  `yaml:"metrics"`
	Jaeger   Jaeger   `yaml:"jaeger"`
	Grpc     Grpc     `yaml:"grpc"`
}

// ReadConfigYML - read configurations from file and init instance Config.
func ReadConfigYML(filePath string) error {
	if cfg != nil {
		return nil
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	return nil
}
