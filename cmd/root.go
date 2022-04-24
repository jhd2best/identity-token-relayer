package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"identity-token-relayer/config"
	"identity-token-relayer/log"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "itoken",
		Short: "itoken can mapping your NFT from Ethereum to Harmony",
		Long: `Itoken means identity token which can mapping your NFT from Ethereum to Harmony
You can get more information at https://github.com/harmony-one/identity-token-relayer`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config.yaml", "the path of config file")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.GetLogger().Fatal("execute cobra failed.", zap.String("error", err.Error()))
	}
}

func initConfig() {
	// load and inti config
	config.InitConfig(cfgFile)

	// init logger
	log.InitLog()
}
