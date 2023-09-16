package controller

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"strings"
)

func filter(main TokenInfo, sub TokenInfo) bool {
	replaceSpaceUser := trimAllSpace.ReplaceAllString(sub.User, "")
	if len(replaceSpaceUser) != 0 {
		if !strings.Contains(main.User, replaceSpaceUser) {
			return false
		}
	}

	replaceSpaceToken := trimAllSpace.ReplaceAllString(sub.Token, "")
	if len(replaceSpaceToken) != 0 {
		if !strings.Contains(main.Token, replaceSpaceToken) {
			return false
		}
	}

	replaceSpaceComment := trimAllSpace.ReplaceAllString(sub.Comment, "")
	if len(replaceSpaceComment) != 0 {
		if !strings.Contains(main.Comment, replaceSpaceComment) {
			return false
		}
	}
	return true
}

func trimString(str string) string {
	return strings.TrimSpace(str)
}

func (c *HandleController) verifyToken(token TokenInfo, operate int) OperationResponse {
	response := OperationResponse{
		Success: true,
		Code:    Success,
		Message: "operate success",
	}

	var (
		validateExist      = false
		validateNotExist   = false
		validateUser       = false
		validateToken      = false
		validateComment    = false
		validatePorts      = false
		validateDomains    = false
		validateSubdomains = false
	)

	if operate == TOKEN_ADD {
		validateExist = true
		validateUser = true
		validateToken = true
		validateComment = true
		validatePorts = true
		validateDomains = true
		validateSubdomains = true
	} else if operate == TOKEN_UPDATE {
		validateNotExist = true
		validateUser = true
		validateToken = true
		validateComment = true
		validatePorts = true
		validateDomains = true
		validateSubdomains = true
	} else if operate == TOKEN_ENABLE || operate == TOKEN_DISABLE || operate == TOKEN_REMOVE {
		validateNotExist = true
	}

	if validateUser && !userFormat.MatchString(token.User) {
		response.Success = false
		response.Code = UserFormatError
		response.Message = fmt.Sprintf("operate failed, user [%s] format error", token.User)
		log.Printf(response.Message)
		return response
	}

	if validateExist {
		if _, exist := c.Tokens[token.User]; exist {
			response.Success = false
			response.Code = UserExist
			response.Message = fmt.Sprintf("operate failed, user [%s] exist ", token.User)
			log.Printf(response.Message)
			return response
		}
	}

	if validateNotExist {
		if _, exist := c.Tokens[token.User]; !exist {
			response.Success = false
			response.Code = UserNotExist
			response.Message = fmt.Sprintf("operate failed, user [%s] not exist ", token.User)
			log.Printf(response.Message)
			return response
		}
	}

	if validateToken && !tokenFormat.MatchString(token.Token) {
		response.Success = false
		response.Code = TokenFormatError
		response.Message = fmt.Sprintf("operate failed, token [%s] format error", token.Token)
		log.Printf(response.Message)
		return response
	}

	trimmedComment := trimString(token.Comment)
	if validateComment && trimmedComment != "" && commentFormat.MatchString(trimmedComment) {
		response.Success = false
		response.Code = CommentFormatError
		response.Message = fmt.Sprintf("operate failed, comment [%s] format error", token.Comment)
		log.Printf(response.Message)
		return response
	}

	if validatePorts {
		for _, port := range token.Ports {
			trimmedPort := trimString(port)
			if trimmedPort != "" && !portsFormatSingle.MatchString(trimmedPort) && !portsFormatRange.MatchString(trimmedPort) {
				response.Success = false
				response.Code = PortsFormatError
				response.Message = fmt.Sprintf("operate failed, ports [%v] format error", token.Ports)
				log.Printf(response.Message)
				return response
			}
		}
	}

	if validateDomains {
		for _, domain := range token.Domains {
			trimmedDomain := trimString(domain)
			if trimmedDomain != "" && !domainFormat.MatchString(trimmedDomain) {
				response.Success = false
				response.Code = DomainsFormatError
				response.Message = fmt.Sprintf("operate failed, domains [%v] format error", token.Domains)
				log.Printf(response.Message)
				return response
			}
		}
	}

	if validateSubdomains {
		for _, subdomain := range token.Subdomains {
			trimmedSubdomain := trimString(subdomain)
			if trimmedSubdomain != "" && !subdomainFormat.MatchString(trimmedSubdomain) {
				response.Success = false
				response.Code = SubdomainsFormatError
				response.Message = fmt.Sprintf("operate failed, subdomains [%v] format error", token.Subdomains)
				log.Printf(response.Message)
				return response
			}
		}
	}

	return response
}

func cleanStrings(originalStrings []string) []string {
	cleanedStrings := make([]string, len(originalStrings))
	for i, str := range originalStrings {
		cleanedStrings[i] = cleanString(str)
	}
	return cleanedStrings
}

func cleanString(originalString string) string {
	return trimAllSpace.ReplaceAllString(originalString, "")
}

func stringContains(element string, data []string) bool {
	for _, v := range data {
		if element == v {
			return true
		}
	}
	return false
}

func tokensList(tokens map[string]TokenInfo) Tokens {
	return Tokens{
		tokens,
	}
}

func (c *HandleController) saveToken() error {
	tokenFile, err := os.Create(c.TokensFile)
	if err != nil {
		log.Printf("error to crate file %v: %v", c.TokensFile, err)
	}

	if err = toml.NewEncoder(tokenFile).Encode(tokensList(c.Tokens)); err != nil {
		log.Printf("error to encode tokens: %v", err)
	}
	if err = tokenFile.Close(); err != nil {
		log.Printf("error to close file %v: %v", c.TokensFile, err)
	}

	return err
}
