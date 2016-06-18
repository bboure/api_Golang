package phapi

import "fmt"

//ErrorResult base struct to be able to Handle errors.
//ErrorCode is 0 if no error
type ErrorResult struct {
	ErrorMessage string `json:"error"`
	ErrorCode    int    `json:"error_code"`
}

//Error return the formated error
func (err ErrorResult) Error() string {
	if err.ErrorCode == 0 {
		return "<nil>"
	}
	return fmt.Sprintf("[%d] - %s", err.ErrorCode, err.ErrorMessage)
}

// http://apidoc.planethoster.net/index.php?title=Account_information
type AccountInformation struct {
	ErrorResult

	Message          string `json:"message"`
	CreditRemaining  string `json:"credit_remaining"`
	CreditCurrency   string `json:"credit_currency"`
	NumActiveOrders  int    `json:"num_active_orders"`
	NumActiveDomains int    `json:"num_active_domains"`
}

// http://apidoc.planethoster.net/index.php?title=Check_domain_availability
type DomainAvailability struct {
	ErrorResult

	Available            bool   `json:"available"`
	Message              string `json:"message"`
	IsPremium            bool   `json:"is_premium"`
	PremiumRegisterPrice string `json:"premium_register_price"`
	PremiumRenewPrice    string `json:"premium_renew_price"`
}

// http://apidoc.planethoster.net/index.php?title=Domain_information
type DomainInformation struct {
	ErrorResult

	Message                     string   `json:"message"`
	OrderID                     int      `json:"order_id"`
	IsTransfer                  bool     `json:"is_transfer"`
	IsRegistration              bool     `json:"is_registration"`
	RegistrationDate            string   `json:"registration_date"`
	ExpiryDate                  string   `json:"expiry_date"`
	RegistrationStatusInfo      string   `json:"registration_status_info"`
	PurchaseStatus              string   `json:"purchase_status"`
	IDProtection                bool     `json:"id_protection"`
	DomainStatuses              []string `json:"domain_statuses"`
	TransferRequestStatus       string   `json:"transfer_request_status"`
	TransferRequestDeniedReason string   `json:"transfer_request_denied_reason"`
	TransferRequestDeniedAt     string   `json:"transfer_request_denied_at"`
	TransferRequestConfirmedAt  string   `json:"transfer_request_confirmed_at"`
}

// http://apidoc.planethoster.net/index.php?title=Retrieve_contact_details
type DomainContactDetails struct {
	ErrorResult

	Message  string `json:"message"`
	Contacts []struct {
		Name        string `json:"name"`
		CompanyName string `json:"company_name"`
		Addr        struct {
			Address1   string `json:"address1"`
			Address2   string `json:"address2"`
			Address3   string `json:"address3"`
			City       string `json:"city"`
			State      string `json:"state"`
			PostalCode string `json:"postal_code"`
			Country    string `json:"country"`
		} `json:"addr"`
		PhoneNumber    string   `json:"phone_number"`
		Fax            string   `json:"fax"`
		Email          string   `json:"email"`
		ContactType    string   `json:"contact_type"`
		PhoneExtension string   `json:"phone_extension"`
		Statuses       []string `json:"statuses"`
	} `json:"contacts"`
}

// http://apidoc.planethoster.net/index.php?title=Show_nameservers
type NameserversResult struct {
	ErrorResult

	Message     string `json:"message"`
	Nameservers []struct {
		Host string `json:"host"`
	} `json:"nameservers"`
}

// http://apidoc.planethoster.net/index.php?title=Show_PlanetHoster_DNS_records
type DNSRecords struct {
	ErrorResult

	Message string `json:"message"`
	Records []struct {
		Type     string `json:"type"`
		Hostname string `json:"hostname"`
		Address  string `json:"address"`
	} `json:"records"`
}

// http://apidoc.planethoster.net/index.php?title=Show_lock_status
type LockResult struct {
	ErrorResult

	Message  string `json:"message"`
	IsLocked bool   `json:"is_locked"`
}

// http://apidoc.planethoster.net/index.php?title=Test_API_connection
type ConnectionTestResult struct {
	ErrorResult

	Message              string `json:"message"`
	SuccessfulConnection bool   `json:"successful_connection"`
}

// http://apidoc.planethoster.net/index.php?title=Show_TLD_Prices
type TLDPrices struct {
	ErrorResult

	Message string              `json:"message"`
	TLDs    map[string]TLDPrice `json:"tlds"`
}
type TLDPrice struct {
	Register                string `json:"register"`
	Transfer                string `json:"transfer"`
	Renew                   string `json:"renew"`
	TransferRequiresEPPCode bool   `json:"transfer_requires_epp_code"`
	IDProtectionSupported   bool   `json:"id_protection_supported"`
}
