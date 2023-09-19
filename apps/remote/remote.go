package remote

type radioTransmiter interface {
	Start()
	Transmit(payload []byte) error
	PayloadSize() int
}

type remote struct {
	transmitter radioTransmiter
}

type RemoteCongigs struct {
	Transmitter radioTransmiter
}

func NewRemote(configs RemoteCongigs) *remote {
	return &remote{
		transmitter: configs.Transmitter,
	}
}

func (r *remote) Start() {

}
