package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

var stopProcess uint64 = 5

type Valores struct {
	Id    uint64
	Val   uint64
	state bool
}

/*
func ProcessStarter(index chan int, procesos *[]Valores) {
	for {
		time.Sleep(time.Millisecond * 500)
		select {
		case cmp := <-index:
			go StartProcess(&(*procesos)[cmp].Id, &(*procesos)[cmp].Val)
		default:
		}
	}

}
*/

func StartProcess(id *uint64, val *uint64, NewProcess chan Valores) {
	for (*id) != stopProcess {
		(*val) = (*val) + 1
		time.Sleep(time.Millisecond * 500)
		select {
		case tmp := <-NewProcess:
			go StartProcess(&tmp.Id, &tmp.Val, NewProcess)
		default:
		}
	}
}
func Printer(valores *[]Valores, process chan int) {
	//var id *uint64 = nil
	//var valor *uint64
	for {
		//id = nil
		fmt.Println("----------")
		for _, val := range *valores {
			/*
				if !val.state {
					id = &val.Id
					valor = &val.Val
				}
			*/
			fmt.Printf("id %d: %d\n", val.Id, val.Val)
		}
		time.Sleep(time.Millisecond * 500)
		/*
			if id != nil {
				go StartProcess(id, valor)
			}
		*/
	}
}

func servidor() {
	s, err := net.Listen("tcp", ":9999")
	process := make(chan int)
	NewProcess := make(chan Valores)
	if err != nil {
		fmt.Println(err)
		return
	}
	var procesos []Valores
	procesos = append(procesos, Valores{Id: 0, Val: 0, state: true})
	procesos = append(procesos, Valores{Id: 1, Val: 0, state: true})
	procesos = append(procesos, Valores{Id: 2, Val: 0, state: true})
	procesos = append(procesos, Valores{Id: 3, Val: 0, state: true})
	procesos = append(procesos, Valores{Id: 4, Val: 0, state: true})
	//go ProcessStarter(process, &procesos)
	for i := 0; i < 5; i++ {
		go StartProcess(&procesos[i].Id, &procesos[i].Val, NewProcess)
	}
	go Printer(&procesos, process)
	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		var msg string
		E := gob.NewDecoder(c).Decode(&msg)
		if E != nil {
			fmt.Println(E)
			return
		}
		switch msg {
		case "Recibir":
			proceso := Valores{0, 0, false}
			err := gob.NewDecoder(c).Decode(&proceso)
			if err != nil {
				fmt.Println(err)
				return
			} else {
				//fmt.Println(proceso)
				c.Close()
				procesos = append(procesos, proceso)

				//fmt.Println(procesos)
				NewProcess <- proceso
			}
		case "Mandar":
			err := gob.NewEncoder(c).Encode(procesos[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			stopProcess = procesos[0].Id
			copy(procesos[0:], procesos[0+1:])
			procesos = procesos[:len(procesos)-1]
			stopProcess = 5
			//fmt.Println(procesos)
			c.Close()
		}
	}
}

/*
func handleClient(c net.Conn, val *[]Valores) {
	proceso := Valores{0, 0}
	err := gob.NewDecoder(c).Decode(&proceso)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(proceso)
		c.Close()
			val = append(val, proceso)
			go ProcesoPrincipal(&val[len(val)-1])
	}
}*/

func main() {
	go servidor()
	var input string
	fmt.Scanln(&input)
}
