package cmd

import (
	"os"
	"runtime/pprof"

	"github.com/rai-project/server"
)

var (
	interrupt = make(chan os.Signal, 2)
)

func handleInterrupt(interrupt chan os.Signal, quitting chan struct{}, svr *server.Server) {
	for _ = range interrupt {
		svr.Lock()
		defer svr.Unlock()

		if svr.Interrupted {
			println("already shutting down")
			continue
		}
		println("shutdown initiated")
		svr.Interrupted = true
		if svr.BeforeShutdown != nil {
			if !svr.BeforeShutdown() {
				svr.Interrupted = false
				continue
			}
		}

		close(quitting)

		if err := svr.Disconnect(); err != nil {
			log.WithError(err).Error("failed to disconnect server")
		}

		if svr.ShutdownInitiated != nil {
			svr.ShutdownInitiated()
		}
		pprof.StopCPUProfile()
	}
}
