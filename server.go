package main

import (
	"fmt"
	"github.com/codecat/go-enet"
)


type MESSAGE uint8
const(
	PROFILE_SEND = 0
	PROFILE_RECEIVED = 1
	PROFILE_LOGIN
	PROFILE_REGISTER
	PROFILE_INSUFFICIENT_CHARACTER
	PROFILE_USERNAME_ERR
	PROFILE_PASSWORD_ERR
	
	MESSAGE_SEND = 7
	MESSAGE_RECEIVED
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
	var _message_byte []byte = _packet.GetPacket().GetData()
	var message []interface{} = GdDeserialize(_message_byte)

	switch message[0]{
		case uint8(MESSAGE_SEND):
			send_all_client([]interface{}{uint8(MESSAGE_RECEIVED), message[1]})
			//_packet.GetPeer().SendBytes(GdSerialize([]interface{}{uint8(MESSAGE_RECEIVED), "Oi"}),0,0)
		case PROFILE_SEND:
			fmt.Print(message,"\n")

			//var _type_message uint8
			//var _err error

			//_type_message, _err = deserialize_profile(message)

			//if _err != nil{
			//	packet.GetPeer().SendBytes()
			//}
		default:
			fmt.Println("O tipo", message[0],"não existe :(")
			
	}
}

/*
func deserialize_profile(_message []byte) (uint8, error) {
	//var real_message []interface{}

	//var _password_byte [2]byte = [2]byte{_message[3],_message[4]}
	//var _password_size uint8 = _password_byte[1]-_password_byte[0]

	var _username string
	var _err error
	
	_username, _err = verify_username(_message)

	if _err != nil{
		return _err
	}

	fmt.Print(_username,"\n");

}
)
func verify_username(_message []byte) (string, error) {
	var _username_byte [2]byte = [2]byte{_message[1],_message[2]}
	var _username_size uint8 = _username_byte[1]-_username_byte[0]

	if _username_size < 4{
		fmt.Print("Nome precisa de ao menos 4 caracteres!")
		return "", errors.New("Nome precisa de ao menos 4 caracteres!")
	}

	var _username string = string(_message[_message[1]:_message[2]])

	var character uint8
	for character = 0; character < _username_size; character += 1{
		if _username[character] == 32{
			fmt.Print("O nome não pode conter espaços!")
			return "", errors.New("O nome não pode conter espaços!")
		}
	}

	return _username, nil
}
*/
func loaderMessage(peer enet.Peer){
	for i := 0; i < len(save_message); i += 1{
		peer.SendBytes(save_message[i].([]byte), 0, enet.PacketFlagReliable)
	}
}


func GdDeserialize(_message_byte[]byte) ([]interface{}) {

	if _message_byte[0] != 28{
		return nil
	}

	var _serial_message []interface{}
	var _var_pass uint8 = uint8(len(_message_byte))

	var byte uint8
	for byte = 0; byte < _var_pass; byte += 4{
		var bytes = _message_byte[byte:byte+4]

		if byte > 4{
			switch bytes[0]{
				case 28:
					fmt.Println("Array")
				case 2:
					_serial_message = append(_serial_message, uint8( _message_byte[byte + 4]))
				case 4:
					_serial_message = append(_serial_message, string(_message_byte[byte + 8:byte + 8 + (_message_byte[byte + 4])]))
			}
		}

	}

	fmt.Println(_serial_message)
	return _serial_message
}

func GdSerialize(new_array[]interface{}) []byte {
	var _serial_message []byte

	_serial_message = append(_serial_message, 28, 0, 0, 0)
	_serial_message = append(_serial_message, byte(len(new_array)), 0, 0, 0)

	var _key uint8
	for _key = 0; _key < uint8(len(new_array)); _key++{
		switch _value := new_array[_key].(type){
			case uint8:
				_serial_message = append(_serial_message, 2, 0, 0, 0)
				_serial_message = append(_serial_message, byte(_value), 0, 0, 0)
			case string:
				_serial_message = append(_serial_message, 4, 0, 0, 0)
				_serial_message = append(_serial_message, byte(len(_value)), 0, 0, 0)
				
				_serial_message = append(_serial_message, string_to_bytes(_value)...)
		}

	}

	_serial_message = append(_serial_message, 0, 0)

	return _serial_message
}

func bytes_to_string(_message []byte) string {
	return string(_message)
}
func string_to_bytes(_message string) []byte {
	var _message_byte []byte

	var _character uint8
	for _character = 0; _character < uint8(len(_message)); _character++{
		_message_byte = append(_message_byte, byte(_message[_character]))
	}

	return _message_byte
}

func send_all_client(_message []interface{}) {	
	var _new_message []byte = GdSerialize(_message)

	var _peer_number uint8
	for _peer_number = 0; _peer_number < uint8(len(allPeers)); _peer_number ++{
		allPeers[_peer_number].SendBytes( _new_message, 0, enet.PacketFlagReliable)
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
