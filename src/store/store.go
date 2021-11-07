package store

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Option string

const (
	DropOwned Option = "drop_owned"
	Migrate   Option = "migrate_schema"
)

type Store struct {
	orm *gorm.DB
}

var ErrRecordNotFound = gorm.ErrRecordNotFound

func NewStore(host, database string, options ...Option) (*Store, error) {
	db := postgres.New(postgres.Config{
		DSN:                  fmt.Sprintf("postgresql:///%s?host=%s", database, host),
		PreferSimpleProtocol: true,
	})
	orm, err := gorm.Open(db, &gorm.Config{
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
	})
	if err != nil {
		return nil, err
	}
	for _, o := range options {
		switch o {
		case DropOwned:
			if err := orm.Exec("drop owned by current_user").Error; err != nil {
				return nil, err
			}
		case Migrate:
			if err := orm.AutoMigrate(&FilterList{}, &FilterInstance{}); err != nil {
				return nil, err
			}
		}
	}

	return &Store{orm: orm}, nil
}

func (s *Store) CountFilters(user string) (int64, error) {
	var filterCount int64
	err := s.orm.Model(&FilterInstance{}).Where("user_id = ?", user).Count(&filterCount).Error
	return filterCount, err
}

func (s *Store) GetActiveFilterNames(user string) map[string]bool {
	if user == "" {
		return nil
	}
	var names []string
	s.orm.Model(&FilterInstance{}).Where("user_id = ?", user).
		Distinct().Pluck("FilterName", &names)
	if len(names) == 0 {
		return nil
	}

	out := make(map[string]bool)
	for _, n := range names {
		out[n] = true
	}
	return out
}

func (s *Store) GetOrCreateFilterList(user string) (*FilterList, error) {
	var list FilterList
	err := s.orm.Where("user_id = ?", user).First(&list).Error

	switch err {
	case nil:
		return &list, nil
	case gorm.ErrRecordNotFound:
		list = FilterList{
			UserID: user,
			Name:   "My filters",
			Token:  uuid.NewString(),
		}
		return &list, s.orm.Create(&list).Error
	default:
		return nil, err
	}
}

func (s *Store) GetListForToken(token string) (*FilterList, error) {
	list := FilterList{
		Token: token,
	}
	err := s.orm.Where(&list).Preload("FilterInstances").First(&list).Error
	return &list, err
}

func (s *Store) UpsertFilterInstance(user, filterName string, params JSONMap) error {
	f := &FilterInstance{
		UserID:     user,
		FilterName: filterName,
	}
	err := s.orm.Where(f).First(f).Error
	f.Params = params

	switch err {
	case nil:
		return s.orm.Save(&f).Error
	case gorm.ErrRecordNotFound:
		list, err := s.GetOrCreateFilterList(user)
		if err != nil {
			return nil
		}
		f.FilterListID = list.ID
		return s.orm.Create(&f).Error
	default:
		return err
	}
}

func (s *Store) DropFilterInstance(user string, filterName string) error {
	target := &FilterInstance{
		UserID:     user,
		FilterName: filterName,
	}
	return s.orm.Where(target).Delete(target).Error
}

func (s *Store) GetFilterInstance(user string, filterName string) (*FilterInstance, error) {
	f := &FilterInstance{
		UserID:     user,
		FilterName: filterName,
	}
	e := s.orm.Where(f).First(f)
	return f, e.Error
}
