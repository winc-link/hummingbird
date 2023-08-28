package dtos

import (
	//"gitlab.com/tedge/edgex/internal/models"
	//"gitlab.com/tedge/edgex/proto/devicelibrary"
	"github.com/winc-link/hummingbird/internal/models"
)

type DeviceLibrarySupportVersion struct {
	Version            string `json:"version"`
	IsDefault          bool   `json:"is_default"`
	DockerParamsSwitch bool   `json:"docker_params_switch"`
	DockerParams       string `json:"docker_params"`
	ExpertMode         bool   `json:"expert_mode"`
	ExpertModeContent  string `json:"expert_mode_content"`
	ConfigFile         string `json:"config_file"`
	ConfigJson         string `json:"config_json"`
}

func SupperVersionsFromModel(versions []models.SupportVersion) []DeviceLibrarySupportVersion {
	ret := make([]DeviceLibrarySupportVersion, 0)
	for _, v := range versions {
		ret = append(ret, DeviceLibrarySupportVersion{
			Version:            v.Version,
			IsDefault:          v.IsDefault,
			DockerParamsSwitch: v.DockerParamsSwitch,
			DockerParams:       v.DockerParams,
			ExpertMode:         v.ExpertMode,
			ExpertModeContent:  v.ExpertModeContent,
			ConfigJson:         v.ConfigJson,
			ConfigFile:         v.ConfigFile,
		})
	}
	return ret
}

//func FromDeviceLibrarySupperVersionsToRpc(versions []models.SupportVersion) []*devicelibrary.SupportVersion {
//	ret := make([]*devicelibrary.SupportVersion, 0)
//	for _, v := range versions {
//		ret = append(ret, &devicelibrary.SupportVersion{
//			Version:            v.Version,
//			IsDefault:          v.IsDefault,
//			DockerParamsSwitch: v.DockerParamsSwitch,
//			DockerParams:       v.DockerParams,
//			ExpertMode:         v.ExpertMode,
//			ExpertModeContent:  v.ExpertModeContent,
//			ConfigJson:         v.ConfigJson,
//			ConfigFile:         v.ConfigFile,
//		})
//	}
//	return ret
//}

type DeviceLibrarySupportVersionSimple struct {
	Version    string `json:"version"`
	IsDefault  bool   `json:"is_default"`
	ConfigFile string `json:"config_file"`
}

func DeviceLibrarySupportVersionSimpleFromModel(versions models.SupportVersions) []DeviceLibrarySupportVersionSimple {
	ret := make([]DeviceLibrarySupportVersionSimple, len(versions))
	for i, v := range versions {
		ret[i] = DeviceLibrarySupportVersionSimple{
			Version:   v.Version,
			IsDefault: v.IsDefault,
			//ConfigFile: v.ConfigFile,
		}
	}
	return ret
}

//func FromSupportVersionSimpleRpcToDto(resp *devicelibrary.DeviceLibrary) []DeviceLibrarySupportVersionSimple {
//	ret := make([]DeviceLibrarySupportVersionSimple, 0)
//	for _, v := range resp.SupportVersions {
//		ret = append(ret, DeviceLibrarySupportVersionSimple{
//			Version:    v.Version,
//			IsDefault:  v.IsDefault,
//			ConfigFile: v.ConfigFile,
//		})
//	}
//	return ret
//}
//
//func ModelSupportVersionFromRPC(s *devicelibrary.SupportVersion) models.SupportVersion {
//	return models.SupportVersion{
//		Version:            s.Version,
//		IsDefault:          s.IsDefault,
//		ConfigJson:         s.ConfigJson,
//		ConfigFile:         s.ConfigFile,
//		DockerParamsSwitch: s.DockerParamsSwitch,
//		DockerParams:       s.DockerParams,
//		ExpertMode:         s.ExpertMode,
//		ExpertModeContent:  s.ExpertModeContent,
//	}
//}
