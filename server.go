package main

import (
	"fmt"
	"github.com/codecat/go-enet"
	"encoding/binary"
	"bytes"
)

var enetConnection enet.Host

type MESSAGE int16
const(
	MESSAGE_REQUEST = iota
	MESSAGE_RECEIVED = 1
	MESSAGE_SEND = 2
)


func main() {

	var err = new_host()

	for err{
		var events enet.Event

		events = enetConnection.Service(100)
		
		switch events.GetType() {
			case enet.EventConnect:
				fmt.Print("OIIIIIIIIII\n")
			case enet.EventDisconnect:
				fmt.Print("BROXOU\n")
			case enet.EventReceive:
				var message_type int32
				var message_string string
				//var message string
				
				packet := events.GetPacket().GetData()
				buf := bytes.NewReader(packet)
				binary.Read(buf ,binary.LittleEndian, &message_type)
				binary.Read(buf ,binary.LittleEndian, &message_string)
				//binary.Read(perae ,binary.LittleEndian, &message)

				fmt.Print(buf,"\n")
				fmt.Print(message_type,"\n")
				fmt.Print(message_string,"\n")
				//fmt.Print(message,"\n")
		}
	}
}

func new_host() bool {
	enet.Initialize()
	
	host, err := enet.NewHost(enet.NewAddress("127.0.0.1",8000),5,0,0,0)
	
	if err != nil{
		fmt.Print("\nDeu erro ai meu nobre: ",err)
		return false
	}

	enetConnection = host

	return true
}
