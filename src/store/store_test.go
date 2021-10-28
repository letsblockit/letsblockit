package store

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	user  string
	store *Store
}

func (s *StoreTestSuite) SetupTest() {
	var err error
	s.user = uuid.NewString()
	s.store, err = NewMemStore()
	s.Require().NoError(err)
}

func (s *StoreTestSuite) TestDefaultState() {
	count, err := s.store.CountFilters(s.user)
	s.NoError(err)
	s.Equal(int64(0), count)
	names := s.store.GetActiveFilterNames(s.user)
	s.Empty(names)
}

func (s *StoreTestSuite) TestFilterList() {
	list, err := s.store.GetOrCreateFilterList(s.user)
	s.NoError(err)
	s.NotEmpty(list.ID)
	s.Equal(s.user, list.UserID)

	dupe, err := s.store.GetOrCreateFilterList(s.user)
	s.NoError(err)
	s.Equal(list.ID, dupe.ID)

	byToken, err := s.store.GetListForToken(list.Token)
	s.NoError(err)
	s.Equal(list.ID, byToken.ID)
}

func (s *StoreTestSuite) TestFilterInstances() {
	// Insert two filter instances and count them
	s.NoError(s.store.UpsertFilterInstance(s.user, "first", nil))
	s.NoError(s.store.UpsertFilterInstance(s.user, "second", map[string]interface{}{"a": true, "b": false}))
	count, err := s.store.CountFilters(s.user)
	s.NoError(err)
	s.Equal(int64(2), count)
	names := s.store.GetActiveFilterNames(s.user)
	s.EqualValues(map[string]bool{"first": true, "second": true}, names)

	// Get the filter list and one of the filters
	list, err := s.store.GetOrCreateFilterList(s.user)
	s.NoError(err)
	filter, err := s.store.GetFilterInstance(s.user, "second")
	s.NoError(err)
	s.EqualValues(map[string]interface{}{"a": true, "b": false}, filter.Params)
	s.Equal(list.ID, filter.FilterListID)

	// Delete first filter, and count the remaining one
	s.NoError(s.store.DropFilterInstance(s.user, "second"))
	count, err = s.store.CountFilters(s.user)
	s.NoError(err)
	s.Equal(int64(1), count)
	names = s.store.GetActiveFilterNames(s.user)
	s.EqualValues(map[string]bool{"first": true}, names)
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}
