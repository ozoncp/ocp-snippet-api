package configuration

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Grpc struct {
		GrpcPort int `yaml:"grpcPort"`
		HttpPort int `yaml:"httpPort"`
	}
	Db struct {
		Table     string `yaml:"table"`
		ChunkSize uint   `yaml:"chunkSize"`
	}
}

var configInstance *Config = nil
var once sync.Once

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func GetInstance() *Config {
	once.Do(func() {
		var err error
		configInstance, err = new()
		if err != nil {
			log.Err(err)
		}
	})
	return configInstance
}

func validatePath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return errors.New("configuration file not found")
	}
	if s.IsDir() {
		return errors.New("configuration path is a directory")
	}
	return nil
}
func new() (*Config, error) {
	const path string = "../../configuration.yml"

	if err := validatePath(path); err != nil {
		log.Err(err)
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error().Msg(err.Error())
		return nil, err
	}

	res := &Config{}

	err = yaml.Unmarshal([]byte(data), &res)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	return res, nil
}
