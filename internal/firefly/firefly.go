/*
Copyright © 2025 Artur Taranchiev <artur.taranchiev@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package firefly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// const ExchangeRateByDateTemplate = "https://%s/v1/exchange-rates/by-date/%s"
const ExchangeRateTemplate = "%s/exchange-rates"

// ApiConfig holds configuration for the Firefly III API.
type Api struct {
	// Config contains the API configuration details.
	Config ApiConfig
}

// NewApi creates a new Api instance with the provided configuration.
// Parameters:
//   - config: an ApiConfig struct containing the API configuration details.
//
// Returns:
//   - A pointer to an Api struct initialized with the provided configuration.
func NewApi(config ApiConfig) *Api {
	return &Api{
		Config: config,
	}
}

// SendExchangeRates sends the exchange rate data to the Firefly API.
//
// Parameters:
//   - rate: the exchange rate value to be sent.
//   - fromCurrency: the source currency code (e.g., "USD").
//   - toCurrency: the target currency code (e.g., "EUR").
//   - date: the date for which the exchange rate is applicable (in "YYYY-MM-DD" format). If empty, the current date is used.
//
// Returns:
//   - An error if the operation fails; otherwise, nil.
func (api *Api) SendExchangeRates(rate float64, fromCurrency string, toCurrency string, date string) error {

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	payload := map[string]string{
		"date": date,
		"from": strings.ToUpper(fromCurrency),
		"to":   strings.ToUpper(toCurrency),
		"rate": fmt.Sprintf("%.6f", rate),
	}

	endpoint := fmt.Sprintf(ExchangeRateTemplate, api.Config.ApiUrl)

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.Config.ApiKey))

	client := &http.Client{Timeout: time.Duration(api.Config.TimeoutSeconds) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send exchange rate: %d", resp.StatusCode)
	}

	return nil
}
