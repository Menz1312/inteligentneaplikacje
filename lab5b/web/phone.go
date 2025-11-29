package web

type Phone struct {
	Company   string
	Model     string
	Weight    float64
	RAM       float64
	CamFront  float64
	CamBack   float64
	Processor string `json:"-"`
	Battery   float64
	Screen    float64
	PriceUSD  float64 `json:"-"`
	PricePLN  float64 `json:"Price"`
	Year      float64 `json:"-"`
	Value     float64
}
