package main

import (
	"frps-panel/pkg/server"
	"frps-panel/pkg/server/controller"
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const version = "1.6.0"

var (
	showVersion bool
	configFile  string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of frps-panel")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "./frps-panel.toml", "config file of frps-panel")
}

var rootCmd = &cobra.Command{
	Use:   "frps-panel",
	Short: "frps-panel is the server plugin of frp to support multiple users.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			log.Println(version)
			return nil
		}
		executable, err := os.Executable()
		if err != nil {
			log.Printf("error get program path: %v", err)
			return err
		}
		rootDir := filepath.Dir(executable)

		configDir := filepath.Dir(configFile)
		tokensFile := filepath.Join(configDir, "frps-tokens.toml")

		config, tls, err := parseConfigFile(configFile, tokensFile)
		if err != nil {
			log.Printf("fail to start frps-panel : %v", err)
			return err
		}

		s, err := server.New(
			rootDir,
			config,
			tls,
		)
		if err != nil {
			return err
		}
		err = s.Run()
		if err != nil {
			return err
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func parseConfigFile(configFile, tokensFile string) (controller.HandleController, server.TLS, error) {
	var common controller.Common
	var tokens controller.Tokens
	_, err := toml.DecodeFile(configFile, &common)
	if err != nil {
		log.Fatalf("decode config file %v error: %v", configFile, err)
	}

	_, err = toml.DecodeFile(tokensFile, &tokens)
	if err != nil {
		log.Fatalf("decode token file %v error: %v", tokensFile, err)
	}

	common.Common.DashboardTls = strings.HasPrefix("https://", strings.ToLower(common.Common.DashboardAddr))

	tls := server.TLS{
		Enable:   common.Common.TlsMode,
		Protocol: "HTTP",
	}

	if tls.Enable {
		tls.Protocol = "HTTPS"

		if strings.TrimSpace(tls.Cert) == "" || strings.TrimSpace(tls.Key) == "" {
			tls.Enable = false
			tls.Protocol = "HTTP"
			log.Printf("fail to enable tls: tls cert or key not exist, use http as default.")
		}
	}

	return controller.HandleController{
		CommonInfo: common.Common,
		Tokens:     tokens.Tokens,
		Version:    version,
		ConfigFile: configFile,
		TokensFile: tokensFile,
	}, tls, nil
}
