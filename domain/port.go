package domain

import (
	"encoding/json"
	"time"
)

/* NOTES: postgres */
type Port struct {
	ID          string          `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   *time.Time      `json:"deleted_at,omitempty" gorm:"index"`
	Name        string          `json:"name"`
	City        string          `json:"city"`
	Country     string          `json:"country"`
	Alias       json.RawMessage `json:"alias" gorm:"type:JSON"`
	Regions     json.RawMessage `json:"regions" gorm:"type:JSON"`
	Coordinates json.RawMessage `json:"coordinates" gorm:"type:JSON"` /* NOTES: postgis */
	Province    string          `json:"province"`
	Timezone    string          `json:"timezone"`
	Unlocs      json.RawMessage `json:"unlocs" gorm:"type:JSON"`
	Code        string          `json:"code"`
}
