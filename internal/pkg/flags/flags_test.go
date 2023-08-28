/*******************************************************************************
 * Copyright 2020 Intel Corp.
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

package flags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// newSUT creates and returns a new "system under test" instance.
func newSUT(args []string) *Default {
	actual := New()
	actual.Parse(args)
	return actual
}

func TestNewAllFlags(t *testing.T) {
	expectedProfile := "docker"
	expectedConfigDirectory := "/res"
	expectedFileName := "config.toml"

	actual := newSUT(
		[]string{
			"-o",
			"-r",
			"-p=" + expectedProfile,
			"-c=" + expectedConfigDirectory,
			"-f=" + expectedFileName,
		},
	)

	assert.Equal(t, true, actual.OverwriteConfig())
	assert.True(t, actual.UseRegistry())
	assert.Equal(t, expectedProfile, actual.Profile())
	assert.Equal(t, expectedConfigDirectory, actual.ConfigDirectory())
	assert.Equal(t, expectedFileName, actual.ConfigFileName())
}

func TestNewDefaultsNoFlags(t *testing.T) {
	actual := newSUT([]string{})

	assert.Equal(t, false, actual.OverwriteConfig())
	assert.False(t, actual.UseRegistry())
	assert.Equal(t, "", actual.ConfigProviderUrl())
	assert.Equal(t, "", actual.Profile())
	assert.Equal(t, "", actual.ConfigDirectory())
	assert.Equal(t, DefaultConfigFile, actual.ConfigFileName())
}

func TestNewDefaultForCP(t *testing.T) {
	actual := newSUT([]string{"-cp"})

	assert.Equal(t, DefaultConfigProvider, actual.ConfigProviderUrl())
}

func TestNewOverrideForCP(t *testing.T) {
	expectedConfigProviderUrl := "consul.http://docker-core-consul:8500"

	actual := newSUT([]string{"-cp=" + expectedConfigProviderUrl})

	assert.Equal(t, expectedConfigProviderUrl, actual.ConfigProviderUrl())
}

func TestNewDefaultForConfigProvider(t *testing.T) {
	actual := newSUT([]string{"-configProvider"})

	assert.Equal(t, DefaultConfigProvider, actual.ConfigProviderUrl())
}

func TestNewOverrideConfigProvider(t *testing.T) {
	expectedConfigProviderUrl := "consul.http://docker-core-consul:8500"

	actual := newSUT([]string{"-configProvider=" + expectedConfigProviderUrl})

	assert.Equal(t, expectedConfigProviderUrl, actual.ConfigProviderUrl())
}

func TestDashR(t *testing.T) {
	expectedConfigDirectory := "/foo/ba-r/"
	actual := newSUT([]string{"-confdir", "/foo/ba-r/"})

	assert.Equal(t, expectedConfigDirectory, actual.ConfigDirectory())
}

func TestConfDirEquals(t *testing.T) {
	expectedConfigDirectory := "/foo/ba-r/"
	actual := newSUT([]string{"-confdir=/foo/ba-r/"})

	assert.Equal(t, expectedConfigDirectory, actual.ConfigDirectory())
}

func TestConfCommonScenario(t *testing.T) {
	expectedConfigProviderUrl := "consul.http://edgex-core-consul:8500"
	expectedConfigDirectory := "/res"

	actual := newSUT([]string{"-cp=consul.http://edgex-core-consul:8500", "--registry", "--confdir=/res"})

	assert.Equal(t, expectedConfigProviderUrl, actual.ConfigProviderUrl())
	assert.True(t, actual.UseRegistry())
	assert.Equal(t, expectedConfigDirectory, actual.ConfigDirectory())
}
