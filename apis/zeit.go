package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Zeit struct {
	Team       string
	Token      string
	HttpClient *http.Client
}

const SecretsUrl = "https://api.zeit.co/v2/now/secrets"

func (client *Zeit) getSecretsUrl() string {
	res := SecretsUrl

	if client.Team != "" {
		res += "?teamId=" + client.Team
	}

	return res
}

func (client *Zeit) SetSecret(secretName string, secretValue string) error {
	httpClient := client.HttpClient

	if httpClient == nil {
		httpClient = &http.Client{}
	}

	secretUrl := client.getSecretsUrl()
	request, err := http.NewRequest("DELETE", secretUrl, nil)
	if err != nil {
		return err
	}

	httpClient.Do(request)

	requestBody, err := json.Marshal(map[string]interface{}{
		"name":  secretName,
		"value": secretValue,
	})
	if err != nil {
		return err
	}

	request, err = http.NewRequest("POST", secretUrl, bytes.NewReader(requestBody))
	if err != nil {
		return err
	}

	res, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Could not set secret.")
	}

	return nil
}
