package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (c *HandleController) BasicAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		if strings.TrimSpace(c.CommonInfo.User) == "" || strings.TrimSpace(c.CommonInfo.Pwd) == "" {
			ClearLogin(context)
			if context.Request.RequestURI == LoginUrl {
				context.Redirect(http.StatusTemporaryRedirect, LoginSuccessUrl)
			}
			return
		}

		auth, err := context.Request.Cookie("token")

		if err == nil {
			username, password, _ := ParseBasicAuth(auth.Value)

			usernameMatch := username == c.CommonInfo.User
			passwordMatch := password == c.CommonInfo.Pwd

			if usernameMatch && passwordMatch {
				context.Next()
				return
			}
		}

		isAjax := context.GetHeader("X-Requested-With") == "XMLHttpRequest"

		if !isAjax && context.Request.RequestURI != LoginUrl {
			context.Redirect(http.StatusTemporaryRedirect, LoginUrl)
		} else {
			context.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func ParseBasicAuth(auth string) (username, password string, ok bool) {
	if len(auth) < len(AuthPrefix) || auth[:len(AuthPrefix)] != AuthPrefix {
		return "", "", false
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(AuthPrefix):])
	if err != nil {
		return "", "", false
	}
	cs := string(c)
	username, password, ok = strings.Cut(cs, ":")
	if !ok {
		return "", "", false
	}
	return username, password, true
}

func EncodeBasicAuth(username, password string) string {
	authString := fmt.Sprintf("%s:%s", username, password)
	return AuthPrefix + base64.StdEncoding.EncodeToString([]byte(authString))
}

func ClearLogin(context *gin.Context) {
	context.SetCookie("token", "", -1, "/", context.Request.Host, false, false)
}
