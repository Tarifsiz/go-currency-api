package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Currency represents a currency with its properties
type Currency struct {
	ID                  uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code                string    `json:"code" gorm:"type:varchar(3);unique;not null;index"`
	Description         string    `json:"description" gorm:"type:varchar(255);not null"`
	AmountDisplayFormat string    `json:"amount_display_format" gorm:"type:varchar(50);default:'###,###.##'"`
	HtmlEncodedSymbol   string    `json:"html_encoded_symbol" gorm:"type:varchar(50)"`
	Factor              int       `json:"factor" gorm:"default:100"` // For decimal precision (100 = 2 decimal places)
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy           uuid.UUID `json:"created_by" gorm:"type:uuid"`
}

// BeforeCreate hook for Currency
func (c *Currency) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// TableName method for explicit table naming
func (Currency) TableName() string {
	return "currencies"
}