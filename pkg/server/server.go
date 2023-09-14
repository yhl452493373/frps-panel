package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"frps-panel/pkg/server/controller"
	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Server struct {
	cfg     controller.HandleController
	s       *http.Server
	tls     TLS
	done    chan struct{}
	rootDir string
}

type TLS struct {
	Enable   bool
	Cert     string
	Key      string
	Protocol string
}

func New(rootDir string, cfg controller.HandleController, tls TLS) (*Server, error) {
	s := &Server{
		cfg:     cfg,
		done:    make(chan struct{}),
		rootDir: rootDir,
		tls:     tls,
	}

	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Run() error {
	bindAddress := s.cfg.CommonInfo.PluginAddr + ":" + strconv.Itoa(s.cfg.CommonInfo.PluginPort)
	l, err := net.Listen("tcp", bindAddress)
	if err != nil {
		return err
	}
	log.Printf("%s server listen on %s", s.tls.Protocol, l.Addr().String())
	go func() {
		if s.tls.Enable {
			configDir := filepath.Dir(s.cfg.ConfigFile)

			cert := filepath.Join(configDir, s.tls.Cert)
			_, err := os.Stat(cert)
			if err != nil && !os.IsExist(err) {
				cert = s.tls.Cert
			}

			key := filepath.Join(configDir, s.tls.Key)
			_, err = os.Stat(key)
			if err != nil && !os.IsExist(err) {
				key = s.tls.Key
			}

			if err = s.s.ServeTLS(l, cert, key); !errors.Is(http.ErrServerClosed, err) {
				log.Printf("error shutdown %s server: %v", s.tls.Protocol, err)
				_ = s.Stop()
			}
		} else {
			if err = s.s.Serve(l); !errors.Is(http.ErrServerClosed, err) {
				log.Printf("error shutdown %s server: %v", s.tls.Protocol, err)
				_ = s.Stop()
			}
		}
	}()
	<-s.done
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.s.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown %s server error: %v", s.tls.Protocol, err)
	}
	log.Printf("%s server exited", s.tls.Protocol)

	close(s.done)
	return nil
}

func (s *Server) init() error {
	if err := s.initHTTPServer(); err != nil {
		log.Printf("init %s server error: %v", s.tls.Protocol, err)
		return err
	}
	return nil
}

func LoadSupportLanguage(dir string) ([]language.Tag, error) {
	var tags []language.Tag

	files, err := os.Open(dir)

	fileList, err := files.Readdir(-1)
	if err != nil {
		log.Printf("error read lang directory: %v", err)
		return tags, err
	}

	err = files.Close()
	if err != nil {
		return nil, err
	}

	for _, file := range fileList {
		name, _ := strings.CutSuffix(file.Name(), ".json")
		parsedLang, _ := language.Parse(name)
		tags = append(tags, parsedLang)
	}

	if len(tags) == 0 {
		return tags, fmt.Errorf("not found any language file in directory: %v", dir)
	}

	return tags, nil
}

func GinI18nLocalize(rootDir string) gin.HandlerFunc {
	assets := filepath.Join(rootDir, "assets")
	_, err := os.Stat(assets)
	if err != nil && !os.IsExist(err) {
		assets = "./assets"
	}
	lang := filepath.Join(assets, "lang")
	tags, err := LoadSupportLanguage(lang)
	if err != nil {
		log.Panicf("language file is not found: %v", err)
	}

	return ginI18n.Localize(
		ginI18n.WithBundle(&ginI18n.BundleCfg{
			RootPath:         lang,
			AcceptLanguage:   tags,
			DefaultLanguage:  language.Chinese,
			FormatBundleFile: "json",
			UnmarshalFunc:    json.Unmarshal,
		}),
		ginI18n.WithGetLngHandle(
			func(context *gin.Context, defaultLng string) string {
				header := context.GetHeader("Accept-Language")
				lang, _, err := language.ParseAcceptLanguage(header)
				if err != nil {
					return defaultLng
				}
				return lang[0].String()
			},
		),
	)
}

func (s *Server) initHTTPServer() error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	authStore := cookie.NewStore([]byte("frps-panel"))
	authStore.Options(sessions.Options{
		Secure:   true,
		HttpOnly: false,
		SameSite: 4,
		Path:     "/",
		MaxAge:   s.cfg.CommonInfo.AdminKeepTime,
	})
	engine.Use(sessions.Sessions(controller.SessionName, authStore))
	engine.Use(GinI18nLocalize(s.rootDir))
	s.s = &http.Server{
		Handler: engine,
	}
	controller.NewHandleController(&s.cfg).Register(s.rootDir, engine)
	return nil
}
