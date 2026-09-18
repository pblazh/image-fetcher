package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {
	var config Config
	err := config.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}

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

	var wg errgroup.Group
	for range config.IdWorkers {
		wg.Go(func() error {
			return fetchForIds(ctx, config, idsChan)
		})
	}

	err = wg.Wait()
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	os.Exit(0)
}
