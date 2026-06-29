package config

import (
	"errors"

	"github.com/spf13/viper"
)

func ErrConfigNotFound(err string) error {
	return errors.New("config not found: " + err)
}

/*
Env описывает конфигурацию приложения.
*/
type Env struct {
	Portainer portainerConfig `mapstructure:"portainer"`
	Git       gitConfig       `mapstructure:"git"`
	Cors      corsConfig      `mapstructure:"cors"`
}

type portainerConfig struct {
	Realm string `mapstructure:"realm"`
	Token string `mapstructure:"token"`
	Teams []int  `mapstructure:"teams"`
}

type gitConfig struct {
	Realm string `mapstructure:"realm"`
	Token string `mapstructure:"token"`
}

type corsConfig struct {
	Origin string `mapstructure:"origin"`
}

func NewEnv() (*Env, error) {
	var env Env

	viper.SetConfigFile("./config.yml")

	err := viper.ReadInConfig()
	if err != nil {
		return &env, err
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		return &env, err
	}

	if err := env.Validate(); err != nil {
		return &env, err
	}

	return &env, nil
}

/*
Validate валидирует файл конфигурации.

Параметры: текущий экземпляр Env.
Результат: nil, если конфигурация валидна, иначе ошибка с описанием отсутствующего поля.
*/
func (env *Env) Validate() error {
	if env.Portainer.Realm == "" {
		return ErrConfigNotFound("Portainer Realm")
	}
	if env.Portainer.Token == "" {
		return ErrConfigNotFound("Portainer Token")
	}
	if len(env.Portainer.Teams) == 0 {
		return ErrConfigNotFound("Portainer Teams")
	}

	if env.Git.Realm == "" {
		return ErrConfigNotFound("Git Realm")
	}
	if env.Git.Token == "" {
		return ErrConfigNotFound("Git Token")
	}
	return nil
}
