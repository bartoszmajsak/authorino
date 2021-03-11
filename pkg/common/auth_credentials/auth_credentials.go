package auth_credentials

import (
	"fmt"
	"regexp"
	"strings"

	envoyServiceAuthV3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	ctrl "sigs.k8s.io/controller-runtime"
)

type AuthCredentials interface {
	GetCredentialsFromReq(*envoyServiceAuthV3.AttributeContext_HttpRequest) (string, error)
}

type AuthCredential struct {
	KeySelector string `yaml:"key_selector"`
	In          string `yaml:"in"`
}

const (
	inCustomHeader = "custom_header"
	inAuthHeader   = "authorization_header"
	inQuery        = "query"
)

var (
	authCredLog = ctrl.Log.WithName("Authorino").WithName("AuthCredential")
	notFoundErr = fmt.Errorf("credential not found")
)

func (c *AuthCredential) GetCredentialsFromReq(httpReq *envoyServiceAuthV3.AttributeContext_HttpRequest) (string, error) {
	switch c.In {
	case inCustomHeader:
		return getCredFromCustomHeader(httpReq.GetHeaders(), c.KeySelector)
	case inAuthHeader:
		return getCredFromAuthHeader(httpReq.GetHeaders(), c.KeySelector)
	case inQuery:
		return getCredFromQuery(httpReq.GetPath(), c.KeySelector)
	default:
		return "", fmt.Errorf("the credential location is not supported")
	}
}

func getCredFromCustomHeader(headers map[string]string, keyName string) (string, error) {
	cred, ok := headers[keyName]
	if !ok {
		authCredLog.Error(notFoundErr, "the credential was not found in the request header")
		return "", notFoundErr
	}
	return cred, nil
}
func getCredFromAuthHeader(headers map[string]string, keyName string) (string, error) {
	authHeader, ok := headers["authorization"]

	if !ok {
		authCredLog.Error(notFoundErr, "the Authorization header is not set")
		return "", notFoundErr
	}
	prefix := keyName + " "
	if strings.HasPrefix(authHeader, prefix) {
		return strings.TrimPrefix(authHeader, prefix), nil
	}
	return "", notFoundErr
}

func getCredFromQuery(path string, keyName string) (string, error) {
	const credValue = "credValue"
	regex := regexp.MustCompile("([?&]" + keyName + "=)(?P<" + credValue + ">[^&]*)")
	matches := regex.FindStringSubmatch(path)
	if len(matches) == 0 {
		return "", notFoundErr
	}
	return matches[regex.SubexpIndex(credValue)], nil
}