package config

import (
	"github.com/spf13/viper"
	"strings"
)

var config Config

func init() {
	instance := viper.New()

	instance.AddConfigPath(".")

	instance.SetConfigType("yaml")
	instance.SetConfigName("config.yaml")

	instance.SetEnvPrefix("itr")
	instance.AutomaticEnv()
	instance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := instance.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	config = Config{
		Debug: &DebugConfig{
			Verbose:       false,
			DisableCron:   false,
			DisableSentry: false,
			SentryDSN:     "",
		},
		Db: &DbConfig{
			ServiceAccountPath: "./firebase-service-account.json",
		},
		Aws: &AwsConfig{
			Profile: "mainnet",
			Region: "us-west-2",
		},
		Eth: &EthConfig{
			RpcEndpoints: "",
		},
		Hmy: &HmyConfig{
			RpcEndpoints: "https://api.harmony.one",
			PrivateKeyPath: "",
			OpenKMS: false,
			OwnershipValidatorAddress: "",
		},
	}

	err = instance.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
}

func Get() *Config {
	return &config
}
