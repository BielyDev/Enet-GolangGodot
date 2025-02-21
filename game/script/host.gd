extends Control

var enet: ENetConnection = ENetConnection.new()
@onready var Message_edit: LineEdit = $hbox/MessageEdit

enum MESSAGE {
	MESSAGE_REQUEST,
	MESSAGE_RECEIVED,
	MESSAGE_SEND,
}

func _ready() -> void:
	set_process(false)
	
	var err = enet.create_host()
	if err != OK:
		OS.alert("Deu errado em criar o host")
		return
	
	var connect = enet.connect_to_host("127.0.0.1",8000)
	if connect == null:
		OS.alert("Deu errado em conectar-se ao host")
	
	await get_tree().create_timer(1).timeout
	
	set_process(true)

func _process(delta: float) -> void:
	enet.service()


func _on_send_button_pressed() -> void:
	if Message_edit.text.length() > 0:
		enet.broadcast(0,var_to_bytes([MESSAGE.MESSAGE_SEND, Message_edit.text]),0)
		enet.flush()
		Message_edit.text = ""


func _on_send_string_pressed() -> void:
	enet.broadcast(0,var_to_bytes(Message_edit.text),0)
	enet.flush()
