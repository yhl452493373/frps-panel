package controller

import (
	"fmt"
	plugin "github.com/fatedier/frp/pkg/plugin/server"
	"log"
	"strconv"
	"strings"
)

func (c *HandleController) HandleLogin(content *plugin.LoginContent) plugin.Response {
	token := content.Metas["token"]
	user := content.User
	return c.JudgeToken(user, token)
}

func (c *HandleController) HandleNewProxy(content *plugin.NewProxyContent) plugin.Response {
	token := content.User.Metas["token"]
	user := content.User.User
	judgeToken := c.JudgeToken(user, token)
	if judgeToken.Reject {
		return judgeToken
	}
	return c.JudgePort(content)
}

func (c *HandleController) HandlePing(content *plugin.PingContent) plugin.Response {
	token := content.User.Metas["token"]
	user := content.User.User
	return c.JudgeToken(user, token)
}

func (c *HandleController) HandleNewWorkConn(content *plugin.NewWorkConnContent) plugin.Response {
	token := content.User.Metas["token"]
	user := content.User.User
	return c.JudgeToken(user, token)
}

func (c *HandleController) HandleNewUserConn(content *plugin.NewUserConnContent) plugin.Response {
	token := content.User.Metas["token"]
	user := content.User.User
	return c.JudgeToken(user, token)
}

func (c *HandleController) JudgeToken(user string, token string) plugin.Response {
	var res plugin.Response
	if len(c.Tokens) == 0 {
		res.Unchange = true
	} else if user == "" || token == "" {
		res.Reject = true
		res.RejectReason = "user or meta token can not be empty"
	} else if info, exist := c.Tokens[user]; exist {
		if !info.Enable {
			res.Reject = true
			res.RejectReason = fmt.Sprintf("user [%s] is disabled", user)
		} else {
			if info.Token != token {
				res.Reject = true
				res.RejectReason = fmt.Sprintf("invalid meta token for user [%s]", user)
			} else {
				res.Unchange = true
			}
		}
	} else {
		res.Reject = true
		res.RejectReason = fmt.Sprintf("user [%s] not exist", user)
	}

	return res
}

func (c *HandleController) JudgePort(content *plugin.NewProxyContent) plugin.Response {
	var res plugin.Response
	var portErr error
	var reject = false
	supportProxyTypes := []string{
		"tcp", "tcpmux", "udp", "http", "https",
	}
	proxyType := content.ProxyType
	if stringContains(proxyType, supportProxyTypes) {
		log.Printf("proxy type [%v] not support, plugin do nothing", proxyType)
		res.Unchange = true
		return res
	}

	user := content.User.User
	userPort := content.RemotePort
	userDomains := content.CustomDomains
	userSubdomain := content.SubDomain

	portAllowed := true
	if proxyType == "tcp" || proxyType == "udp" {
		portAllowed = false
		if token, exist := c.Tokens[user]; exist {
			for _, port := range token.Ports {
				if str, ok := port.(string); ok {
					if strings.Contains(str, "-") {
						allowedRanges := strings.Split(str, "-")
						if len(allowedRanges) != 2 {
							portErr = fmt.Errorf("user [%v] port range [%v] format error", user, port)
							break
						}
						start, err := strconv.Atoi(trimString(allowedRanges[0]))
						if err != nil {
							portErr = fmt.Errorf("user [%v] port rang [%v] start port [%v] is not a number", user, port, allowedRanges[0])
							break
						}
						end, err := strconv.Atoi(trimString(allowedRanges[1]))
						if err != nil {
							portErr = fmt.Errorf("user [%v] port rang [%v] end port [%v] is not a number", user, port, allowedRanges[0])
							break
						}
						if max(userPort, start) == userPort && min(userPort, end) == userPort {
							portAllowed = true
							break
						}
					} else {
						allowed, err := strconv.Atoi(str)
						if err != nil {
							portErr = fmt.Errorf("user [%v] allowed port [%v] is not a number", user, port)
						}
						if allowed == userPort {
							portAllowed = true
							break
						}
					}
				} else {
					allowed := port
					if allowed == userPort {
						portAllowed = true
						break
					}
				}

			}
		} else {
			portAllowed = true
		}
	}
	if !portAllowed {
		if portErr == nil {
			portErr = fmt.Errorf("user [%v] port [%v] is not allowed", user, userPort)
		}
		reject = true
	}

	domainAllowed := true
	if proxyType == "http" || proxyType == "https" || proxyType == "tcpmux" {
		if portAllowed {
			if token, exist := c.Tokens[user]; exist {
				for _, userDomain := range userDomains {
					if stringContains(userDomain, token.Domains) {
						domainAllowed = false
						break
					}
				}
			}
			if !domainAllowed {
				portErr = fmt.Errorf("user [%v] domain [%v] is not allowed", user, strings.Join(userDomains, ","))
				reject = true
			}
		}
	}

	subdomainAllowed := true
	if proxyType == "http" || proxyType == "https" {
		subdomainAllowed = false
		if portAllowed && domainAllowed {
			if token, exist := c.Tokens[user]; exist {
				for _, subdomain := range token.Subdomains {
					if subdomain == userSubdomain {
						subdomainAllowed = true
						break
					}
				}
			} else {
				subdomainAllowed = true
			}
			if !subdomainAllowed {
				portErr = fmt.Errorf("user [%v] subdomain [%v] is not allowed", user, userSubdomain)
				reject = true
			}
		}
	}

	if reject {
		res.Reject = true
		res.RejectReason = portErr.Error()
	} else {
		res.Unchange = true
	}
	return res
}
