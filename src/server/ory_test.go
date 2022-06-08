package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/pages"
	"github.com/letsblockit/letsblockit/src/users/auth"
	"github.com/stretchr/testify/assert"
)

// TODO: rewrite these tests out of the main suite and into the users/auth package

func (s *ServerTestSuite) TestRenderKratosForm_OK() {
	req := httptest.NewRequest(http.MethodGet, "/user/forms/login?flow=123456", nil)
	s.expectRender("kratos-form", pages.ContextData{
		"type": "login",
		"ui": map[string]interface{}{
			"a": "1",
			"b": "2",
		},
		"return_to": "https://target",
		"settings":  auth.SupportedForms["login"],
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestRenderKratosForm_ErrFormType() {
	req := httptest.NewRequest(http.MethodGet, "/user/forms/unknown?flow=123456", nil)
	s.expectP.Redirect(gomock.Any(), 302, s.kratosServer.URL+"/ui/unknown?flow=123456")
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestRenderKratosForm_ErrBadFlow() {
	req := httptest.NewRequest(http.MethodGet, "/user/forms/login?flow=666", nil)
	s.expectP.Redirect(gomock.Any(), 302, s.kratosServer.URL+"/ui/login?flow=666")
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestStartKratosFlow_Settings() {
	req := httptest.NewRequest(http.MethodPost, "/user/action/settings", s.csrfBody())
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/settings/browser")
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestStartKratosFlow_LoginOrRegistration_Register() {
	req := httptest.NewRequest(http.MethodPost, "/user/action/loginOrRegistration", s.csrfBody())
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/registration/browser")
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestStartKratosFlow_LoginOrRegistration_Login() {
	req := httptest.NewRequest(http.MethodPost, "/user/action/loginOrRegistration", s.csrfBody())
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(&http.Cookie{
		Name:  "has_account",
		Value: "true",
	})
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/login/browser")
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestStartKratosFlow_Login() {
	req := httptest.NewRequest(http.MethodPost, "/user/action/login", s.csrfBody())
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/login/browser")
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestStartKratosFlow_Login_ReturnToFromForm() {
	form := make(url.Values)
	form.Add(csrfLookup, s.csrf)
	form.Add("return_to", "https://myserver/page")
	req := httptest.NewRequest(http.MethodPost, "https://myserver/user/action/login",
		strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set("Referer", "https://myserver/ignore")
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/login/browser?return_to=https://myserver/page")
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestStartKratosFlow_Login_ReturnToFromReferer() {
	req := httptest.NewRequest(http.MethodPost, "https://myserver/user/action/login", s.csrfBody())
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set("Referer", "https://myserver/page")
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/login/browser?return_to=https://myserver/page")
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestStartKratosFlow_Login_ReturnToNotInDomain() {
	req := httptest.NewRequest(http.MethodPost, "https://myserver/user/action/login", s.csrfBody())
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set("Referer", "https://anotherserver/page")
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/login/browser")
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestStartKratosFlow_Logout() {
	req := httptest.NewRequest(http.MethodPost, "/user/action/logout", s.csrfBody())
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(verifiedCookie)
	s.expectP.Redirect(gomock.Any(), 303, "targetURL")
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestStartKratosFlow_MissingCSRF() {
	req := httptest.NewRequest(http.MethodPost, "/user/action/logout", nil)
	req.AddCookie(verifiedCookie)
	s.runRequest(req, func(t *testing.T, recorder *httptest.ResponseRecorder) {
		assert.Equal(t, 400, recorder.Result().StatusCode)
	})
}

func (s *ServerTestSuite) csrfBody() io.Reader {
	f := url.Values{}
	f.Add(csrfLookup, s.csrf)
	return strings.NewReader(f.Encode())
}
