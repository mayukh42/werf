package lib

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Terminator(service string) chan bool {
	stop := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	// notify runtime to send the following signals to stop channel
	signal.Notify(stop,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGABRT,
	)

	go func() {
		// block but in another goroutine
		s := <-stop
		log.Fatalf("terminating due to %v", s)
		done <- true
	}()

	return done
}
