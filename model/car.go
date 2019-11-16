package model

type CarType uint

const (
	UnknownCarType CarType = iota
	Hatchback
	Sedan
	SUV
)

func (c *Car) EnumToStr() {
	switch c.Car {
	case Hatchback:
		c.CarType = "Hatchback"
	case Sedan:
		c.CarType = "Sedan"
	case SUV:
		c.CarType = "SUV"
	default:
		c.CarType = "Unknown"
	}
}

func (c *Car) StrToEnum() {
	switch c.CarType {
	case "Hatchback":
		c.Car = Hatchback
	case "Sedan":
		c.Car = Sedan
	case "SUV":
		c.Car = SUV
	default:
		c.Car = UnknownCarType
	}
}
