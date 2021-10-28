package store

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	InMemoryTarget = ":memory:"
	mainDBFile     = "main.db"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type Store struct {
	orm *gorm.DB
}

func NewStore(target string, migrations bool) (*Store, error) {
	if target != InMemoryTarget {
		if err := os.MkdirAll(target, 0700); err != nil {
			return nil, err
		}
		target = filepath.Join(target, mainDBFile)
	}

	db := sqlite.Open(target)
	orm, err := gorm.Open(db, &gorm.Config{
		PrepareStmt:                              true,
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}
	if migrations {
		err = orm.AutoMigrate(&FilterList{}, &FilterInstance{})
		if err != nil {
			return nil, err
		}
	}
	return &Store{orm: orm}, nil
}

func NewMemStore() (*Store, error) {
	return NewStore(InMemoryTarget, true)
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
