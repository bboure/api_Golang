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

//type Result map[string]interface{}

func makeError(r map[string]interface{}) error {
	if code, ok := r["error_code"]; ok {
		return fmt.Errorf("[%f] - %s", code, r["error"])
	}
	//return fmt.Sprint(r)
	return nil
}

func haveError(r map[string]interface{}) bool {
	i, ok := r["error_code"].(float64)
	return (ok && i != 0)
}

func New(apiUrl, apikey, apiuser string) *API {
	return &API{
		url:     apiUrl,	
		key:     apikey,
		user:    apiuser,
		timeout: 30 * time.Second,
	}
}

type API struct {
	url     string
	key     string
	user    string
	timeout time.Duration
}

type Param struct {
	name  string
	value string
}

//NewDomainData create a new DomainData (for domaine registration) with the minimum datas
func NewDomainData(registrant *ContactDomain, ns1, ns2 string) *DomainData {
	return &DomainData{
		Registrant:        registrant,
		NS1:               ns1,
		NS2:               ns2,
		IDProtection:      false, //should be false by default....
		RegisterIfPremium: false, //should be false by default....
	}
}

//DomainData used to register a domain
type DomainData struct {
	Registrant   *ContactDomain
	IDProtection bool
	NS1          string
	NS2          string

	//optionnal
	NS3               string
	NS4               string
	NS5               string
	RegisterIfPremium bool
	Admin             *ContactDomain
	Tech              *ContactDomain
	Billing           *ContactDomain
	//TODO
	//	addtl_fields (assoc. array), associative array (hash map) of key-value pairs that represent additional fields specific for the TLD that is being registered.

}
type ContactDomain struct {
	FirstName,
	LastName,
	Email,
	CompanyName,
	Address1,
	Address2,
	City,
	PostalCode,
	State,
	CountryCode,
	Phone string
}

func (cdom *ContactDomain) Params(prefix string) []Param {
	return []Param{
		{prefix + "_first_name", cdom.FirstName},
		{prefix + "_last_name", cdom.LastName},
		{prefix + "_email", cdom.Email},
		{prefix + "_company_name", cdom.CompanyName},
		{prefix + "_address1", cdom.Address1},
		{prefix + "_address2", cdom.Address2},
		{prefix + "_city", cdom.City},
		{prefix + "_postal_code", cdom.PostalCode},
		{prefix + "_state", cdom.State},
		{prefix + "_country_code", cdom.CountryCode},
		{prefix + "_phone", cdom.Phone},
	}
}

func (d *DomainData) Valid() bool {
	return (len(d.NS1) > 0 && d.Registrant != nil)
}

func (d *DomainData) Params() []Param {
	ps := make([]Param, 0, 15)
	//probably not the best way to do things...
	ps = append(ps, d.Registrant.Params("registrant")...)
	if d.IDProtection {
		ps = append(ps, Param{"id_protection", "1"})
	} else {
		ps = append(ps, Param{"id_protection", "0"})
	}

	ps = append(ps, Param{"ns1", d.NS1})
	ps = append(ps, Param{"ns2", d.NS2})

	if len(d.NS3) > 0 {
		ps = append(ps, Param{"ns3", d.NS3})
	}
	if len(d.NS4) > 0 {
		ps = append(ps, Param{"ns4", d.NS4})
	}
	if len(d.NS5) > 0 {
		ps = append(ps, Param{"ns5", d.NS5})
	}
	if d.RegisterIfPremium {
		ps = append(ps, Param{"register_if_premium", "1"})
	}
	if d.Admin != nil {
		ps = append(ps, d.Admin.Params("admin")...)
	}
	if d.Tech != nil {
		ps = append(ps, d.Tech.Params("tech")...)
	}
	if d.Billing != nil {
		ps = append(ps, d.Billing.Params("billing")...)
	}

	return ps
}

func (api *API) SetTimeout(t time.Duration) {
	api.timeout = t
}

