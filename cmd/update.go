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
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"ffiii-rate-updater/internal/exchange"
	"ffiii-rate-updater/internal/firefly"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Fetch and update exchange rates in Firefly III",
	Long:  `Fetch exchange rates for specified currencies and update them in Firefly III.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		currencies := viper.GetStringSlice("currencies")

		if len(currencies) < 2 {
			return fmt.Errorf("please provide at least two currencies to fetch exchange rates")
		}

		apiKey := viper.GetString("firefly.api_key")
		if apiKey == "" {
			return fmt.Errorf("firefly API key is not set")
		}

		apiUrl := viper.GetString("firefly.api_url")
		if apiUrl == "" {
			return fmt.Errorf("firefly API URL is not set")
		}

		exchangeApi, err := exchange.NewApi(currencies, viper.GetString("date"))
		if err != nil {
			return fmt.Errorf("failed to initialize exchange API: %v", err)
		}

		fireflyApi := firefly.NewApi(firefly.ApiConfig{
			ApiKey:         apiKey,
			ApiUrl:         apiUrl,
			TimeoutSeconds: 10,
		})

		// Send exchange rates as batch
		for i := range currencies {
			fromCurrency := currencies[i]
			rates := make(map[string]float64)
			date := ""

			for j := range currencies {
				if i != j {
					toCurrency := currencies[j]
					rate, err := exchangeApi.GetRate(fromCurrency, toCurrency)
					if err != nil {
						log.Printf("Error fetching rate for %s/%s: %v", fromCurrency, toCurrency, err)
						continue
					}
					rates[toCurrency] = rate.Value
					// if not set yet, set the date
					if date == "" {
						date = rate.Date
					}
				}
			}

			err = fireflyApi.SendExchangeRateByDate(fromCurrency, rates, date)
			if err != nil {
				log.Printf("Error sending batch rates for %s: %v", fromCurrency, err)
				break
			}
			log.Printf("Sent batch exchange rates for %s on %s", fromCurrency, date)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
