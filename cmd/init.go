package cmd

import (
	"github.com/Sirupsen/logrus"

	"github.com/rai-project/config"
	logger "github.com/rai-project/logger"
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
	config.OnInit(func() {
		log = dlog{logger.New().WithField("pkg", "raid")}
	})
}
