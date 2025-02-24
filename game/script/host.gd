extends Control


@onready var Message: VBoxContainer = $vbox/message
@onready var Message_edit: LineEdit = $vbox/hbox/MessageEdit

var enet: ENetConnection = ENetConnection.new()

enum MESSAGE {
	MESSAGE_RECEIVED,
	MESSAGE_SEND,
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
			var message: Array = bytes_to_var(event[1].get_packet())
			match message[0]:
				MESSAGE.MESSAGE_RECEIVED:
					var rich = RichTextLabel.new()
					rich.scroll_active = false
					rich.bbcode_enabled = true
					rich.fit_content = true
					rich.text = message[1]
					
					Message.add_child(rich)


func _on_send_string_pressed() -> void:
	var message_send = PackedByteArray()
	if Message_edit.text.length() > 0:
		message_send.append(MESSAGE.MESSAGE_SEND)
		add_string_bytes(message_send,Message_edit.text)
		
		enet.broadcast(0,message_send,0)
		enet.flush() 

func add_string_bytes(packedByte: PackedByteArray, message: String) -> void:
	for i in message.to_utf8_buffer():
		packedByte.append(i)
