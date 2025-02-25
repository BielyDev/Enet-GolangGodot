extends Control

signal NewMessage(_message: String)

var enet: ENetConnection = ENetConnection.new()

enum MESSAGE {
	PROFILE_SEND,
	PROFILE_RECEIVED,
	PROFILE_LOGIN,
	PROFILE_REGISTER,
	PROFILE_INSUFFICIENT_CHARACTER,
	PROFILE_USERNAME_ERR,
	PROFILE_PASSWORD_ERR,
	
	MESSAGE_SEND,
	MESSAGE_RECEIVED,
	LOADER_MESSAGE,
}


func _ready() -> void:
	set_process(false)
	
	var err = enet.create_host()
	if err != OK:
		OS.alert("Deu errado em criar o host")
		return
	
	var connect = enet.connect_to_host("127.0.0.1",8300)
	if connect == null:
		OS.alert("Deu errado em conectar-se ao host")
	
	await get_tree().create_timer(1).timeout
	
	set_process(true)

func _process(delta: float) -> void:
	var event = enet.service()
	
	match event[0]:
		ENetConnection.EVENT_RECEIVE:
			
			var message: Array = get_message(event[1].get_packet())
			match message[0]:
				MESSAGE.MESSAGE_RECEIVED:
					NewMessage.emit(message[1])


func send_profile(username: String, password: String) -> void:
	var message_send = PackedByteArray()
	
	var username_size: int = 1 + username.length()
	var password_size: int = (1 + username_size) + password.length()
	
	message_send.append(MESSAGE.PROFILE_SEND)
	message_send.append(1 + 4)
	message_send.append(username_size + 4)
	message_send.append(username_size + 1 + 4)
	message_send.append(password_size + 4)
	
	add_string_bytes(message_send,username)
	add_string_bytes(message_send,password)
	
	enet.broadcast(0,message_send,0)
	enet.flush() 

func send_message(text: String) -> void:
	var message_send = PackedByteArray()
	
	message_send.append(MESSAGE.MESSAGE_SEND)
	add_string_bytes(message_send,text)
	
	enet.broadcast(0,var_to_bytes([MESSAGE.MESSAGE_SEND,text]),0)
	enet.flush() 

func add_string_bytes(packedByte: PackedByteArray, message: String) -> void:
	for i in message.to_utf8_buffer():
		packedByte.append(i)

func get_message(_packet: PackedByteArray) -> Array:
	var bytes: PackedByteArray = PackedByteArray()
	
	for byte in _packet:
		bytes.append(int(byte))
	
	return bytes_to_var(bytes)
