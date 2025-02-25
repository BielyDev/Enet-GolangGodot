extends Control

@onready var PasswordEdit: LineEdit = $Panel/Margin/Margin/vbox/hbox/vbox/Password/Edit
@onready var UsernameEdit: LineEdit = $Panel/Margin/Margin/vbox/hbox/vbox/Username/Edit
@onready var Warning: Label = $Panel/Margin/Margin/vbox/Warning

func _on_view_pressed() -> void:
	PasswordEdit.secret = !PasswordEdit.secret


func _on_send_pressed() -> void:
	Host.send_profile(UsernameEdit.text,PasswordEdit.text)
