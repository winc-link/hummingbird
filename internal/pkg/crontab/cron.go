package crontab

import "github.com/robfig/cron/v3"

var (
	Schedule *cron.Cron
)

func init() {
	Schedule = cron.New()
}

func Start() {
	Schedule.Start()
}

func Stop() {
	Schedule.Stop()
}
