package B0116

type OcrData struct {
	Images []Images `json:"images" form:"images" gorm:"工号"`
}

type Images struct {
	Pages []Pages `json:"pages" form:"pages" gorm:"工号"`
}

type Pages struct {
	Blocks []Blocks `json:"blocks" form:"blocks" gorm:"工号"`
}

type Blocks struct {
	Area []float64 `json:"area" form:"area" gorm:"工号"`
	Text string    `json:"text" form:"text" gorm:"工号"`
	N    int       `json:"n" form:"n" gorm:"工号"`
}
