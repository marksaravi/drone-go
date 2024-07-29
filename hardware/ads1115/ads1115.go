package ads1115

type ads1115AtoD struct {

}

func NewADS1115() *ads1115AtoD {
	return &ads1115AtoD {

	}
}

func (d *ads1115AtoD) Read(channel int) int {
	return 0
}