package server

import (
	"sync"

	"github.com/labstack/echo/v4"
)

type helpSection struct {
	Title string
	Pages []*helpPageDescription
}

type helpPageDescription struct {
	Code  string
	Title string
}

var mainPageDescription = helpPageDescription{
	Code:  "main",
	Title: "Help center",
}

var helpMenu = []helpSection{{
	Title: "Usage",
	Pages: []*helpPageDescription{{
		Code:  "use-list",
		Title: "Use your filter list with uBlock",
	}, {
		Code:  "refresh-list",
		Title: "Manually refresh the filters",
	}, {
		Code:  "remove-list",
		Title: "Remove letsblock.it filters from uBlock",
	}},
}, {
	Title: "About",
	Pages: []*helpPageDescription{{
		Code:  "privacy",
		Title: "Privacy: where is my data?",
	}},
}}

// Index built on the first page load, protected by a sync.Once
var helpPageIndex map[string]*helpPageDescription
var helpPageIndexOnce sync.Once

func (s *Server) helpPages(c echo.Context) error {
	helpPageIndexOnce.Do(func() {
		helpPageIndex = make(map[string]*helpPageDescription)
		for _, s := range helpMenu {
			for _, p := range s.Pages {
				helpPageIndex[p.Code] = p
			}
		}
	})

	page := helpPageIndex[c.Param("page")]
	if page == nil {
		page = &mainPageDescription
	}

	hc := s.buildPageContext(c, page.Title)
	hc.Add("page", page)
	hc.Add("menu_sections", helpMenu)

	if page.Code == "use-list" && hc.UserLoggedIn {
		info, err := s.store.GetListForUser(c.Request().Context(), hc.UserID)
		if err == nil {
			hc.Add("has_filters", info.InstanceCount > 0)
			hc.Add("list_token", info.Token.String())
		}
	}

	if page.Code == mainPageDescription.Code {
		return s.pages.Render(c, "help-"+page.Code, hc)
	} else {
		return s.pages.RenderWithSidebar(c, "help-"+page.Code, "help-sidebar", hc)
	}
}
