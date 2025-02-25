extends Control

@onready var Message: VBoxContainer = $vbox/message
@onready var Message_edit: LineEdit = $vbox/hbox/MessageEdit

func _ready() -> void:
	Host.NewMessage.connect(new_message)


func new_message(_message: String) -> void:
	print(_message)
	var rich = RichTextLabel.new()
	
	rich.scroll_active = false
	rich.bbcode_enabled = true
	rich.fit_content = true
	rich.text = _message
	
	Message.add_child(rich)


func _on_send_pressed() -> void:
	if Message_edit.text.length() > 0:
		Host.send_message(Message_edit.text)
