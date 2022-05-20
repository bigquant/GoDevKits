package app_context

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// AppContext interface
type AppContext interface {
	Context() context.Context
	EndContext()
	End()
	WaitGroup() *sync.WaitGroup
}

type appContext struct {
	ctx    context.Context
	endCtx func()
	wg     *sync.WaitGroup
}

func (a *appContext) Context() context.Context {
	return a.ctx
}

func (a *appContext) EndContext() {
	a.endCtx()
}

func (a *appContext) End() {
	defer a.wg.Wait()
	a.EndContext()
}

func (a *appContext) WaitGroup() *sync.WaitGroup {
	return a.wg
}

// Start func
func Start() AppContext {

	logger := log.New(os.Stderr, "[app_context]", log.Flags())
	interruptChan := make(chan os.Signal, 2)
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGTERM)
	ctx, endCtx := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case sig := <-interruptChan:
			logger.Printf("signal %v", sig)
			endCtx()
		case <-ctx.Done():
		}
		signal.Stop(interruptChan)
	}()

	return &appContext{ctx, endCtx, wg}
}
