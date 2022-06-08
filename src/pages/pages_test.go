package pages

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func nilHandler(_ echo.Context) error { return nil }

func TestRedirect(t *testing.T) {
	e := echo.New()
	e.GET("/filters", nilHandler).Name = "list-filters"
	e.GET("/filters/tag/:tag", nilHandler).Name = "filters-for-tag"
	pages := &Pages{}

	for name, asserts := range map[string]func(t *testing.T, c echo.Context, rec *httptest.ResponseRecorder){
		"html no param": func(t *testing.T, c echo.Context, rec *httptest.ResponseRecorder) {
			assert.NoError(t, pages.RedirectToPage(c, "list-filters"))
			assert.Equal(t, 302, rec.Code, rec.Body)
			assert.Equal(t, "/filters", rec.Header().Get("Location"))
		},
		"html params": func(t *testing.T, c echo.Context, rec *httptest.ResponseRecorder) {
			assert.NoError(t, pages.RedirectToPage(c, "filters-for-tag", "test-tag"))
			assert.Equal(t, 302, rec.Code, rec.Body)
			assert.Equal(t, "/filters/tag/test-tag", rec.Header().Get("Location"))
		},
		"html 303": func(t *testing.T, c echo.Context, rec *httptest.ResponseRecorder) {
			assert.NoError(t, pages.Redirect(c, http.StatusSeeOther, "/news"))
			assert.Equal(t, 303, rec.Code, rec.Body)
			assert.Equal(t, "/news", rec.Header().Get("Location"))
		},
		"htmx no param": func(t *testing.T, c echo.Context, rec *httptest.ResponseRecorder) {
			c.Request().Header.Set("HX-Request", "true")
			assert.NoError(t, pages.RedirectToPage(c, "list-filters"))
			assert.Equal(t, 200, rec.Code, rec.Body)
			assert.Equal(t, "/filters", rec.Header().Get("HX-Redirect"))
		},
	} {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			asserts(t, c, rec)
		})
	}
}
