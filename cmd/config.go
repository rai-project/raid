package cmd

import (
	"github.com/k0kubun/pp"
	"github.com/rai-project/config"
	"github.com/rai-project/vipertags"
)

type raidConfig struct {
	QueueName string `json:"queue_name" config:"raid.queue_name" default:"rai"`
	done      chan struct{}
}

// Config ...
var (
	Config = &raidConfig{
		done: make(chan struct{}),
	}
)

// ConfigName ...
func (raidConfig) ConfigName() string {
	return "raid"
}

// SetDefaults ...
func (a *raidConfig) SetDefaults() {
	vipertags.SetDefaults(a)
}

// Read ...
func (a *raidConfig) Read() {
	defer close(a.done)
	vipertags.Fill(a)
}

// Wait ...
func (c raidConfig) Wait() {
	<-c.done
}

// String ...
func (c raidConfig) String() string {
	return pp.Sprintln(c)
}

// Debug ...
func (c raidConfig) Debug() {
	log.Debug("raid Config = ", c)
}

func init() {
	config.Register(Config)
}
