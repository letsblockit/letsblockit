package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/letsblockit/letsblockit/src/filters"
	"github.com/letsblockit/letsblockit/src/pages"
	"github.com/stretchr/testify/require"
)

func (s *ServerTestSuite) TestHelpMain_OK() {
	req := httptest.NewRequest(http.MethodGet, "/help", nil)
	s.expectRender("help-main", pages.ContextData{
		"page":          &mainPageDescription,
		"menu_sections": helpMenu,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestHelpUseList_OK() {
	token, err := s.store.CreateListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)

	req := httptest.NewRequest(http.MethodGet, "http://myhost/help/use-list", nil)
	s.expectRenderWithSidebar("help-use-list", "help-sidebar", pages.ContextData{
		"has_filters":   false,
		"list_url":      fmt.Sprintf("http://myhost/list/%s.txt", token.String()),
		"page":          helpMenu[0].Pages[0],
		"menu_sections": helpMenu,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestHelpUseList_DownloadDomainOK() {
	token, err := s.store.CreateListForUser(context.Background(), s.user)
	require.NoError(s.T(), err)
	require.NoError(s.T(), s.server.upsertFilterParams(s.c, s.user, &filters.Instance{Template: "one"}))

	req := httptest.NewRequest(http.MethodGet, "https://myhost/help/use-list", nil)
	s.server.options.ListDownloadDomain = "get.letsblock.it"
	s.expectRenderWithSidebar("help-use-list", "help-sidebar", pages.ContextData{
		"has_filters":   true,
		"list_url":      fmt.Sprintf("https://get.letsblock.it/list/%s.txt", token.String()),
		"page":          helpMenu[0].Pages[0],
		"menu_sections": helpMenu,
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestHelpFallback_OK() {
	req := httptest.NewRequest(http.MethodGet, "/help/invalid", nil)
	s.expectRender("help-main", pages.ContextData{
		"page":          &mainPageDescription,
		"menu_sections": helpMenu,
	})
	s.runRequest(req, assertOk)
}
