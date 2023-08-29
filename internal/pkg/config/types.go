/*******************************************************************************
 * Copyright 2023 Winc link Inc.
 * Copyright 2020 Intel Inc.
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
package config

import (
	"fmt"
	"strings"
	"time"
)

// ServiceInfo contains configuration settings necessary for the basic operation of any EdgeX service.
type ServiceInfo struct {
	// BootTimeout indicates, in milliseconds, how long the service will retry connecting to upstream dependencies
	// before giving up. Default is 30,000.
	BootTimeout int
	// Health check interval
	CheckInterval string
	// Host is the hostname or IP address of the service.
	Host string
	// Port is the HTTP port of the service.
	Port int
	// ServerBindAddr specifies an IP address or hostname
	// for ListenAndServe to bind to, such as 0.0.0.0
	ServerBindAddr string
	// The protocol that should be used to call this service
	Protocol string
	// StartupMsg specifies a string to log once service
	// initialization and startup is completed.
	StartupMsg string
	// MaxResultCount specifies the maximum size list supported
	// in response to REST calls to other services.
	MaxResultCount int
	// Timeout specifies a timeout (in milliseconds) for
	// processing REST calls from other services.
	Timeout int
}

// Url provides a way to obtain the full url of the host service for use in initialization or, in some cases,
// responses to a caller.
func (s ServiceInfo) Url() string {
	url := fmt.Sprintf("%s://%s:%v", s.Protocol, s.Host, s.Port)
	return url
}

// ConfigProviderInfo defines the type and location (via host/port) of the desired configuration provider (e.g. Consul, Eureka)
type ConfigProviderInfo struct {
	Host string
	Port int
	Type string
}

// RegistryInfo defines the type and location (via host/port) of the desired service registry (e.g. Consul, Eureka)
type RegistryInfo struct {
	Host string
	Port int
	Type string
}

// ClientInfo provides the host and port of another service in the eco-system.
type ClientInfo struct {
	// Host is the hostname or IP address of a service.
	Host string
	// Port defines the port on which to access a given service
	Port int
	// Protocol indicates the protocol to use when accessing a given service
	Protocol string
	// 是否启用tls
	UseTLS bool
	// ca cert
	CertFilePath string
}

func (c ClientInfo) Address() string {
	lower := strings.ToLower(c.Protocol)
	if lower == "http" {
		return c.Url()
	} else if lower == "grpc" {
		return c.RpcUrl()
	}
	return ""
}

func (c ClientInfo) Url() string {
	return fmt.Sprintf("%s://%s:%v", c.Protocol, c.Host, c.Port)
}

func (c ClientInfo) RpcUrl() string {
	return fmt.Sprintf("%s:%v", c.Host, c.Port)
}

// SecretStoreInfo encapsulates configuration properties used to create a SecretClient.
type SecretStoreInfo struct {
	Host           string
	Port           int
	Path           string
	Protocol       string
	Namespace      string
	RootCaCertPath string
	ServerName     string
	//Authentication          types.AuthenticationInfo
	Authentication          AuthenticationInfo
	AdditionalRetryAttempts int
	RetryWaitPeriod         string
	retryWaitPeriodTime     time.Duration
	// TokenFile provides a location to a token file.
	TokenFile string
}

type Database struct {
	Type       string
	Timeout    int
	Host       string
	Port       int
	Name       string
	DataSource string
	Dsn        string
	Cluster    []string
	Username   string
	Password   string
}

// Credentials encapsulates username-password attributes.
type Credentials struct {
	Username string
	Password string
}

//CertKeyPair encapsulates public certificate/private key pair for an SSL certificate
type CertKeyPair struct {
	Cert string
	Key  string
}

// InsecureSecrets is used to hold the secrets stored in the configuration
type InsecureSecrets map[string]InsecureSecretsInfo

// InsecureSecretsInfo encapsulates info used to retrieve insecure secrets
type InsecureSecretsInfo struct {
	Path    string
	Secrets map[string]string
}

// rpc server info
type RPCServiceInfo struct {
	Address  string
	UseTLS   bool
	CertFile string
	KeyFile  string
	// 服务器响应超时时间，单位为毫秒
	Timeout      int
	LimitSetting LimitSetting
}

type LimitSetting struct {
	LimitTime int64 // 发放的间隔时间
	Burst     int   // 发放的个数
}

// BootstrapConfiguration defines the configuration elements required by the bootstrap.
type BootstrapConfiguration struct {
	Clients     map[string]ClientInfo
	Service     ServiceInfo
	RpcServer   RPCServiceInfo
	Config      ConfigProviderInfo
	Registry    RegistryInfo
	SecretStore SecretStoreInfo
}

//test!!!
// AuthenticationInfo contains authentication information to be used when communicating with an HTTP based provider
type AuthenticationInfo struct {
	AuthType  string
	AuthToken string
}
