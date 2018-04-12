package cmd

import (
	"github.com/sirupsen/logrus"

	"github.com/rai-project/config"
	logger "github.com/rai-project/logger"
	_ "github.com/rai-project/logger/hooks"
)


var (
	log *logrus.Entry
)


func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "raid")
	})
}
