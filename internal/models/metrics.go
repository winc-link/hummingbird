package models

type Metrics struct {
	Key            string `gorm:"column:key;pk"`
	Timestamp      int64  `gorm:"column:timestamp"` // 时间戳
	CpuUsedPercent float64
	MemoryUsed     int64
}

func (table *Metrics) TableName() string {
	return "metrics"
}

func (table *Metrics) Get() interface{} {
	return *table
}

type SystemMetrics struct {
	ID        int64  `gorm:"column:id;pk;autoIncrement"`
	Data      string `gorm:"column:data"`
	Timestamp int64  `gorm:"column:timestamp"` // 时间戳
}

func (table *SystemMetrics) TableName() string {
	return "system_metrics"
}

func (table *SystemMetrics) Get() interface{} {
	return *table
}
