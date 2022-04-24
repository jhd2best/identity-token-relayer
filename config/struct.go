package config

type Config struct {
	Debug *DebugConfig
	Db    *DbConfig
	RPC   *RPCConfig
	Aws   *AwsConfig
	Eth   *EthConfig
	Hmy   *HmyConfig
}

type DebugConfig struct {
	Verbose      bool
	LogPath      string
	DisableCron  bool
	SentryDSN    string
	PagerDutyKey string
}

type DbConfig struct {
	ServiceAccountPath string
}

type AwsConfig struct {
	Profile string
	Region  string
}

type EthConfig struct {
	RpcEndpoints  string
	MoralisApiKey string
}

type HmyConfig struct {
	RpcEndpoints              string
	PrivateKeyPath            string
	OpenKMS                   bool
	OwnershipValidatorAddress string
}

type RPCConfig struct {
	Listen string
	Port   uint16
}
