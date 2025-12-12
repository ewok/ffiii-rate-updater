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

import "fmt"

type ApiConfig struct {
	URL            string
	FallbackURL    string
	TimeoutSeconds int
}

func GetApiConfig() ApiConfig {
	return ApiConfig{
		URL: "https://cdn.jsdelivr.net/npm/@fawazahmed0/currency-api@%s/v1/%s/%s.min.json",
		// FallbackURL: "https://%s.currency-api.pages.dev/v1/%s/%s.min.json",
		TimeoutSeconds: 10,
	}
}

func (apiconfig *ApiConfig) GetURL(date string, currency string, endpoint string) string {
	return fmt.Sprintf(apiconfig.URL, date, endpoint, currency)
}

// func (apiconfig *ApiConfig) GetFallbackURL(date string, currency string, endpoint string) string {
// 	return fmt.Sprintf(apiconfig.FallbackURL, date, endpoint, currency)
// }
