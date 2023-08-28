package config

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type ConfigurationStruct struct {
	Writable     WritableInfo
	Clients      map[string]ClientInfo
	Proxy        map[string]ClientInfo
	Databases    map[string]Database
	Service      ServiceInfo
	RpcServer    RPCServiceInfo
	MessageQueue MessageQueueInfo
	Topics       struct {
		CommandTopic  TopicInfo
		CloudnetTopic TopicInfo
	}
	ApplicationSettings struct {
		CloseAuthToken bool
	}
}

type TopicInfo struct {
	Topic string
	Group string
}

type WritableInfo struct {
	LogLevel                        string
	LogPath                         string
	EnableValueDescriptorManagement bool
	DebugProfile                    bool
	IsNewModel                      bool
}

// MessageQueueInfo provides parameters related to connecting to a message bus
type MessageQueueInfo struct {
	// Host is the hostname or IP address of the broker, if applicable.
	Host string
	// Port defines the port on which to access the message queue.
	Port int
	// Protocol indicates the protocol to use when accessing the message queue.
	Protocol string
	// Indicates the message queue platform being used.
	Type string
	// Indicates the topic the data is published/subscribed
	// TODO this configuration shall be removed once v1 API is deprecated.
	EventTopic string
	// Indicates the topic the data is published/subscribed
	SubscribeTopics []string
	// sync device status to cloud by message queue
	DeviceStatusTopic string
	// Indicates the topic prefix the data is published to. Note that /<device-profile-name>/<device-name> will be
	// added to this Publish Topic prefix as the complete publish topic
	PublishTopicPrefix string
	// Provides additional configuration properties which do not fit within the existing field.
	// Typically the key is the name of the configuration property and the value is a string representation of the
	// desired value for the configuration property.
	Optional map[string]string
}

func (m MessageQueueInfo) Enable() bool {
	return !(m.Host == "" || m.Port <= 0)
}

var (
	config ConfigurationStruct
)

func GetConfig() ConfigurationStruct {
	return config
}

func ReadConfig() (err error) {
	var data []byte
	if data, err = ioutil.ReadFile(ConfDir + "/configuration.toml"); err != nil {
		err = errors.New(fmt.Sprintf("failed read configuration file: %s \n", err))
		return
	}
	if err = toml.Unmarshal(data, &config); err != nil {
		err = errors.New(fmt.Sprintf("failed unmarshal configuration file: %s \n", err))
	}
	return
}

var (
	ConfDir string
)
