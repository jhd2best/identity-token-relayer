package config

import (
	"github.com/spf13/viper"
	"strings"
)

var config Config

func InitConfig(path string) {
	instance := viper.New()
	instance.SetConfigFile(path)
	instance.SetConfigType("yaml")
	instance.SetEnvPrefix("ITR")
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
			Verbose:      true,
			LogPath:      "",
			DisableCron:  true,
			SentryDSN:    "",
			PagerDutyKey: "",
		},
		Db: &DbConfig{
			ServiceAccountPath: "./firebase-service-account.json",
		},
		RPC: &RPCConfig{
			Listen: "127.0.0.1",
			Port:   28888,
		},
		Aws: &AwsConfig{
			Profile: "mainnet",
			Region:  "us-west-2",
		},
		Eth: &EthConfig{
			RpcEndpoints:  "https://kovan.infura.io/v3",
			MoralisApiKey: "",
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
