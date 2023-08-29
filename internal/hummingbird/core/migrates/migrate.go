package migrates

import (
	"github.com/go-gormigrate/gormigrate/v2"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/winc-link/hummingbird/internal/pkg/container"
	"github.com/winc-link/hummingbird/internal/pkg/di"

	bootstrapContainer "github.com/winc-link/hummingbird/internal/hummingbird/core/container"
)

func Migrate(dic *di.Container) {
	dbClient := bootstrapContainer.DBClientFrom(dic.Get)
	lc := container.LoggingClientFrom(dic.Get)
	m := gormigrate.New(dbClient.GetDBInstance(), gormigrate.DefaultOptions, migrations())

	if err := m.Migrate(); err != nil {
		lc.Errorf("Migration run err: %v", err)
	} else {
		lc.Info("Migration run successfully")
	}
}

func migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		//m_1627356992_funcpoint_properties(),
		//m_1630660539_device_expand_data(),
		//m_1630660539_update_screen_device(),
		//m_1637287851_update_library_internal(),
		//m_1641866282_upgrade_cloud_market(),
		//m_1648619146_upgrade_driver_app(),
	}
}
