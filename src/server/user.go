package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/db"
)

var (
	hasAccountCookieName  = "has_account"
	hasAccountCookieValue = "true"
)

func (s *Server) loadBannedUsers() error {
	users, err := s.store.GetBannedUsers(context.Background())
	if err != nil {
		return err
	}
	s.banned = make(map[string]struct{}, len(users))
	for _, u := range users {
		s.banned[u] = struct{}{}
	}
	return nil
}

func (s *Server) isUserBanned(id string) bool {
	_, found := s.banned[id]
	return found
}

func (s *Server) userAccount(c echo.Context) error {
	hc := s.buildPageContext(c, "My account")
	hc.NoBoost = true
	if hc.UserLoggedIn {
		if err := s.store.RunTx(c, func(ctx context.Context, q db.Querier) error {
			info, err := q.GetListForUser(ctx, hc.UserID)
			switch err {
			case nil:
				hc.Add("filter_count", info.InstanceCount)
				hc.Add("list_downloaded", info.Downloaded)
				hc.Add("list_token", info.Token.String())
				return nil
			case db.NotFound:
				token, err := q.CreateListForUser(ctx, hc.UserID)
				hc.Add("filter_count", 0)
				hc.Add("list_downloaded", false)
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

func (s *Server) rotateListToken(c echo.Context) error {
	if err := s.store.RunTx(c, func(ctx context.Context, q db.Querier) error {
		if c.Request().Method != http.MethodPost {
			return nil
		}
		formParams, err := c.FormParams()
		if err != nil {
			return err
		}
		confirmation := formParams.Get("confirm")
		tokenString := formParams.Get("token")
		token, err := uuid.Parse(tokenString)
		u := getUser(c)
		if !u.IsActive() || err != nil || confirmation != "on" {
			return errors.New("invalid arguments")
		}
		return q.RotateListToken(ctx, db.RotateListTokenParams{
			UserID: u.Id(),
			Token:  token,
		})
	}); err != nil {
		return err
	}

	return s.redirect(c, http.StatusSeeOther, s.echo.Reverse("user-account"))
}
