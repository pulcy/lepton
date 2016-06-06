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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/op/go-logging"

	"github.com/pulcy/lepton/adapters"
)

const (
	logglyAddr          = "https://logs-01.loggly.com"
	logglyEventEndpoint = "/bulk"
	flushTimeout        = 10 * time.Second
)

// Adapter satisfies the router.LogAdapter interface by providing Stream which
// passes all messages to loggly.
type Adapter struct {
	bufferSize int
	log        *logging.Logger
	logglyURL  string
	queue      chan logglyMessage
}

// New returns an Adapter that receives messages from lepton. Additionally,
// it launches a goroutine to buffer and flush messages to loggly.
func New(logglyToken string, tags string, bufferSize int) *Adapter {
	adapter := &Adapter{
		bufferSize: bufferSize,
		log:        logging.MustGetLogger("loggly-adapter"),
		logglyURL:  buildLogglyURL(logglyToken, tags),
		queue:      make(chan logglyMessage),
	}

	go adapter.readQueue()

	return adapter
}

// Stream satisfies the router.LogAdapter interface and passes all messages to
// Loggly
func (l *Adapter) Stream(logstream chan *adapters.Message) {
	for m := range logstream {
		l.queue <- logglyMessage{
			Timestamp:   m.Timestamp,
			Message:     m.Message,
			Unit:        m.Unit,
			Hostname:    m.Hostname,
			MachineID:   m.MachineID,
			J2JobName:   m.J2JobName,
			J2GroupName: m.J2GroupName,
			J2GroupFull: m.J2GroupFull,
			J2TaskName:  m.J2TaskName,
			J2TaskFull:  m.J2TaskFull,
			J2Kind:      m.J2Kind,
			J2Instance:  m.J2Instance,
		}
	}
}

func (l *Adapter) readQueue() {
	buffer := l.newBuffer()

	timeout := time.NewTimer(flushTimeout)

	for {
		select {
		case msg := <-l.queue:
			if len(buffer) == cap(buffer) {
				timeout.Stop()
				l.flushBuffer(buffer)
				buffer = l.newBuffer()
			}

			buffer = append(buffer, msg)

		case <-timeout.C:
			if len(buffer) > 0 {
				l.flushBuffer(buffer)
				buffer = l.newBuffer()
			}
		}

		timeout.Reset(flushTimeout)
	}
}

func (l *Adapter) newBuffer() []logglyMessage {
	return make([]logglyMessage, 0, l.bufferSize)
}

func (l *Adapter) flushBuffer(buffer []logglyMessage) {
	var data bytes.Buffer

	for _, msg := range buffer {
		j, _ := json.Marshal(msg)
		data.Write(j)
		data.WriteString("\n")
	}

	req, _ := http.NewRequest(
		"POST",
		l.logglyURL,
		&data,
	)

	go l.sendRequestToLoggly(req)
}

func (l *Adapter) sendRequestToLoggly(req *http.Request) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		l.log.Errorf("error from client: %#v", err)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		l.log.Errorf("received a %s status code when sending message. response: %s", resp.StatusCode, resp.Body)
	}
}

func buildLogglyURL(token, tags string) string {
	var url string
	url = fmt.Sprintf(
		"%s%s/%s",
		logglyAddr,
		logglyEventEndpoint,
		token,
	)

	if tags != "" {
		url = fmt.Sprintf(
			"%s/tag/%s/",
			url,
			tags,
		)
	}
	return url
}
