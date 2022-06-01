package mocks

import (
	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/db"
)

type MockStore struct {
	*MockQuerier
}

func NewMockStore(q *MockQuerier) *MockStore {
	return &MockStore{q}
}

func (m MockStore) RunTx(e echo.Context, f db.TxFunc) error {
	return f(e.Request().Context(), m)
}
