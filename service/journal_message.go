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

type JournalEntry struct {
	Pid                int    `json:"_PID,string,omitempty"`
	Uid                int    `json:"_UID,string,omitempty"`
	Gid                int    `json:"_GID,string,omitempty"`
	Exe                string `json:"_EXE,omitempty"`
	Cmdline            string `json:"_CMDLINE,omitempty"`
	SystemdUnit        string `json:"_SYSTEMD_UNIT,omitempty"`
	SystemdSlice       string `json:"_SYSTEMD_SLICE,omitempty"`
	BootId             string `json:"_BOOT_ID,omitempty"`
	MachineId          string `json:"_MACHINE_ID,omitempty"`
	Hostname           string `json:"_HOSTNAME,omitempty"`
	RealtimeTimestamp  int64  `json:"__REALTIME_TIMESTAMP,string,omitempty"`
	MonotonicTimestamp int64  `json:"__MONOTONIC_TIMESTAMP,string,omitempty"`
	Message            string `json:"MESSAGE,omitempty"`
	Priority           int    `json:"PRIORITY,string,omitempty"`
}
