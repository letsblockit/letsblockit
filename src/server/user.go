package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/models"
)

func (s *Server) userLogin(c echo.Context) error {
	if getUser(c) != nil {
		return c.Redirect(http.StatusFound, "/user/account")
	}
	hc := s.buildHandlebarsContext(c, "Login")
	return s.pages.render(c, "user-login", hc)
}

func (s *Server) userLogout(c echo.Context) error {
	if getUser(c) == nil {
		return c.Redirect(http.StatusFound, "/user/login")
	}
	logout, err := getLogoutUrl(s.options.OryProject, c)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusFound, logout)
}

func (s *Server) userAccount(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		// Cannot find user session, redirect to login
		return c.Redirect(http.StatusFound, "/user/login")
	}

	hc := s.buildHandlebarsContext(c, "My account")
	if user.IsVerified() {
		// Retrieve filter lists, create one if none exists
		var lists []models.FilterList
		s.gorm.Where("user_id = ?", user.Id()).Find(&lists)
		if len(lists) == 0 {
			lists = append(lists, models.FilterList{
				UserID: user.Id(),
				Name:   "My filters",
				Token:  uuid.NewString(),
			})
			s.gorm.Create(&lists)
		}

		hc["filter_lists"] = lists
		hc["verified"] = true
	}
	return s.pages.render(c, "user-account", hc)
}
