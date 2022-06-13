package crons

import (
	"whale-tracker/src/crons/cronjob"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

func Init() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	log.Info("Create new cron")
	c := cron.New()
	c.AddFunc("* 1 * * *", func() { log.Info("[Job 1]Every minute job\n") })

	c.AddFunc("@every 1m", cronjob.WhaleTrackerCron)

	// Start cron with one scheduled job
	log.Info("Start cron")
	c.Start()
}
