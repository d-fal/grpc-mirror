package control

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// SignalListener signal listener
func SignalListener(systemWideContext context.Context, systemWideCancel context.CancelFunc) {

	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGSTOP)

	go func() {
		select {
		case <-c:
			Shutdown(systemWideCancel)
		}
	}()

}
