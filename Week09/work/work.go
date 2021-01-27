package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

/*
1. 用 Go 实现一个 tcp server ，用两个 goroutine 读写 conn，两个 goroutine 通过 chan 可以传递 message，能够正确退出
以上作业，要求提交到 GitHub 上面，Week09 作业地址：
*/

func main() { // errgroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	wg, ctx := errgroup.WithContext(ctx)

	// server tcp
	wg.Go(func() error {
		return serverTcp(ctx)
	})

	// listen signal
	wg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case sign := <-exitChan:
			cancel()
			return fmt.Errorf("recv signal: %v", sign)
		}
	})

	if err := wg.Wait(); err != nil {
		log.Println("recv error: ", err)
	}

	log.Println("server stop")

}

func serverTcp(ctx context.Context) error {
	lister, err := net.Listen("tcp", ":8888")
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return lister.Close()

		default:
			tcpConn, err := lister.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}

			NewTcpHandle(tcpConn, 100).handle(ctx) // errgroup
		}
	}
}

func NewTcpHandle(conn net.Conn, caps int) *tcpHandle {
	return &tcpHandle{conn: conn, buf: make(chan []byte, caps)}
}

type tcpHandle struct {
	conn net.Conn
	buf  chan []byte
}

func (h *tcpHandle) handle(ctx context.Context) error {
	defer h.conn.Close()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return h.read(ctx)
	})

	eg.Go(func() error {
		return h.write(ctx)
	})

	return eg.Wait()
}

func (h *tcpHandle) read(ctx context.Context) error {
	reader := bufio.NewReader(h.conn)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			d, _, err := reader.ReadLine()
			if err != nil {

			}

			h.buf <- d
		}
	}
}

func (h *tcpHandle) write(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil

		case data := <-h.buf:
			h.conn.Write(data)
		}
	}
}
