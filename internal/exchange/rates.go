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

type Pair struct {
	From Currency
	To   Currency
}

type Rate struct {
	Date  string
	Pair  Pair
	Value float64
}

func (r Rate) String() string {
	return r.Pair.From.String() + "/" + r.Pair.To.String() + ": " + fmt.Sprintf("%f", r.Value) + " on " + r.Date
}
