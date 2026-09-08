package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
)

func main() {
	var config Config
	config.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func(signals chan os.Signal, cancel context.CancelFunc) {
		<-signals
		cancel()
	}(signals, cancel)

	idsChan := make(chan string)

	go ListIds(ctx, idsChan)

	var wg sync.WaitGroup
	for range config.IdWorkers {
		wg.Go(func() {
			fetchForIds(ctx, config, idsChan)
		})
	}

	wg.Wait()
}
