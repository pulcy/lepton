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
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/op/go-logging"

	"github.com/pulcy/lepton/adapters"
)

type ServiceConfig struct {
	AdapterKey  string
	JournalPort int
}

type ServiceDependencies struct {
	Logger *logging.Logger
}

type service struct {
	ServiceConfig
	ServiceDependencies

	adapter adapters.LogAdapter
}

func NewService(config ServiceConfig, deps ServiceDependencies) (*service, error) {
	adapter, err := adapters.CreateAdapter(config.AdapterKey)
	if err != nil {
		return nil, maskAny(err)
	}
	return &service{
		ServiceConfig:       config,
		ServiceDependencies: deps,
		adapter:             adapter,
	}, nil
}

func (s *service) Run() error {
	entries := make(chan *JournalEntry)
	messages := make(chan *adapters.Message)

	go s.processMessages(entries, messages)
	go s.adapter.Stream(messages)

	s.Logger.Debugf("start processing messages to %s", s.AdapterKey)
	if err := s.collectEntries(entries); err != nil {
		return maskAny(err)
	}
	return nil
}

func (s *service) processMessages(entries chan *JournalEntry, messages chan *adapters.Message) {
	for entry := range entries {
		message := &adapters.Message{
			// Journald timestamp is in microseconds, so they are being converted to nano.
			Timestamp: time.Unix(0, entry.RealtimeTimestamp*1000),
			Message:   entry.Message,
			Unit:      entry.SystemdUnit,
			Hostname:  entry.Hostname,
			MachineID: entry.MachineId,
		}
		j2names, ok := ParseUnitName(entry.SystemdUnit)
		if ok {
			message.J2JobName = j2names.JobName
			message.J2GroupName = j2names.GroupName
			message.J2GroupFull = fmt.Sprintf("%s.%s", j2names.JobName, j2names.GroupName)
			message.J2TaskName = j2names.TaskName
			message.J2TaskFull = fmt.Sprintf("%s.%s.%s", j2names.JobName, j2names.GroupName, j2names.TaskName)
			message.J2Kind = j2names.Kind
			message.J2Instance = j2names.Instance
		}
		messages <- message
	}
}

func (s *service) collectEntries(entries chan *JournalEntry) error {
	url := fmt.Sprintf("http://localhost:%d/entries?follow", s.JournalPort)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return maskAny(err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return maskAny(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return maskAny(fmt.Errorf("Invalid status code %d", resp.StatusCode))
	}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		msg := scanner.Text()
		var entry JournalEntry
		err := json.Unmarshal([]byte(strings.Replace(msg, "data:", "", 1)), &entry)
		if err != nil {
			s.Logger.Debugf("failed to parse journal message: %#v", err)
			// Ignore blank lines
		} else {
			entries <- &entry
		}
	}

	return nil
}
