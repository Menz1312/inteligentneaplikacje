package apts

type Apartment struct {
	Year     int
	Age      int
	Area     int
	Floor    int
	Parking  int
	Bus      int
	Metro    int
	Location int
	Parks    int
	Schools  int
	Price    int
}

func (a Apartment) GetData() ([]float64, float64) {
	return []float64{
		float64(a.Year),
		float64(a.Age),
		float64(a.Area),
		float64(a.Floor),
		float64(a.Parking),
		float64(a.Bus),
		float64(a.Metro),
		float64(a.Location),
		float64(a.Parks),
		float64(a.Schools),
	}, float64(a.Price)
}
