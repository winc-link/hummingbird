/*******************************************************************************
 * Copyright 2017 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/winc-link/hummingbird/cmd/mqtt-broker/initcmd"

	_ "github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/persistence"
	_ "github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/topicalias/fifo"

	_ "github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/plugin/admin"
	_ "github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/plugin/aplugin"
	_ "github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/plugin/auth"
	_ "github.com/winc-link/hummingbird/internal/hummingbird/mqttbroker/plugin/federation"

	"net/http"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:     "mqttd",
		Long:    "This is a MQTT broker that fully implements MQTT V5.0 and V3.1.1 protocol",
		Version: "",
	}
	enablePprof bool
	pprofAddr   = "127.0.0.1:60600"
)

func main() {
	if enablePprof {
		go func() {
			http.ListenAndServe(pprofAddr, nil)
		}()
	}
	initcmd.Init(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
