package store

import (
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
	FilterListID uint   `gorm:"index:idx_filters_by_list"`
	UserID       string `gorm:"index:idx_filters_by_user_filter"`
	FilterName   string `gorm:"index:idx_filters_by_user_filter"`
	Params       JSONMap
}
