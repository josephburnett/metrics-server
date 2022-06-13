// Copyright 2018 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"os"
	"runtime"

	genericapiserver "k8s.io/apiserver/pkg/server"
	restclient "k8s.io/client-go/rest"
	"k8s.io/component-base/logs"

	"sigs.k8s.io/metrics-server/cmd/metrics-server/app"
	"sigs.k8s.io/metrics-server/pkg/podautoscaler"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	go func() {
		config, err := restclient.InClusterConfig()
		if err != nil {
			panic(err)
		}
		controllerFactory := podautoscaler.ControllerFactory{
			StopCh:     make(<-chan struct{}),
			KubeConfig: config,
		}
		horizontalController, err := controllerFactory.Make()
		if err != nil {
			panic(err)
		}
		horizontalController.Run(make(<-chan struct{}))
	}()

	cmd := app.NewMetricsServerCommand(genericapiserver.SetupSignalHandler())
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
