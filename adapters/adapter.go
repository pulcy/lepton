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

package adapters

type AdapterFactory func() (LogAdapter, error)

type LogAdapter interface {
	Stream(logstream chan *Message)
}

var (
	factories = make(map[string]AdapterFactory)
)

func RegisterAdapter(key string, factory AdapterFactory) {
	factories[key] = factory
}

func CreateAdapter(key string) (LogAdapter, error) {
	factory, ok := factories[key]
	if !ok {
		return nil, maskAny(NoSuchAdapterError)
	}
	result, err := factory()
	if err != nil {
		return nil, maskAny(err)
	}
	return result, nil
}
