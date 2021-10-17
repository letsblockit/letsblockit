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
		hc["filter_list"] = s.getOrCreateFilterList(user)
	}
	return s.pages.render(c, "user-account", hc)
}

func (s *Server) getOrCreateFilterList(user *oryUser) *models.FilterList {
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
	return &list
}

func (s *Server) getActiveFilterNames(user *oryUser) map[string]bool {
	var names []string
	s.gorm.Model(&models.FilterInstance{}).Where("user_id = ?", user.Id()).
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