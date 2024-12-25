package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Project struct {
	ID             string         `json:"id,omitempty" db:"id"`
	UserID         string         `json:"user_id" db:"user_id"`
	Name           string         `json:"name" db:"name"`
	Networking     NetworkingInfo `json:"networking" gorm:"embedded"`       // структура для Address и Ports
	Status         ProjectStatus  `json:"status" db:"status"`               // тип для статуса проекта
	TariffID       int32          `json:"tariff_id" db:"tariff_id"`         // ID тарифа
	TariffStatus   TariffStatus   `json:"tariff_status" db:"tariff_status"` // структура для тарифа
	ExpirationTime time.Time      `json:"expiration_time" db:"expiration_time"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
}

type NetworkingInfo struct {
	Address string     `json:"address"`
	Ports   Int32Slice `json:"ports" gorm:"type:json"`
}

type ProjectStatus int

const (
	ProjectStopped  ProjectStatus = 0 // Выключённый
	ProjectRunning  ProjectStatus = 1 // Включённый
	ProjectStarting ProjectStatus = 2 // Запускается
	ProjectStopping ProjectStatus = 3 // Останавливается
)

type TariffStatus int

const (
	TariffActive  TariffStatus = 0 // Активный
	TariffExpired TariffStatus = 1 // Истёк
)

type Int32Slice []int32

func (s Int32Slice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Int32Slice) Scan(value interface{}) error {
	bytes, ok := value.([]uint8)
	if !ok {
		return fmt.Errorf("cannot scan type %T into Int32Slice", value)
	}

	return json.Unmarshal(bytes, s)
}
