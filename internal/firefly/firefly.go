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
	"log"
	"net/http"
	"strings"
	"time"
)

// const ExchangeRateByDateTemplate = "https://%s/v1/exchange-rates/by-date/%s"
const ExchangeRateTemplate = "%s/v1/exchange-rates"

type Api struct {
	Config ApiConfig
}

func NewApi(config ApiConfig) *Api {
	return &Api{
		Config: config,
	}
}

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
	log.Printf("Sending exchange rate to %s", endpoint)

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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send exchange rate: %s", resp.Status)
	}

	log.Println("Exchange rate sent successfully")
	return nil
}
