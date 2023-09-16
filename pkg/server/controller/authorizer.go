package controller

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func (c *HandleController) BasicAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		if trimString(c.CommonInfo.AdminUser) == "" || trimString(c.CommonInfo.AdminPwd) == "" {
			if context.Request.RequestURI == LoginUrl {
				context.Redirect(http.StatusTemporaryRedirect, LoginSuccessUrl)
			}
			return
		}

		session := sessions.Default(context)
		auth := session.Get(AuthName)

		if auth != nil {
			if c.CommonInfo.AdminKeepTime > 0 {
				cookie, _ := context.Request.Cookie(SessionName)
				if cookie != nil {
					//important thx https://blog.csdn.net/zhanghongxia8285/article/details/107321838/
					cookie.Expires = time.Now().Add(time.Second * time.Duration(c.CommonInfo.AdminKeepTime))
					http.SetCookie(context.Writer, cookie)
				}
			}

			username, password, _ := parseBasicAuth(fmt.Sprintf("%v", auth))

			usernameMatch := username == c.CommonInfo.AdminUser
			passwordMatch := password == c.CommonInfo.AdminPwd

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

func (c *HandleController) LoginAuth(username, password string, context *gin.Context) bool {
	if trimString(c.CommonInfo.AdminUser) == "" || trimString(c.CommonInfo.AdminPwd) == "" {
		return true
	}

	session := sessions.Default(context)

	sessionAuth := session.Get(AuthName)
	internalAuth := encodeBasicAuth(c.CommonInfo.AdminUser, c.CommonInfo.AdminPwd)

	if sessionAuth == internalAuth {
		return true
	} else {
		basicAuth := encodeBasicAuth(username, password)
		if basicAuth == internalAuth {
			session.Set(AuthName, basicAuth)
			_ = session.Save()
			return true
		} else {
			session.Delete(AuthName)
			_ = session.Save()
			return false
		}
	}
}

func ClearAuth(context *gin.Context) {
	session := sessions.Default(context)
	session.Delete(AuthName)
	_ = session.Save()
}

func parseBasicAuth(auth string) (username, password string, ok bool) {
	c, err := base64.StdEncoding.DecodeString(auth)
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

func encodeBasicAuth(username, password string) string {
	authString := fmt.Sprintf("%s:%s", username, password)
	return base64.StdEncoding.EncodeToString([]byte(authString))
}
