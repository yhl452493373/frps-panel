package controller

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
)

type HandleController struct {
	CommonInfo CommonInfo
	Tokens     map[string]TokenInfo
	Ports      map[string][]string
	Domains    map[string][]string
	Subdomains map[string][]string
	ConfigFile string
	IniFile    *ini.File
	Version    string
}

func NewHandleController(config *HandleController) *HandleController {
	return config
}

func (c *HandleController) Register(rootDir string, engine *gin.Engine) {
	assets := filepath.Join(rootDir, "assets")
	_, err := os.Stat(assets)
	if err != nil && !os.IsExist(err) {
		assets = "./assets"
	}

	engine.Delims("${", "}")
	engine.LoadHTMLGlob(filepath.Join(assets, "templates/*"))
	engine.POST("/handler", c.MakeHandlerFunc())
	engine.Static("/static", filepath.Join(assets, "static"))
	engine.GET("/lang.json", c.MakeLangFunc())
	engine.GET(LoginUrl, c.MakeLoginFunc())
	engine.POST(LoginUrl, c.MakeLoginFunc())
	engine.GET(LogoutUrl, c.MakeLogoutFunc())

	var group *gin.RouterGroup
	if len(c.CommonInfo.User) != 0 {
		//group = engine.Group("/", gin.BasicAuthForRealm(gin.Accounts{
		//	c.CommonInfo.User: c.CommonInfo.Pwd,
		//}, "Restricted"))
		group = engine.Group("/", c.BasicAuth())
	} else {
		group = engine.Group("/")
	}
	group.GET("/", c.MakeIndexFunc())
	group.GET("/tokens", c.MakeQueryTokensFunc())
	group.POST("/add", c.MakeAddTokenFunc())
	group.POST("/update", c.MakeUpdateTokensFunc())
	group.POST("/remove", c.MakeRemoveTokensFunc())
	group.POST("/disable", c.MakeDisableTokensFunc())
	group.POST("/enable", c.MakeEnableTokensFunc())
	group.GET("/proxy/*serverApi", c.MakeProxyFunc())
}
