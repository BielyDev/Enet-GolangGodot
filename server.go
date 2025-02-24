package main

import (
	"fmt"
	"github.com/codecat/go-enet"
)


type MESSAGE uint16
const(
	MESSAGE_RECEIVED = iota
	MESSAGE_SEND
	LOADER_MESSAGE
)


var enetConnection enet.Host
var allPeers []enet.Peer
var save_message []interface{}

func main() {

	new_host()

	for true{
		var packet enet.Event = enetConnection.Service(100)
		
		switch packet.GetType() {
			case enet.EventConnect:
				fmt.Print("New connection!\n")

				allPeers = append(allPeers, packet.GetPeer())
				loaderMessage(packet.GetPeer())
			case enet.EventDisconnect:
				fmt.Print("New disconnection\n")

				erasePeers(packet.GetPeer())
			case enet.EventReceive:
				received_packet(packet)


		}
	}
}

func received_packet(_packet enet.Event){
	var message []byte = _packet.GetPacket().GetData()
	
	switch message[0]{
		case MESSAGE_SEND:
			send_all_client(MESSAGE_RECEIVED,get_message(message))
	}

}

func loaderMessage(peer enet.Peer){
	for i := 0; i < len(save_message); i += 1{
		peer.SendBytes(save_message[i].([]byte), 0, enet.PacketFlagReliable)
	}
}

func send_all_client(_type_message uint8, _message []byte) {
	var finished_message []byte

	finished_message = append(finished_message, 28, 0, 0, 0)
	finished_message = append(finished_message, 2, 0, 0, 0)
	finished_message = append(finished_message, 2, 0, 0, 0)
	finished_message = append(finished_message, _type_message, 0, 0, 0)
	finished_message = append(finished_message, 4, 0, 0, 0)
	
	
	finished_message = append(finished_message, byte(len(_message)), 0, 0, 0)

	for i := 0; i < len(_message); i += 1{
		finished_message = append(finished_message, _message[i])
	}

	finished_message = append(finished_message, 0, 0)
	save_message = append(save_message, finished_message)

	for peer_number := 0; peer_number < len(allPeers); peer_number += 1{
		allPeers[peer_number].SendBytes( finished_message, 0, enet.PacketFlagReliable)
	}

}

func new_host() bool {
	enet.Initialize()
	
	host, err := enet.NewHost(enet.NewAddress("127.0.0.1",8300),5,0,0,0)
	
	if err != nil{
		fmt.Print("Deu erro ai meu nobre: ",err)
		return false
	}
	fmt.Print("Servidor iniciado!\n")

	enetConnection = host
	return true
}

func filter_message(_message[]byte) []interface{} {
	var return_message []interface{}

	return_message = append(return_message, _message[0])
	return_message = append(return_message,string(_message[1:]))

	return return_message
}

func get_message(_message[]byte) []byte {
	return _message[1:]
}


func erasePeers(_peer enet.Peer){
	saveAllPeers := allPeers
	allPeers = []enet.Peer{}

	for i := 0;i < len(saveAllPeers); i+=1{
		if _peer != saveAllPeers[i]{
			allPeers = append(allPeers ,saveAllPeers[i])
		}
	}
}
