package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/models"
)

func (s *Server) userLogin(c echo.Context) error {
	if getUser(c) != nil {
		return s.redirect(c, "user-account")
	}
	hc := s.buildHandlebarsContext(c, "Login")
	return s.pages.render(c, "user-login", hc)
}

func (s *Server) userLogout(c echo.Context) error {
	if getUser(c) == nil {
		return s.redirect(c, "user-login")
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
		return s.redirect(c, "user-login")
	}

	hc := s.buildHandlebarsContext(c, "My account")
	if user.IsVerified() {
		// Retrieve or create filter list
		var list models.FilterList
		s.gorm.Where("user_id = ?", user.Id()).First(&list)
		if list.Token == "" {
			list = models.FilterList{
				UserID: user.Id(),
				Name:   "My filters",
				Token:  uuid.NewString(),
			}
			s.gorm.Create(&list)
		}

		hc["filter_list"] = &list
	}
	return s.pages.render(c, "user-account", hc)
}
