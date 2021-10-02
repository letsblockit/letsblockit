package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type FilterList struct {
	gorm.Model
	UserID          string `gorm:"index:idx_lists_by_user"`
	Name            string
	Token           string `gorm:"uniqueIndex:idx_lists_by_token"`
	FilterInstances []*FilterInstance
}

type FilterInstance struct {
	gorm.Model
	FilterListID uint         `gorm:"index:idx_filters_by_list"`
	UserID       string       `gorm:"index:idx_filters_by_user_filter"`
	FilterName   string       `gorm:"index:idx_filters_by_user_filter"`
	Params       FilterParams `gorm:"type:bytes"`
}

type FilterParams map[string]interface{}

// Scan scans value into a new map, implements sql.Scanner interface
func (p *FilterParams) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("not valid JSON: %s", value)
	}
	*p = make(map[string]interface{})
	return json.Unmarshal(bytes, p)
}

// Value return json value, implement driver.Valuer interface
func (p FilterParams) Value() (driver.Value, error) {
	if len(p) == 0 {
		return nil, nil
	}
	return json.Marshal(p)
}
