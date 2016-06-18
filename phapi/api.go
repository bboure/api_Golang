package phapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const APIUrl = "https://api.planethoster.net"
const Version = 0.1

func New(apikey, apiuser string) *API {
	return &API{
		key:     apikey,
		user:    apiuser,
		timeout: 30 * time.Second,
	}
}

type API struct {
	key     string
	user    string
	timeout time.Duration
}

func (api *API) SetTimeout(t time.Duration) {
	api.timeout = t
}

func (api *API) Prepare(method, path string, params map[string]string) (*http.Client, *http.Request, error) {
	client := &http.Client{
		Timeout: api.timeout,
		//tmp fix ...
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	ps := url.Values{}
	ps.Add("api_key", api.key)
	ps.Add("api_user", api.user)
	if params != nil {
		for k, v := range params {
			ps.Add(k, v)
		}
	}
	req, err := http.NewRequest(method, APIUrl+path, bytes.NewBufferString(ps.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("GOphapi/%f", Version))
	return client, req, err
}

func (api *API) Request(method, path string, params map[string]string, v interface{}) error {
	client, req, err := api.Prepare(method, path, params)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, v)
	if err != nil {
		return err
	}

	return nil
}

//Test test the connection to the API. return nil on successful connection
func (api *API) Test() error {
	var result ConnectionTestResult
	err := api.Request("GET", "/reseller-api/test-connection", nil, &result)
	if err != nil {
		return err
	}

	if result.ErrorCode != 0 || !result.SuccessfulConnection {
		return result.ErrorResult
	}

	return nil
}
