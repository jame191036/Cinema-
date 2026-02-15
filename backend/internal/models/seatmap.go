package models

type Seat struct {
	SeatCode string `bson:"seat_code" json:"seatCode"`
	Type     string `bson:"type" json:"type"`
	Active   bool   `bson:"active" json:"active"`
}

type Row struct {
	RowLabel string `bson:"row_label" json:"rowLabel"`
	Seats    []Seat `bson:"seats" json:"seats"`
}

type Seatmap struct {
	ID   string `bson:"_id" json:"id"`
	Rows []Row  `bson:"rows" json:"rows"`
}
