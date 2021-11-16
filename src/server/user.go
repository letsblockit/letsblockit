package server

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/db"
)

var (
	hasAccountCookieName  = "has_account"
	hasAccountCookieValue = "true"
)

func (s *Server) userAccount(c echo.Context) error {
	hc := s.buildPageContext(c, "My account")
	hc.NoBoost = true
	if hc.UserVerified {
		if err := s.store.RunTx(c, func(ctx context.Context, q db.Querier) error {
			info, err := q.GetListForUser(ctx, hc.UserID)
			switch err {
			case nil:
				hc.Add("filter_count", info.InstanceCount)
				hc.Add("list_token", info.Token.String())
				return nil
			case db.NotFound:
				token, err := q.CreateListForUser(ctx, hc.UserID)
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

	c.SetCookie(&http.Cookie{
		Name:     hasAccountCookieName,
		Value:    hasAccountCookieValue,
		Path:     "/",
		Expires:  time.Now().AddDate(10, 0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	return s.pages.Render(c, "user-account", hc)
}
