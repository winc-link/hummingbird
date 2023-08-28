package leveldb

import (
	"strings"
)

const (
	DeviceIdDpId = "DeviceIdDpId"
	Tab          = ":"
)

func genDeviceIdDpId(deviceId, dpId string) string {
	return getJointIds(DeviceIdDpId, deviceId, dpId)
}

func getJointIds(ids ...string) (ret string) {
	if len(ids) == 0 {
		return ""
	}
	ret = strings.Join(ids, Tab)
	return
}
