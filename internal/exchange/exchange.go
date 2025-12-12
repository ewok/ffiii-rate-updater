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
package exchange

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Api struct {
	Config ApiConfig
	Rates  []Rate
}

type ApiResponse struct {
	Date  string
	Rates map[string]float64
}

func NewApi(rawCurrencies []string, date string) *Api {

	api := Api{
		Config: GetApiConfig(),
	}

	// Convert rawCurrencies to []Currency
	var exCurrencies []Currency
	for _, curr := range rawCurrencies {
		exCurrencies = append(exCurrencies, NewCurrency(curr))
	}

	// Initialize exchange rates
	rates, err := api.getExchangeRates(exCurrencies, date)
	if err != nil {
		log.Fatalf("Error initializing API rates: %v", err)
	}

	api.Rates = rates

	return &api
}

func (api *Api) GetRate(from string, to string) (Rate, error) {

	var rate Rate

	pair := Pair{From: NewCurrency(from), To: NewCurrency(to)}
	for _, r := range api.Rates {
		if r.Pair == pair {
			rate = r
			break
		}
	}

	if rate.Pair.From == NewCurrency("") || rate.Pair.To == NewCurrency("") {
		return Rate{}, fmt.Errorf("rate not found for pair %s/%s", from, to)
	}

	return rate, nil
}

func (api *Api) getExchangeRates(currencies []Currency, date string) ([]Rate, error) {

	rates := []Rate{}

	// fetch exchange rates for the given currencies
	for _, currency := range currencies {
		resp, err := api.fetchRates(currency.GetLCode(), date)
		if err != nil {
			return nil, err
		}

		for k, v := range resp.Rates {

			rates = append(rates, Rate{
				Date: resp.Date,
				Pair: Pair{
					From: currency,
					To:   NewCurrency(k),
				},
				Value: v,
			})
		}

	}

	return rates, nil
}

func (api *Api) fetchRates(currency string, date string) (response ApiResponse, err error) {

	if date == "" {
		date = "latest"
	}

	// TODO: Implement fallback mechanism
	url := api.Config.GetURL(date, currency, "currencies")

	log.Printf("Fetching rates for %s on %s", currency, date)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ApiResponse{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ApiResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ApiResponse{}, fmt.Errorf("failed to fetch rates: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ApiResponse{}, err
	}

	var rawJson map[string]interface{}
	err = json.Unmarshal(body, &rawJson)
	if err != nil {
		return ApiResponse{}, err
	}

	// Extract date and rates
	// Safely extract "date" as string
	rawDate, ok := rawJson["date"]
	if !ok {
		return ApiResponse{}, fmt.Errorf("missing 'date' field in API response")
	}
	date, ok = rawDate.(string)
	if !ok {
		return ApiResponse{}, fmt.Errorf("'date' field is not a string in API response")
	}

	// Safely extract rates map
	rawRates, ok := rawJson[currency]
	if !ok {
		return ApiResponse{}, fmt.Errorf("missing '%s' field in API response", currency)
	}
	ratesMap, ok := rawRates.(map[string]any)
	if !ok {
		return ApiResponse{}, fmt.Errorf("'%s' field is not a map in API response", currency)
	}
	var rates = make(map[string]float64)
	for key, value := range ratesMap {
		floatVal, ok := value.(float64)
		if !ok {
			return ApiResponse{}, fmt.Errorf("rate for '%s' is not a float64 in API response", key)
		}
		rates[key] = floatVal
	}

	return ApiResponse{
		Date:  date,
		Rates: rates,
	}, nil
}
