package modelsControlPanel

import "time"

type Server struct {
	ID           string    `gorm:"primaryKey;type:varchar(255);not null"`
	Name         string    `gorm:"not null"`
	Status       string    `gorm:"not null"`
	TariffStatus string    `gorm:"not null"`
	Address      string    `gorm:"not null"`
	Port         string    `gorm:"not null"`
	UserID       uint      `json:"user_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	TariffID     int       `json:"tariff_id" gorm:"tariff_id"`
}
