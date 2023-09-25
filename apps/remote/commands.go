package remote

func (r *remoteControl) ReadCommands() {
	roll := r.roll.Read()
	pitch := r.pitch.Read()
	yaw := r.yaw.Read()
	throttle := r.throttle.Read()

	r.commands = commands{
		roll:     roll,
		pitch:    pitch,
		yaw:      yaw,
		throttle: throttle,
	}
}

func (r *remoteControl) ReadButtons() {
	for i, button := range r.buttons {
		pressed:=button.IsPressed()
		if pressed {
			r.buttonsPressed[i]=byte(1)
		} else {
			r.buttonsPressed[i]=byte(0)
		}
	}
}