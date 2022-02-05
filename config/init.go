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
			LogPath:       "",
			DisableCron:   false,
			DisableSentry: false,
			SentryDSN:     "",
		},
		Db: &DbConfig{
			ServiceAccountPath: "./firebase-service-account.json",
		},
		Aws: &AwsConfig{
			Profile: "mainnet",
			Region:  "us-west-2",
		},
		Eth: &EthConfig{
			RpcEndpoints: "https://kovan.infura.io/v3",
		},
		Hmy: &HmyConfig{
			RpcEndpoints:              "https://api.s0.b.hmny.io",
			PrivateKeyPath:            "./harmony-testnet.key",
			OpenKMS:                   false,
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
