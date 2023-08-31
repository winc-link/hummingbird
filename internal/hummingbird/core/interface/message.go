package interfaces

import (
	"context"
	"github.com/winc-link/edge-driver-proto/drivercommon"
	"github.com/winc-link/hummingbird/internal/dtos"
	"github.com/winc-link/hummingbird/internal/hummingbird/core/application/messagestore"
	//msgTypes "gitlab.com/tedge/edgex/internal/pkg/messaging/types"
	//
	//"gitlab.com/tedge/edgex/proto/thingmodel"
	//
	//pkgMQTT "gitlab.com/tedge/edgex/internal/tools/mqttclient"
	//
	//mqtt "github.com/eclipse/paho.mqtt.golang"
	//
	//"gitlab.com/tedge/edgex/internal/dtos"
)

type MessageStores interface {
	StoreRange()
	StoreMsgId(id string, ch string)
	LoadMsgChan(id string) (interface{}, bool)
	DeleteMsgId(id string)
	GenAckChan(id string) *messagestore.MsgAckChan
}

type MessageItf interface {
	TyCloudMqttItf
}

type PublishCallback func(ctx context.Context, params ...interface{}) (bool, interface{})

type TyCloudMqttItf interface {
	// ThingModelMsgReport 物模型消息上报到云端
	ThingModelMsgReport(ctx context.Context, msg dtos.ThingModelMessage) (*drivercommon.CommonResponse, error)
	DeviceStatusToMessageBus(ctx context.Context, deviceId, deviceStatus string)
}
