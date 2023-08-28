package dtos

import (
	"github.com/winc-link/hummingbird/internal/models"
	"github.com/winc-link/hummingbird/internal/pkg/constants"
	//"github.com/winc-link/hummingbird/proto/resource"
	//"github.com/winc-link/hummingbird/proto/strategy"
)

type AdvanceConfig struct {
	// 日志级别 默认为DEBUG
	LogLevel constants.LogLevel
	// 持久化存储开关 默认关闭
	PersistStorage bool
	// 存储时长 默认为0
	StorageHour int32
}

func AdvanceConfigFromModelToDTO(config models.AdvanceConfig) AdvanceConfig {
	return AdvanceConfig{
		LogLevel:       config.LogLevel,
		PersistStorage: config.PersistStorage,
		StorageHour:    config.StorageHour,
	}
}

func AdvanceConfigFromDTOToModel(config AdvanceConfig) models.AdvanceConfig {
	return models.AdvanceConfig{
		ID:             constants.DefaultAdvanceConfigID,
		LogLevel:       config.LogLevel,
		PersistStorage: config.PersistStorage,
		StorageHour:    config.StorageHour,
	}
}
