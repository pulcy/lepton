// Copyright (c) 2016 Pulcy.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package loggly

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pulcy/lepton/adapters"
)

func init() {
	adapters.RegisterAdapter("loggly", newLogglyAdapter)
}

func newLogglyAdapter() (adapters.LogAdapter, error) {
	token := os.Getenv("LOGGLY_TOKEN")
	tokenPath := os.Getenv("LOGGLY_TOKEN_FILE")
	tags := os.Getenv("LOGGLY_TAGS")

	if token == "" && tokenPath != "" {
		raw, err := ioutil.ReadFile(tokenPath)
		if err != nil {
			return nil, maskAny(err)
		}
		token = strings.TrimSpace(string(raw))
	}

	if token == "" {
		return nil, maskAny(errors.New("Environment variable LOGGLY_TOKEN not set"))
	}

	return New(token, tags, 100), nil
}
