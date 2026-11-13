/*
Copyright 2017 Heptio Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sinks

import (
	"testing"

	"github.com/spf13/viper"
)

// The sink is hardcoded; ManufactureSink must return the stdout sink.
func TestManufactureSinkReturnsStdout(t *testing.T) {
	if _, ok := ManufactureSink().(*StdoutSink); !ok {
		t.Fatal("ManufactureSink must return the stdout sink")
	}
}

// Stale ConfigMaps on upgraded clusters may still carry a "sink" key.
// It must be ignored, not honored or panicked on.
func TestManufactureSinkIgnoresStaleSinkConfig(t *testing.T) {
	viper.Set("sink", "http")
	defer viper.Reset()
	if _, ok := ManufactureSink().(*StdoutSink); !ok {
		t.Fatal("ManufactureSink must ignore a stale sink config value")
	}
}
