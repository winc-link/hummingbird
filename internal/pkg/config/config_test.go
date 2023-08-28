//
// Copyright (c) 2019 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUrl(t *testing.T) {
	expected := "https://localhost:8080"
	target := ServiceConfig{
		Protocol: "https",
		Host:     "localhost",
		Port:     8080,
		Type:     "consul",
	}

	actual := target.GetUrl()
	assert.Equal(t, expected, actual)
}

func TestGetProtocol(t *testing.T) {
	testCases := []struct {
		Name     string
		Protocol string
		Expected string
	}{
		{
			Name:     "Protocol Specified",
			Protocol: "https",
			Expected: "https",
		},
		{
			Name:     "Protocol Not Specified",
			Protocol: "",
			Expected: "http",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			target := ServiceConfig{
				Protocol: test.Protocol,
			}

			actual := target.GetProtocol()
			assert.Equal(t, test.Expected, actual)
		})
	}
}

func TestPopulateFromUrl(t *testing.T) {
	testCases := []struct {
		Name             string
		Url              string
		ExpectedType     string
		ExpectedProtocol string
		ExpectedHost     string
		ExpectedPort     int
		ExpectedError    string
	}{
		{
			Name:             "Success, protocol specified",
			Url:              "consul.https://localhost:8080",
			ExpectedType:     "consul",
			ExpectedProtocol: "https",
			ExpectedHost:     "localhost",
			ExpectedPort:     8080,
		},
		{
			Name:             "Success, protocol not specified",
			Url:              "consul://localhost:8080",
			ExpectedType:     "consul",
			ExpectedProtocol: "http",
			ExpectedHost:     "localhost",
			ExpectedPort:     8080,
		},
		{
			Name:          "Bad URL format",
			Url:           "not a url\r\n",
			ExpectedError: "the format of Provider URL is incorrect",
		},
		{
			Name:          "Bad Port",
			Url:           "consul.https:\\localhost:eight",
			ExpectedError: "the port from Provider URL is incorrect",
		},
		{
			Name:          "Missing Type and Protocol spec",
			Url:           "://localhost:800",
			ExpectedError: "missing protocol scheme",
		},
		{
			Name:          "Bad Type and Protocol spec",
			Url:           "xyz.consul.http://localhost:800",
			ExpectedError: "the Type and Protocol spec from Provider URL is incorrect",
		},
	}

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			target := ServiceConfig{}

			err := target.PopulateFromUrl(test.Url)
			if test.ExpectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.ExpectedError)
				return // test is complete
			}

			assert.Equal(t, test.ExpectedType, target.Type)
			assert.Equal(t, test.ExpectedProtocol, target.Protocol)
			assert.Equal(t, test.ExpectedHost, target.Host)
			assert.Equal(t, test.ExpectedPort, target.Port)
		})
	}
}
