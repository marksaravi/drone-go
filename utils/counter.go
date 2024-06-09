package utils

type counter struct {
	limit   int
	counter int
}

func NewCounter(limit int) *counter {
	return &counter{
		limit:   limit,
		counter: 0,
	}
}

func (c *counter) Inc() bool {
	c.counter++
	if c.counter == c.limit {
		c.counter = 0
		return true
	}
	return false
}
