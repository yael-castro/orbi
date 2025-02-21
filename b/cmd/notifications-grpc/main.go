package main

import (
	"context"
	"github.com/yael-castro/orbi/b/internal/container"
	"github.com/yael-castro/orbi/b/internal/runtime"
	"google.golang.org/grpc"
	"log"
	"net"
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

	// Listening through TCP
	listener, err := net.Listen("tcp", ":"+runtime.Port)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		_ = listener.Close()
		log.Println("Listener stopped")
	}()

	c := container.New()

	// DI in action!
	var server *grpc.Server

	if err = c.Inject(ctx, &server); err != nil {
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

		server.GracefulStop()
		log.Println("Server stopped")
	}()

	// Listening for server error
	errCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Printf("Starting server at port %s\n", runtime.Port)
		errCh <- server.Serve(listener)
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
