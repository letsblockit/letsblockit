package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/db"
	"github.com/letsblockit/letsblockit/src/users/auth"
)

func (s *Server) userAccount(c echo.Context) error {
	hc := s.buildPageContext(c, "My account")
	hc.NoBoost = true
	if hc.UserLoggedIn {
		if err := s.store.RunTx(c, func(ctx context.Context, q db.Querier) error {
			info, err := q.GetListForUser(ctx, hc.UserID)
			switch err {
			case nil:
				hc.Add("filter_count", info.InstanceCount)
				hc.Add("list_downloaded", info.DownloadedAt.Valid)
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
	return s.pages.Render(c, "user-account", hc)
}

func (s *Server) updatePreferences(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return c.NoContent(http.StatusMethodNotAllowed)
	}
	user := auth.GetUserId(c)
	if user == "" {
		return errors.New("invalid user session")
	}
	formParams, err := c.FormParams()
	if err != nil {
		return err
	}
	if err := s.preferences.UpdatePreferences(c, db.UpdateUserPreferencesParams{
		UserID:       user,
		ColorMode:    db.ColorMode(formParams.Get("color_mode")),
		BetaFeatures: formParams.Get("beta_features") == "on",
	}); err != nil {
		return err
	}

	return s.pages.Redirect(c, http.StatusSeeOther, s.echo.Reverse("user-account"))
}

func (s *Server) rotateListToken(c echo.Context) error {
	if err := s.store.RunTx(c, func(ctx context.Context, q db.Querier) error {
		if c.Request().Method != http.MethodPost {
			return c.NoContent(http.StatusMethodNotAllowed)
		}
		formParams, err := c.FormParams()
		if err != nil {
			return err
		}
		confirmation := formParams.Get("confirm")
		tokenString := formParams.Get("token")
		token, err := uuid.Parse(tokenString)
		user := auth.GetUserId(c)
		if user == "" || err != nil || confirmation != "on" {
			return errors.New("invalid arguments")
		}
		return q.RotateListToken(ctx, db.RotateListTokenParams{
			UserID: user,
			Token:  token,
		})
	}); err != nil {
		return err
	}

	return s.pages.Redirect(c, http.StatusSeeOther, s.echo.Reverse("user-account"))
}
