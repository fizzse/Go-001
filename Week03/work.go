package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type Handel struct {
}

func (c *Handel) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ping"))
}

func (c Handel) Ping(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ping"))
}

func serveHttp(exitChan <-chan struct{}, cancelFunc context.CancelFunc, addr string) error {
	defer func() {
		log.Println("http server prepare exit")
		cancelFunc()
	}()

	http.HandleFunc("/ping", (&Handel{}).Ping)
	srv := &http.Server{
		Addr:    addr,
		Handler: &Handel{},
	}

	go func() {
		select {
		case <-exitChan:
			ctx1, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			srv.Shutdown(ctx1)
		}
	}()

	err := srv.ListenAndServe()

	return fmt.Errorf("serve http failed: %w", err)
}

func handleSignal(exitChan <-chan struct{}, cancelFunc context.CancelFunc) error {
	defer func() {
		log.Println("handle signal prepare exit")
		cancelFunc()
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sign := <-quit:
		err := fmt.Errorf("handle signal: %v", sign)
		return err

	case <-exitChan:
		return nil
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	addr := "0.0.0.0:8080"
	group, _ := errgroup.WithContext(ctx)
	group.Go(func() error {
		return serveHttp(ctx.Done(), cancel, addr)
	})

	group.Go(func() error {
		return handleSignal(ctx.Done(), cancel)
	})

	log.Println("server run at: ", addr)
	err := group.Wait()
	log.Println("sever recv error", err)
	time.Sleep(2 * time.Second) // 等待安全退出
}
