package dtos

import "encoding/json"

type ConfigNetWork struct {
	NcId    string `json:"ncId"`
	LocalIp string `json:"localIp,omitempty"`
	GwIp    string `json:"gwIp,omitempty"`
	SmIp    string `json:"smIp,omitempty"`
	Netlink bool   `json:"netlink,omitempty"`
}

type ConfigNetworkUpdateRequest struct {
	NcId    string `json:"ncId" binding:"required"`
	LocalIp string `json:"localIp" binding:"required,ipv4"`
	GwIp    string `json:"gwIp" binding:"required,ipv4"`
	SmIp    string `json:"smIp" binding:"required,ipv4"`
}

type ConfigDnsUpdateRequest struct {
	Dns        []string `json:"dns,omitempty" binding:"required"`
	OpenSwitch bool     `json:"openSwitch,omitempty"`
}

type ConfigNetWorkResponse struct {
	List []ConfigNetWork `json:"list"`
}

func NewConfigNetWorkResponse() ConfigNetWorkResponse {
	return ConfigNetWorkResponse{List: make([]ConfigNetWork, 0)}
}

func (d ConfigNetWorkResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *ConfigNetWorkResponse) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &d)
}

type ConfigDnsResponse struct {
	Dns        []string `json:"dns"`
	OpenSwitch bool     `json:"openSwitch"`
}

func (d ConfigDnsResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}
