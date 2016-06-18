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

func (api *API) Prepare(actioname, path string, params map[string]string) (*http.Client, *http.Request, error) {
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
	req, err := http.NewRequest(actioname, APIUrl+path, bytes.NewBufferString(ps.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("GOphapi/%f", Version))
	return client, req, err
}

func (api *API) SetTimeout(t time.Duration) {
	api.timeout = t
}

func request(client *http.Client, req *http.Request, v interface{}) error {
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
	client, req, err := api.Prepare("GET", "/reseller-api/test-connection", nil)
	if err != nil {
		return err
	}

	var result ConnectionTestResult
	err = request(client, req, &result)

	if result.ErrorCode != 0 || !result.SuccessfulConnection {
		return result.ErrorResult
	}

	return nil
}
