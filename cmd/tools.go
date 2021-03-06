package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"identity-token-relayer/hmy"
	"identity-token-relayer/log"
	"identity-token-relayer/model"
	"identity-token-relayer/tools"
)

func init() {
	toolsCmd.Flags().StringVar(&origin, "origin", "", "the address of origin contract")
	toolsCmd.Flags().StringVar(&mapping, "mapping", "", "the address of mapping contract")
	toolsCmd.Flags().StringVar(&name, "name", "", "the name of project")
	toolsCmd.Flags().StringVar(&symbol, "symbol", "", "the symbol of project")
	toolsCmd.Flags().StringVar(&baseUrl, "base_url", "", "the baseUrl of project")
	toolsCmd.Flags().StringVar(&mode, "mode", "normal", "project init mode [normal/deep]")
	toolsCmd.Flags().StringVar(&source, "source", "moralis", "source for getting nft owners [contract/moralis]")
	toolsCmd.Flags().Int64Var(&tokenId, "token_id", 0, "the token id you want to get info")
	toolsCmd.Flags().Int64Var(&indexOffset, "index_offset", 0, "the offset of token index")

	rootCmd.AddCommand(toolsCmd)
}

var (
	origin, mapping, name, symbol, baseUrl, mode, source string
	tokenId, indexOffset                                 int64

	toolsCmd = &cobra.Command{
		Use:   "tools",
		Short: "some tools for itoken",
		Long: `tools as follows are supported:
		- new: add and init new NFT project
	`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var initErr error

			switch args[0] {
			case "new":
				initErr = model.InitDb()
				initErr = hmy.InitClient()
				if initErr != nil {
					log.GetLogger().Fatal("init database failed.", zap.String("error", initErr.Error()))
				}

				tools.AutoImportFromChain(origin, mapping, name, symbol, baseUrl, mode, source, indexOffset)
			case "get_supply":
				if mapping == "" {
					log.GetLogger().Fatal("mapping contract must be provided.")
				}
				tools.GetMappingTotalSupply(mapping)
			case "get_owner":
				tools.GetOriginOwner(origin, tokenId)
			}
		},
	}
)
