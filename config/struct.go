package config

type Config struct {
	Debug *DebugConfig
	Db    *DbConfig
	Aws   *AwsConfig
	Eth   *EthConfig
	Hmy   *HmyConfig
}

type DebugConfig struct {
	Verbose       bool
	LogPath       string
	DisableCron   bool
	DisableSentry bool
	SentryDSN     string
}

type DbConfig struct {
	ServiceAccountPath string
}

type AwsConfig struct {
	Profile string
	Region  string
}

type EthConfig struct {
	RpcEndpoints string
}

type HmyConfig struct {
	RpcEndpoints              string
	PrivateKeyPath            string
	OpenKMS                   bool
	OwnershipValidatorAddress string
}
