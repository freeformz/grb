package graceful

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/freeformz/grb/Godeps/_workspace/src/github.com/go-martini/martini"
)

type GracefulShutdown struct {
	Timeout time.Duration
	wg      sync.WaitGroup
	Log     *log.Logger
}

func NewGracefulShutdown(t time.Duration) *GracefulShutdown {
	return &GracefulShutdown{Timeout: t}
}

func (g *GracefulShutdown) Handler(c martini.Context) {
	g.wg.Add(1)
	c.Next()
	g.wg.Done()
}

func (g *GracefulShutdown) WaitForSignal(signals ...os.Signal) error {
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, signals...)
	<-sigchan

	if g.Log != nil {
		g.Log.Println("Waiting for all requests to finish")
	} else {
		log.Println("Waiting for all requests to finish")
	}

	waitChan := make(chan struct{})
	go func() {
		g.wg.Wait()
		waitChan <- struct{}{}
	}()

	select {
	case <-time.After(g.Timeout):
		return fmt.Errorf("timed out waiting %v for shutdown", g.Timeout)
	case <-waitChan:
		return nil
	}
}
