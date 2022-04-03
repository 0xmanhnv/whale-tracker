package crons

import (
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

func Init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	log.Info("Create new cron")
	c := cron.New()
	c.AddFunc("*/1 * * * *", func() { log.Info("[Job 1]Every minute job\n") })

	// Start cron with one scheduled job
	log.Info("Start cron")
	c.Start()
}
