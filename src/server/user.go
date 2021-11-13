package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/db"
)

func (s *Server) userLogin(c echo.Context) error {
	if getUser(c) != nil {
		return s.redirect(c, "user-account")
	}
	hc := s.buildPageContext(c, "Login")
	return s.pages.Render(c, "user-login", hc)
}

func (s *Server) userLogout(c echo.Context) error {
	if getUser(c) == nil {
		return s.redirect(c, "user-login")
	}
	logout, err := s.getLogoutUrl(c)
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

	hc := s.buildPageContext(c, "My account")
	if user.IsVerified() {
		if err := s.store.RunTx(c, func(ctx context.Context, q db.Querier) error {
			info, err := q.GetListForUser(ctx, user.Id())
			switch err {
			case nil:
				hc.Add("filter_count", info.InstanceCount)
				hc.Add("list_token", info.Token.String())
				return nil
			case db.NotFound:
				token, err := q.CreateListForUser(ctx, user.Id())
				hc.Add("filter_count", 0)
				hc.Add("list_token", token.String())
				return err
			default:
				return err
			}
		}); err != nil {
			return err
		}
	}
	return s.pages.Render(c, "user-account", hc)
}
