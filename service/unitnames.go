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

package service

import (
	"strings"
)

var (
	unitSuffixes = []string{".service", ".timer"}
)

type J2Names struct {
	JobName   string
	GroupName string
	TaskName  string
	Kind      string
	Instance  string
}

// ParseUnitName tries to split the given unit name into J2 generate job/group/task names.
func ParseUnitName(unit string) (J2Names, bool) {
	found, suffix := hasUnitSuffix(unit)
	if !found {
		return J2Names{}, false
	}

	unit = unit[:len(unit)-len(suffix)]
	parts := strings.Split(unit, "-")
	if len(parts) != 4 {
		return J2Names{}, false
	}

	result := J2Names{
		JobName:   parts[0],
		GroupName: parts[1],
		TaskName:  parts[2],
	}

	parts = strings.Split(parts[3], "@")
	result.Kind = parts[0]
	if len(parts) >= 2 {
		result.Instance = parts[1]
	}
	return result, true
}

func hasUnitSuffix(unitName string) (bool, string) {
	for _, suffix := range unitSuffixes {
		if strings.HasSuffix(unitName, suffix) {
			return true, suffix
		}
	}
	return false, ""
}
