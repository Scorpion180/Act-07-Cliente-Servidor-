package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Valores struct {
	Id    uint64
	Val   uint64
	state bool
}

var sendToHost = make(chan bool)

func SetupCloseHandler(proceso *Valores) {
	closeChannel := make(chan os.Signal)
	signal.Notify(closeChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-closeChannel
		sendToHost <- true
		time.Sleep(time.Millisecond * 500)
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}
func ProcesoPrincipal(val *Valores) {
	for {
		select {
		default:
			fmt.Printf("id %d: %d\n", val.Id, val.Val)
			val.Val = val.Val + 1
		}
		time.Sleep(time.Millisecond * 500)
		select {
		case <-sendToHost:
			sendProcessToHost(val)
		default:
		}
	}
}
func cliente(proceso *Valores, c net.Conn, err error) {
	err = gob.NewEncoder(c).Encode("Mandar")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewDecoder(c).Decode(proceso)
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Close()
	go ProcesoPrincipal(proceso)
}
func sendProcessToHost(proceso *Valores) {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode("Recibir")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(proceso)
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Millisecond * 500)
}
func main() {
	var proceso Valores
	proceso = Valores{0, 0, false}
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	SetupCloseHandler(&proceso)
	go cliente(&proceso, c, err)
	for true {

	}
}