func (api *API) Prepare(method, path string, params []Param) (*http.Client, *http.Request, error) {
	client := &http.Client{
		Timeout: api.timeout,
		//tmp fix ...
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	ps := url.Values{}
	ps.Add("api_key", api.key)
	ps.Add("api_user", api.user)
	if params != nil {
		for _, p := range params {
			ps.Add(p.name, p.value)
		}
	}

	req, err := http.NewRequest(method, api.url+path, bytes.NewBufferString(ps.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("GOphapi/%f", Version))
	return client, req, err
}

func (api *API) Request(method, path string, v interface{}, params []Param) error {
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
	var result map[string]interface{}

	err := api.Request("GET", "/reseller-api/test-connection", &result, nil)
	if err != nil {
		return err
	}

	b, ok := result["successful_connection"].(bool)
	if haveError(result) || !ok || !b {
		return makeError(result)
	}

	return nil
}

//AccountInfo return reseller account info
func (api *API) AccountInfo() (map[string]interface{}, error) {
	var result map[string]interface{}
	err := api.Request("GET", "/reseller-api/account-info", &result, nil)

	if err != nil {
		return nil, err
	}

	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

//DomainAvailable check if a domain is available for registration
func (api *API) DomainAvailable(sld, tld string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := api.Request("GET", "/reseller-api/check-availability", &result, []Param{{"sld", sld}, {"tld", tld}})
	if err != nil {
		return nil, err
	}

	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

//DomainInfo return informations about the domain
func (api *API) DomainInfo(sld, tld string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := api.Request("GET", "/reseller-api/domain-info", &result, []Param{{"sld", sld}, {"tld", tld}})
	if err != nil {
		return nil, err
	}

	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

//Whois return WHOIS infos
func (api *API) Whois(sld, tld string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := api.Request("GET", "/reseller-api/get-contact-details", &result, []Param{{"sld", sld}, {"tld", tld}})
	if err != nil {
		return nil, err
	}

	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

//Nameservers return nameservers of the domain
func (api *API) Nameservers(sld, tld string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := api.Request("GET", "/reseller-api/get-nameservers", &result, []Param{{"sld", sld}, {"tld", tld}})
	if err != nil {
		return nil, err
	}

	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

//DNSRecords return DNS records of the domain in PlanetHoster's DNS servers
func (api *API) DNSRecords(sld, tld string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := api.Request("GET", "/reseller-api/get-ph-dns-records", &result, []Param{{"sld", sld}, {"tld", tld}})
	if err != nil {
		return nil, err
	}

	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

//LockStatus check if a domain is locked
func (api *API) LockStatus(sld, tld string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := api.Request("GET", "/reseller-api/get-registrar-lock", &result, []Param{{"sld", sld}, {"tld", tld}})
	if err != nil {
		return nil, err
	}

	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

func (api *API) TLDPrices() (map[string]interface{}, error) {
	var result map[string]interface{}
	err := api.Request("GET", "/reseller-api/tld-prices", &result, nil)
	if err != nil {
		return nil, err
	}

	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

//////////////////

//RequestEPPCode request the EPP and unlock the domain
func (api *API) RequestEPPCode(sld, tld string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := api.Request("POST", "/reseller-api/email-epp-code", &result, []Param{{"sld", sld}, {"tld", tld}})
	if err != nil {
		return nil, err
	}

	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

//RegisterDomain register a domain
func (api *API) RegisterDomain(sld, tld string, period int, domain *DomainData) (map[string]interface{}, error) {
	var result map[string]interface{}
	if !domain.Valid() {
		return nil, fmt.Errorf("Invalid DomainData")
	}
	ps := []Param{{"sld", sld}, {"tld", tld}, {"period", fmt.Sprint(period)}}
	ps = append(ps, domain.Params()...)
	err := api.Request("POST", "/reseller-api/register-domain", &result, ps)
	if err != nil {
		return nil, err
	}

	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

//RequestEPPCode request the EPP and unlock the domain
func (api *API) RenewDomain(sld, tld string, period int) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := api.Request("POST", "/reseller-api/renew-domain", &result, []Param{{"sld", sld}, {"tld", tld}, {"period", fmt.Sprint(period)}})
	if err != nil {
		return nil, err
	}
	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

//RequestEPPCode request the EPP and unlock the domain
func (api *API) ChangeContact(sld, tld string, registrant, admin, tech, billing *ContactDomain) (map[string]interface{}, error) {
	var result map[string]interface{}
	ps := []Param{
		{"sld", sld},
		{"tld", tld},
	}
	var contact_types string
	if registrant != nil {
		contact_types += "registrant"
		ps = append(ps, registrant.Params("registrant")...)
	}
	if admin != nil {
		if len(contact_types) > 0 {
			contact_types += ","
		}
		contact_types += "admin"
		ps = append(ps, admin.Params("admin")...)
	}
	if tech != nil {
		if len(contact_types) > 0 {
			contact_types += ","
		}
		contact_types += "tech"
		ps = append(ps, tech.Params("tech")...)
	}
	if billing != nil {
		if len(contact_types) > 0 {
			contact_types += ","
		}
		contact_types += "billing"
		ps = append(ps, billing.Params("billing")...)
	}

	if len(contact_types) == 0 {
		return nil, fmt.Errorf("Nothing to do")
	}

	ps = append(ps, Param{"contact_types", contact_types})
	err := api.Request("POST", "/reseller-api/save-contact-details", &result, ps)

	if err != nil {
		return nil, err
	}

	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

func (api *API) UpdateNameservers(sld, tld, ns1, ns2, ns3, ns4, ns5 string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := api.Request("POST", "/reseller-api/save-nameservers", &result, []Param{
		{"sld", sld},
		{"tld", tld},
		{"ns1", ns1},
		{"ns2", ns2},
		{"ns3", ns3},
		{"ns4", ns4},
		{"ns5", ns5},
	})
	if err != nil {
		return nil, err
	}
	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

type DNSRecord struct {
	Hostname string
	Address  string
	Type     string
}

func (api *API) UpdateDNS(sld, tld string, dns []DNSRecord) (map[string]interface{}, error) {
	var result map[string]interface{}
	ps := []Param{
		{"sld", sld},
		{"tld", tld},
	}

	for i, d := range dns {
		stri := fmt.Sprint(i + 1)
		ps = append(ps,
			Param{"hostname" + stri, d.Hostname},
			Param{"address" + stri, d.Hostname},
			Param{"type" + stri, d.Hostname})
	}
	err := api.Request("POST", "/reseller-api/save-ph-dns-records", &result, ps)
	if err != nil {
		return nil, err
	}
	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

func (api *API) UpdateDomainLock(sld, tld, lockaction string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := api.Request("POST", "/reseller-api/save-registrar-lock", &result, []Param{{"sld", sld}, {"tld", tld}, {"lock_action", lockaction}})
	if err != nil {
		return nil, err
	}
	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}

func (api *API) DeleteDNS(sld, tld string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := api.Request("POST", "/reseller-api/delete-ph-dns-zone", &result, []Param{{"sld", sld}, {"tld", tld}})
	if err != nil {
		return nil, err
	}
	if haveError(result) {
		return nil, makeError(result)
	}

	return result, nil
}
