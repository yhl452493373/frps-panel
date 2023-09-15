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
	var config controller.Config
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		log.Fatalf("decode config file %v error: %v", configFile, err)
	}

	_, err = toml.DecodeFile(tokensFile, &config)
	if err != nil {
		log.Fatalf("decode token file %v error: %v", tokensFile, err)
	}

	config.Common.DashboardTls = strings.HasPrefix("https://", strings.ToLower(config.Common.DashboardAddr))

	//f, err := os.Create("/Volumes/Work/Git Sources/frps-panel/config/frps-panel-new.toml")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if err := toml.NewEncoder(f).Encode(config); err != nil {
	//	// failed to encode
	//	log.Fatal(err)
	//}
	//if err := f.Close(); err != nil {
	//	// failed to close the file
	//	log.Fatal(err)
	//}

	tls := server.TLS{
		Enable:   config.Common.TlsMode,
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
		CommonInfo: config.Common,
		Tokens:     config.Tokens.Tokens,
		Version:    version,
		ConfigFile: configFile,
		TokensFile: tokensFile,
	}, tls, nil
}
