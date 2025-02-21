package main

import (
	"context"
	"github.com/yael-castro/orbi/b/internal/container"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Setting logger flags
	log.SetFlags(log.Flags() | log.Lshortfile)

	var consumer func(context.Context) error

	// DI in action!
	c := container.New()

	if err := c.Inject(ctx, &consumer); err != nil {
		log.Println(err)
		return
	}

	wg := sync.WaitGroup{}

	// Listening for graceful shutdown
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		const gracePeriod = 3 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()

		if err := c.Close(ctx); err != nil {
			log.Println(err)
			return
		}

		log.Println("Container is closed")
	}()

	// Listening events
	errCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Println("Listening events...")
		errCh <- consumer(ctx)
	}()

	// Waiting for context cancellation OR some error
	select {
	case <-ctx.Done():
	case err := <-errCh:
		log.Println(err)
	}

	stop()
	wg.Wait()
}
