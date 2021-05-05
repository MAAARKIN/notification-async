package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/MAAARKIN/notification-async/basic"
	"github.com/MAAARKIN/notification-async/concurrent"
	"github.com/MAAARKIN/notification-async/model"
)

func main() {
	// create a context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// that cancels at ctrl+C
	go onSignal(os.Interrupt, cancel)

	// parse command line arguments
	op := new(model.Options)

	flag.StringVar(&op.Filename, "filename", "", "src file")
	flag.StringVar(&op.Event, "event", "", "event name")
	numberOfWorkers := flag.Int("workers", 2, "concurrent workers")
	flag.BoolVar(&op.Async, "async", false, "if the process will be async")
	flag.Parse()

	// check arguments
	if op.Filename == "" {
		log.Fatal("filename required")
	}

	if op.Event == "" {
		log.Fatal("event required")
	}

	if op.Async {
		// check arguments
		if numberOfWorkers == nil {
			log.Fatal("workers required")
		}
		concurrent.Start(ctx, *op, *numberOfWorkers)
	} else {
		basic.Start(ctx, *op)
	}

}

func onSignal(s os.Signal, cancel func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, s)
	<-c
	cancel()
}
