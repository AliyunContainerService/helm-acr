// Modifications copyright (C) 2019 Alibaba Group Holding Limited / Yuning Xie (xyn1016@gmail.com)

package chartmuseum

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// DownloadFile downloads a file from ChartMuseum
func (client *Client) DownloadFile(filePath string) (*http.Response, error) {
	u, err := url.Parse(client.opts.url)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(client.opts.contextPath, strings.TrimPrefix(u.Path, client.opts.contextPath), filePath)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	accessToken := client.opts.accessToken

	if client.opts.autoTokenAuth {
		resp, err := client.Do(req)
		if err != nil {
			return resp, err
		} else if resp.StatusCode == http.StatusUnauthorized {
			token, err := client.GetAuthTokenFromResponse(resp)
			if err != nil {
				return nil, err
			}
			accessToken = token
		} else {
			return resp, err
		}
	}

	if accessToken != "" {
		if client.opts.authHeader != "" {
			req.Header.Set(client.opts.authHeader, client.opts.accessToken)
		} else {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		}
	} else if client.opts.username != "" && client.opts.password != "" {
		req.SetBasicAuth(client.opts.username, client.opts.password)
	}

	return client.Do(req)
}

func (client *Client) GetAuthTokenFromResponse(resp *http.Response) (string, error) {
	authHeader := resp.Header.Get("Www-Authenticate")
	authHeader = strings.Split(authHeader, " ")[1]
	tokens := strings.Split(authHeader, ",")
	var realm, service, scope string
	for _, token := range tokens {
		if strings.HasPrefix(token, "realm") {
			realm = strings.Trim(token[len("realm="):], "\"")
		}
		if strings.HasPrefix(token, "service") {
			service = strings.Trim(token[len("service="):], "\"")
		}
		if strings.HasPrefix(token, "scope") {
			scope = strings.Trim(token[len("scope="):], "\"")
		}
	}
	if realm == "" {
		return "", fmt.Errorf("missing realm in bearer auth challenge")
	}
	if service == "" {
		return "", fmt.Errorf("missing service in bearer auth challenge")
	}
	if scope == "" {
		return "", fmt.Errorf("missing scope in bearer auth challenge")
	}
	return client.getBearerToken(realm, service, scope)
}

func (client *Client) getBearerToken(realm, service, scope string) (string, error) {
	authReq, err := http.NewRequest("POST", realm, nil)
	if err != nil {
		return "", err
	}
	getParams := authReq.URL.Query()
	getParams.Add("service", service)
	if scope != "" {
		getParams.Add("scope", scope)
	}
	authReq.URL.RawQuery = getParams.Encode()
	if client.opts.username != "" && client.opts.password != "" {
		authReq.SetBasicAuth(client.opts.username, client.opts.password)
	}
	resp, err := client.Do(authReq)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return "", fmt.Errorf("unable to retrieve auth token: 401 unauthorized")
	case http.StatusOK:
		break
	default:
		return "", fmt.Errorf("unexpected http code: %d, URL: %s", resp.StatusCode, authReq.URL)
	}
	tokenBlob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	token := struct {
		Token string `json:"access_token"`
	}{}
	if err := json.Unmarshal(tokenBlob, &token); err != nil {
		return "", err
	}
	return token.Token, nil
}
