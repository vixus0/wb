package main

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/vixus0/wb/util"
	"github.com/vixus0/wb/wbd"
)

var sessionKey string

func cleanup() {
	log.Println("Cleaning up")
	exec.Command("bw", "lock").Run()
	os.Remove(wbd.Sock)
	os.Exit(0)
}

func main() {
	log.SetPrefix("[wbd] ")
	log.SetFlags(0)

	sessionKey := util.Input()

	if len(sessionKey) == 0 {
		log.Fatal("Empty input")
	}

	log.Printf("Got session key: %v...", util.Trunc(sessionKey, 8))

	serverStop := make(chan bool)
	server := wbd.NewServer(sessionKey, &serverStop)

	// Deal with interrupt signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		cleanup()
	}()

	// Die after lock delay
	go func() {
		time.Sleep(wbd.LockDelay)
		cleanup()
	}()

	// Die on RPC request
	go func() {
		<-serverStop
		cleanup()
	}()

	listener, err := net.Listen("unix", wbd.Sock)
	util.Err("listener error:", err)

	rpc.Register(server)
	rpc.Accept(listener)

	log.Println("started")
}
