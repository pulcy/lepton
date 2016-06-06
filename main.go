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

package main

import (
	"os"

	"github.com/op/go-logging"
	"github.com/spf13/cobra"

	"github.com/pulcy/lepton/service"
)

const (
	defaultAdapterKey  = "loggly"
	defaultJournalPort = 19531
)

var (
	projectName    = "lepton"
	projectVersion = "dev"
	projectBuild   = "dev"
)

var (
	cmdMain = &cobra.Command{
		Use:   projectName,
		Short: "Lepton forwards journal log data to remote destinations",
		Long:  "Lepton forwards journal log data to remote destinations",
		Run:   cmdMainRun,
	}
	log      = logging.MustGetLogger(projectName)
	appFlags struct {
		service.ServiceConfig
	}
)

func init() {
	logging.SetFormatter(logging.MustStringFormatter("[%{level:-5s}] %{message}"))

	cmdMain.Flags().StringVarP(&appFlags.AdapterKey, "adapter", "A", defaultAdapterKey, "Type of adapter to use")
	cmdMain.Flags().IntVarP(&appFlags.JournalPort, "journal-port", "P", defaultJournalPort, "Port of systemd-journal-gatewayd")
}

func main() {
	cmdMain.Execute()
}

func cmdMainRun(cmd *cobra.Command, args []string) {
	assertArgIsSet(appFlags.AdapterKey, "adapter")

	s, err := service.NewService(appFlags.ServiceConfig, service.ServiceDependencies{
		Logger: log,
	})
	if err != nil {
		Exitf("Failed to created service: %#v\n", err)
	}

	if err := s.Run(); err != nil {
		Exitf("Failed to run service: %#v\n", err)
	}
}

func Exitf(format string, args ...interface{}) {
	log.Errorf(format, args...)
	os.Exit(1)
}

func assertArgIsSet(arg, argKey string) {
	if arg == "" {
		Exitf("%s must be set\n", argKey)
	}
}
