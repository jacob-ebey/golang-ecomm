package apis

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Avatax struct {
	BearerToken string
	Username    string
	Password    string
	AccountID   string
	LicenseKey  string
	Development bool
	HttpClient  *http.Client
}

type AvataxAddress struct {
	Line1      string
	Line2      string
	Line3      string
	City       string
	Region     string
	PostalCode string
	Country    string
}

type AvataxRates struct {
	TotalRate float64      `json:"totalRate"`
	Rates     []AvataxRate `json:"rates"`
	Error     *AvataxError
}

type AvataxError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AvataxRate struct {
	Rate float64 `json:"rate"`
	Name string  `json:"name"`
	Type string  `json:"type"`
}

func (err *AvataxError) Error() string {
	if err.Message != "" {
		return err.Message
	}

	if err.Code != "" {
		return err.Code
	}

	return "Unknown error."
}

func (e *AvataxError) Extensions() map[string]interface{} {
	var code *string
	if e.Code != "" {
		code = &e.Code
	}

	return map[string]interface{}{
		"code": code,
	}
}

func (config *Avatax) GetUrl() string {
	if config.Development {
		return "https://sandbox-rest.avatax.com/api/v2/taxrates/byaddress"
	}

	return "https://rest.avatax.com/api/v2/taxrates/byaddress"
}

func createBasicAuthorization(a string, b string) string {
	return "Basic " + b64.URLEncoding.EncodeToString([]byte(a+":"+b))
}

func (config *Avatax) getAuthorization() (string, error) {
	if config.BearerToken != "" {
		return "Bearer " + config.BearerToken, nil
	}

	if config.Username != "" && config.Password != "" {
		return createBasicAuthorization(config.Username, config.Password), nil
	}

	if config.AccountID != "" && config.LicenseKey != "" {
		return createBasicAuthorization(config.AccountID, config.LicenseKey), nil
	}

	return "", fmt.Errorf("No valid authorization configuration provided.")
}

func (config *Avatax) TaxRatesByAddress(address AvataxAddress) (*AvataxRates, error) {
	authorization, err := config.getAuthorization()
	if err != nil {
		return nil, err
	}

	httpClient := config.HttpClient

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	url := config.GetUrl()
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", authorization)

	query := request.URL.Query()
	query.Add("line1", address.Line1)
	query.Add("line2", address.Line2)
	query.Add("line3", address.Line3)
	query.Add("city", address.City)
	query.Add("region", address.Region)
	query.Add("postalCode", address.PostalCode)
	query.Add("country", address.Country)
	request.URL.RawQuery = query.Encode()

	res, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	read, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	body := AvataxRates{}
	if err := json.Unmarshal(read, &body); err != nil || body.Error != nil {
		if err != nil {
			return nil, err
		}
		return nil, body.Error
	}

	return &body, nil
}
