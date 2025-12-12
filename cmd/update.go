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

var date string

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Fetch and update exchange rates in Firefly III",
	Long:  `Fetch exchange rates for specified currencies and update them in Firefly III.`,
	Run: func(cmd *cobra.Command, args []string) {

		currencies := viper.GetStringSlice("currencies")
		if len(currencies) < 2 {
			fmt.Println("Please provide at least two currencies to fetch exchange rates.")
			return
		}
		apiKey := viper.GetString("firefly.api_key")
		if apiKey == "" {
			fmt.Println("Firefly API key is not set. Please set it in the configuration.")
			return
		}
		apiUrl := viper.GetString("firefly.api_url")
		if apiUrl == "" {
			fmt.Println("Firefly API URL is not set. Please set it in the configuration.")
			return
		}

		exchangeApi := exchange.NewApi(currencies, date)
		log.Printf("Fetching exchange rate for %s/%s", currencies[0], currencies[1])

		fireflyApi := firefly.NewApi(firefly.ApiConfig{
			ApiKey: apiKey,
			ApiUrl: apiUrl,
		})

		// Send all exchange variants between the provided currencies
		for i := range currencies {
			for j := range currencies {
				if i != j {
					rate, err := exchangeApi.GetRate(currencies[i], currencies[j])
					if err != nil {
						log.Printf("Error fetching rate for %s/%s: %v", currencies[i], currencies[j], err)
						continue
					}

					err = fireflyApi.SendExchangeRates(rate.Value, currencies[i], currencies[j], rate.Date)
					if err != nil {
						log.Printf("Error sending rate for %s/%s: %v", currencies[i], currencies[j], err)
					}
					log.Printf("Sent exchange rate for %s/%s: %.6f on %s", currencies[i], currencies[j], rate.Value, rate.Date)
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVarP(&date, "date", "d", "", "Date for the exchange rates in YYYY-MM-DD format (default is latest available)")
}
