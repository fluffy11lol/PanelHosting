package models

//message Tariff {
//string ID = 1;
//int64 SSD = 2;
//int64 CPU = 3;
//int64 RAM = 4;
//int64 Price = 5;
//}

type Tariff struct {
	ID    string `json:"id,omitempty" db:"id"`
	Name  string `json:"name,omitempty" db:"name"`
	SSD   int    `json:"ssd,omitempty" db:"ssd"`
	CPU   int    `json:"cpu,omitempty" db:"cpu"`
	RAM   int    `json:"ram,omitempty" db:"ram"`
	Price int    `json:"price,omitempty" db:"price"`
}
