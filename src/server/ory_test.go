package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/xvello/letsblockit/src/pages"
)

func TestNilOryUser(t *testing.T) {
	var user *oryUser
	assert.False(t, user.IsActive())
	assert.EqualValues(t, uuid.Nil, user.Id())
}

func TestVerifiedOryUser(t *testing.T) {
	payload := `{
	  "id": "d631b403-eb29-4a5b-8829-125da6ebdf75",
	  "active": true,
	  "identity": {
		"id": "9a3f8aeb-729a-44cf-bede-f885175344ef",
		"verifiable_addresses": [
		  {
			"id": "0988fc40-3cb1-4174-b867-cac9de28f1a4",
			"value": "test@example.com",
			"verified": true
		  }
		]
	  }
	}`
	user := new(oryUser)
	assert.NoError(t, json.Unmarshal([]byte(payload), user))
	assert.True(t, user.IsActive())
	assert.Equal(t, "9a3f8aeb-729a-44cf-bede-f885175344ef", user.Id().String())
}

func TestInactiveOrySession(t *testing.T) {
	payload := `{
	  "id": "d631b403-eb29-4a5b-8829-125da6ebdf75",
	  "identity": {
		"id": "9a3f8aeb-729a-44cf-bede-f885175344ef",
		"verifiable_addresses": [
		  {
			"id": "0988fc40-3cb1-4174-b867-cac9de28f1a4",
			"value": "test@example.com",
			"verified": true
		  }
		]
	  }
	}`
	user := new(oryUser)
	assert.NoError(t, json.Unmarshal([]byte(payload), user))
	assert.False(t, user.IsActive())
}

func (s *ServerTestSuite) TestRenderKratosForm_OK() {
	req := httptest.NewRequest(http.MethodGet, "/user/forms/login?flow=123456", nil)
	s.expectRender("kratos-form", pages.ContextData{
		"type": "login",
		"ui": map[string]interface{}{
			"a": "1",
			"b": "2",
		},
		"return_to": "https://target",
		"settings":  supportedForms["login"],
	})
	s.runRequest(req, assertOk)
}

func (s *ServerTestSuite) TestRenderKratosForm_KratosDown() {
	s.kratosServer.Close() // Kratos is unresponsive, continue anonymous
	req := httptest.NewRequest(http.MethodGet, "/user/forms/login?flow=123456", nil)
	s.runRequest(req, assertRedirect("/.ory/ui/login?flow=123456"))
}

func (s *ServerTestSuite) TestRenderKratosForm_ErrFormType() {
	req := httptest.NewRequest(http.MethodGet, "/user/forms/unknown?flow=123456", nil)
	s.runRequest(req, assertRedirect("/.ory/ui/unknown?flow=123456"))
}

func (s *ServerTestSuite) TestRenderKratosForm_ErrBadFlow() {
	req := httptest.NewRequest(http.MethodGet, "/user/forms/login?flow=666", nil)
	s.runRequest(req, assertRedirect("/.ory/ui/login?flow=666"))
}

func (s *ServerTestSuite) TestStartKratosFlow_Settings() {
	req := httptest.NewRequest(http.MethodPost, "/user/start/settings", nil)
	s.runRequest(req, assertSeeOther(s.kratosServer.URL+"/api/kratos/public/self-service/settings/browser"))
}

func (s *ServerTestSuite) TestStartKratosFlow_LoginOrRegistration_Register() {
	req := httptest.NewRequest(http.MethodPost, "/user/start/loginOrRegistration", nil)
	s.runRequest(req, assertSeeOther(s.kratosServer.URL+"/api/kratos/public/self-service/registration/browser"))
}

func (s *ServerTestSuite) TestStartKratosFlow_LoginOrRegistration_Login() {
	req := httptest.NewRequest(http.MethodPost, "/user/start/loginOrRegistration", nil)
	req.AddCookie(&http.Cookie{
		Name:  "has_account",
		Value: "true",
	})
	s.runRequest(req, assertSeeOther(s.kratosServer.URL+"/api/kratos/public/self-service/login/browser"))
}

func (s *ServerTestSuite) TestStartKratosFlow_Login() {
	req := httptest.NewRequest(http.MethodPost, "/user/start/login", nil)
	s.runRequest(req, assertSeeOther(s.kratosServer.URL+"/api/kratos/public/self-service/login/browser"))
}

func (s *ServerTestSuite) TestStartKratosFlow_Login_ReturnToFromForm() {
	form := make(url.Values)
	form.Add("return_to", "https://myserver/page")
	req := httptest.NewRequest(http.MethodPost, "https://myserver/user/start/login",
		strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set("Referer", "https://myserver/ignore")
	s.runRequest(req, assertSeeOther(s.kratosServer.URL+"/api/kratos/public/self-service/login/browser?return_to=https://myserver/page"))
}

func (s *ServerTestSuite) TestStartKratosFlow_Login_ReturnToFromReferer() {
	req := httptest.NewRequest(http.MethodPost, "https://myserver/user/start/login", nil)
	req.Header.Set("Referer", "https://myserver/page")
	s.runRequest(req, assertSeeOther(s.kratosServer.URL+"/api/kratos/public/self-service/login/browser?return_to=https://myserver/page"))
}

func (s *ServerTestSuite) TestStartKratosFlow_Login_ReturnToNotInDomain() {
	req := httptest.NewRequest(http.MethodPost, "https://myserver/user/start/login", nil)
	req.Header.Set("Referer", "https://anotherserver/page")
	s.runRequest(req, assertSeeOther(s.kratosServer.URL+"/api/kratos/public/self-service/login/browser"))
}

func (s *ServerTestSuite) TestStartKratosFlow_Logout() {
	req := httptest.NewRequest(http.MethodPost, "/user/start/logout", nil)
	req.AddCookie(verifiedCookie)
	s.runRequest(req, assertSeeOther("targetURL"))
}
