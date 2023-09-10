package main

import (
	"errors"
	"frps-panel/pkg/server"
	"frps-panel/pkg/server/controller"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const version = "1.4.0"

var (
	showVersion bool
	configFile  string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "version of frps-panel")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "./frps-panel.ini", "config file of frps-panel")
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

		common, tokens, ports, domains, subdomains, iniFile, err := ParseConfigFile(configFile)
		if err != nil {
			log.Printf("fail to start frps-panel : %v", err)
			return err
		}
		s, err := server.New(
			rootDir,
			controller.HandleController{
				CommonInfo: common,
				Tokens:     tokens,
				Ports:      ports,
				Domains:    domains,
				Subdomains: subdomains,
				ConfigFile: configFile,
				IniFile:    iniFile,
				Version:    version,
			})
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

func ParseConfigFile(file string) (controller.CommonInfo, map[string]controller.TokenInfo, map[string][]string, map[string][]string, map[string][]string, *ini.File, error) {
	common := controller.CommonInfo{}
	users := make(map[string]controller.TokenInfo)
	ports := make(map[string][]string)
	domains := make(map[string][]string)
	subdomains := make(map[string][]string)

	iniFile, err := ini.LoadSources(ini.LoadOptions{
		Insensitive:         false,
		InsensitiveSections: false,
		InsensitiveKeys:     false,
		IgnoreInlineComment: true,
		AllowBooleanKeys:    true,
	}, file)
	if err != nil {
		var pathError *fs.PathError
		if errors.As(err, &pathError) {
			log.Printf("token file %s not found", file)
		} else {
			log.Printf("fail to parse token file %s : %v", file, err)
		}
		return common, nil, nil, nil, nil, iniFile, err
	}

	commonSection, err := iniFile.GetSection("common")
	if err != nil {
		log.Printf("fail to get [common] section from file %s : %v", file, err)
		return common, nil, nil, nil, nil, iniFile, err
	}
	common.PluginAddr = commonSection.Key("plugin_addr").MustString("0.0.0.0")
	common.PluginPort = commonSection.Key("plugin_port").MustInt(7200)
	common.User = commonSection.Key("admin_user").Value()
	common.Pwd = commonSection.Key("admin_pwd").Value()
	common.DashboardAddr = commonSection.Key("dashboard_addr").MustString("127.0.0.1")
	common.DashboardPort = commonSection.Key("dashboard_port").MustInt(7500)
	common.DashboardUser = commonSection.Key("dashboard_user").Value()
	common.DashboardPwd = commonSection.Key("dashboard_pwd").Value()

	portsSection, err := iniFile.GetSection("ports")
	if err != nil {
		log.Printf("fail to get [ports] section from file %s : %v", file, err)
		return common, nil, nil, nil, nil, iniFile, err
	}
	for _, key := range portsSection.Keys() {
		user := key.Name()
		value := key.Value()
		port := strings.Split(controller.TrimAllSpaceReg.ReplaceAllString(value, ""), ",")
		ports[user] = port
	}

	domainsSection, err := iniFile.GetSection("domains")
	if err != nil {
		log.Printf("fail to get [domains] section from file %s : %v", file, err)
		return common, nil, nil, nil, nil, iniFile, err
	}
	for _, key := range domainsSection.Keys() {
		user := key.Name()
		value := key.Value()
		domain := strings.Split(controller.TrimAllSpaceReg.ReplaceAllString(value, ""), ",")
		domains[user] = domain
	}

	subdomainsSection, err := iniFile.GetSection("subdomains")
	if err != nil {
		log.Printf("fail to get [subdomains] section from file %s : %v", file, err)
		return common, nil, nil, nil, nil, iniFile, err
	}
	for _, key := range subdomainsSection.Keys() {
		user := key.Name()
		value := key.Value()
		subdomain := strings.Split(controller.TrimAllSpaceReg.ReplaceAllString(value, ""), ",")
		subdomains[user] = subdomain
	}

	usersSection, err := iniFile.GetSection("users")
	if err != nil {
		log.Printf("fail to get [users] section from file %s : %v", file, err)
		return common, nil, nil, nil, nil, iniFile, err
	}

	disabledSection, err := iniFile.GetSection("disabled")
	if err != nil {
		log.Printf("fail to get [disabled] section from file %s : %v", file, err)
		return common, nil, nil, nil, nil, iniFile, err
	}

	keys := usersSection.Keys()
	for _, key := range keys {
		comment, found := strings.CutPrefix(key.Comment, ";")
		if !found {
			comment, found = strings.CutPrefix(comment, "#")
		}
		token := controller.TokenInfo{
			User:       key.Name(),
			Token:      key.Value(),
			Comment:    comment,
			Ports:      strings.Join(ports[key.Name()], ","),
			Domains:    strings.Join(domains[key.Name()], ","),
			Subdomains: strings.Join(subdomains[key.Name()], ","),
			Status:     !(disabledSection.HasKey(key.Name()) && disabledSection.Key(key.Name()).Value() == "disable"),
		}
		users[token.User] = token
	}

	return common, users, ports, domains, subdomains, iniFile, nil
}
