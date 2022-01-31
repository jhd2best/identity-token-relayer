package config

type Config struct {
	Debug *DebugConfig
	Db    *DbConfig
	Eth   *EthConfig
}

type DebugConfig struct {
	Verbose       bool
	DisableCron   bool
	DisableSentry bool
	SentryDSN     string
}

type DbConfig struct {
	ServiceAccountPath string
}

type EthConfig struct {
	RpcEndpoints string
}
