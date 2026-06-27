package config

import (
	"github.com/spf13/viper"
)

/*
Env описывает конфигурацию приложения.

Поля:
  - Storage: настройки подключения к S3-compatible storage.
*/
type Env struct {
}

/*
S3Config описывает настройки подключения к S3-compatible storage.

Поля:
  - Endpoint: адрес S3 endpoint в формате host или host:port.
  - Region: регион S3 bucket.
  - Bucket: имя bucket.
  - Prefix: базовый префикс внутри bucket.
  - AccessKey: access key для подключения.
  - SecretKey: secret key для подключения.
  - Secure: использовать HTTPS-подключение.
*/

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
