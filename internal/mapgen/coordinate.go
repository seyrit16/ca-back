package mapgen

var Directions = [6][2]int16{
	{1, 0},
	{1, -1},
	{0, -1},
	{-1, 0},
	{-1, 1},
	{0, 1},
}

type Coordinate struct {
	Q int16
	R int16
}

var InvalidCoord = Coordinate{Q: -32768, R: -32768}

func (c Coordinate) Key() [2]int16 {
	return [2]int16{c.Q, c.R}
}

func (c Coordinate) IsInvalid() bool {
	return c.Q == -32768
}

func (c Coordinate) Equals(other Coordinate) bool {
	return c.Q == other.Q && c.R == other.R
}
