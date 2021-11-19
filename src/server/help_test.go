package server

import (
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/xvello/letsblockit/src/db"
	"github.com/xvello/letsblockit/src/pages"
)

func (s *ServerTestSuite) TestHelpMain_OK() {
	req := httptest.NewRequest(http.MethodGet, "/help", nil)
	req.AddCookie(verifiedCookie)
	s.expectRender("help-main", pages.ContextData{
		"page":          &mainPageDescription,
		"menu_sections": helpMenu,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestHelpUseList_OK() {
	token := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/help/use-list", nil)
	req.AddCookie(verifiedCookie)
	s.expectQ.GetListForUser(gomock.Any(), s.user).Return(db.GetListForUserRow{
		Token:         token,
		InstanceCount: 5,
	}, nil)
	s.expectRenderWithSidebar("help-use-list", "help-sidebar", pages.ContextData{
		"has_filters":   true,
		"list_token":    token.String(),
		"page":          helpMenu[0].Pages[0],
		"menu_sections": helpMenu,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestHelpFallback_OK() {
	req := httptest.NewRequest(http.MethodGet, "/help/invalid", nil)
	req.AddCookie(verifiedCookie)
	s.expectRender("help-main", pages.ContextData{
		"page":          &mainPageDescription,
		"menu_sections": helpMenu,
	})
	s.runRequest(req, assertOk)
}
