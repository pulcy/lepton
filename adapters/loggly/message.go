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
	"time"
)

type logglyMessage struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Unit      string    `json:"unit,omitempty"`
	Hostname  string    `json:"hostname,omitempty"`
	MachineID string    `json:"machine_id,omitempty"`

	J2JobName   string `json:"j2_job,omitempty"`
	J2GroupName string `json:"j2_group,omitempty"`
	J2GroupFull string `json:"j2_group_full,omitempty"`
	J2TaskName  string `json:"j2_task,omitempty"`
	J2TaskFull  string `json:"j2_task_full,omitempty"`
	J2Kind      string `json:"j2_kind,omitempty"`
	J2Instance  string `json:"j2_instance,omitempty"`
}
