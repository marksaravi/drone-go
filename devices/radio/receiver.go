package radio

import "github.com/marksaravi/drone-go/constants"

type radioReceiverLink interface {
	PowerOn()
	ReceiverOn()
	Listen()
	Receive() ([]byte, error)
	IsReceiverDataReady(update bool) bool
}

type radioReceiver struct {
	radiolink radioReceiverLink
}

func NewRadioReceiver(radiolink radioReceiverLink) *radioReceiver {
	return &radioReceiver{
		radiolink: radiolink,
	}
}

func (r *radioReceiver) On() {
	r.radiolink.ReceiverOn()
	r.radiolink.PowerOn()
	r.radiolink.Listen()
}

func (r *radioReceiver) Receive() ([]byte, bool) {
	if r.radiolink.IsReceiverDataReady(true) {
		data, err := r.radiolink.Receive()
		if err != nil || len(data) != int(constants.RADIO_PAYLOAD_SIZE) {
			return nil, false
		}
		r.radiolink.Listen()
		return data, true
	}
	return nil, false
}
