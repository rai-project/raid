package cmd

import (
	"github.com/sirupsen/logrus"

	"github.com/rai-project/config"
	logger "github.com/rai-project/logger"
	_ "github.com/rai-project/logger/hooks"
)

type dlog struct {
	*logrus.Entry
}

var (
	log dlog
)

func (d dlog) Error(v ...interface{}) error {
	d.Error(v...)
	return nil
}
func (d dlog) Warn(v ...interface{}) error {
	d.Warn(v...)
	return nil
}

func init() {
	config.AfterInit(func() {
		log = dlog{logger.New().WithField("pkg", "raid")}
	})
}
