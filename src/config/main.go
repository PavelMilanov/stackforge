package config

import (
	"github.com/spf13/viper"
)

/*
Env описывает конфигурацию приложения.
*/
type Env struct {
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
	return nil
}
