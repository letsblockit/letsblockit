package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/letsblockit/letsblockit/src/db"
	"github.com/letsblockit/letsblockit/src/pages"
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
	req := httptest.NewRequest(http.MethodGet, "http://myhost/help/use-list", nil)
	req.AddCookie(verifiedCookie)
	s.expectQ.GetListForUser(gomock.Any(), s.user).Return(db.GetListForUserRow{
		Token:         token,
		InstanceCount: 5,
	}, nil)
	s.expectRenderWithSidebar("help-use-list", "help-sidebar", pages.ContextData{
		"has_filters":   true,
		"list_url":      fmt.Sprintf("http://myhost/list/%s", token.String()),
		"page":          helpMenu[0].Pages[0],
		"menu_sections": helpMenu,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestHelpUseList_DownloadDomainOK() {
	token := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "https://myhost/help/use-list", nil)
	req.AddCookie(verifiedCookie)
	s.server.options.ListDownloadDomain = "get.letsblock.it"
	s.expectQ.GetListForUser(gomock.Any(), s.user).Return(db.GetListForUserRow{
		Token:         token,
		InstanceCount: 5,
	}, nil)
	s.expectRenderWithSidebar("help-use-list", "help-sidebar", pages.ContextData{
		"has_filters":   true,
		"list_url":      fmt.Sprintf("https://get.letsblock.it/list/%s", token.String()),
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
