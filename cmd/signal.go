//+build !appengine

package cmd

import (
	"os"
	"os/signal"
	"syscall"
)

func signalNotify(interrupt chan<- os.Signal) {
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
}
