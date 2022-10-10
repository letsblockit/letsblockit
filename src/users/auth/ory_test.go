package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/pages"
	"github.com/letsblockit/letsblockit/src/server/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	verifiedCookie = &http.Cookie{
		Name:  "ory_session_verified",
		Value: "true",
	}
	whoAmiPattern = `{
	  "id": "af9b460f-4ca0-453d-8bc7-cf68f30d4174",
	  "active": %t,
	  "identity": {
		"id": "%s",
		"verifiable_addresses": [
		  {
			"verified": %t
		  }
		]
	  }
	}`
)

func TestNilOryUser(t *testing.T) {
	var user *oryUser
	assert.False(t, user.IsActive())
	assert.EqualValues(t, "", user.Id())
}

func TestVerifiedOryUser(t *testing.T) {
	uid := uuid.NewString()
	payload := fmt.Sprintf(whoAmiPattern, true, uid, true)
	user := new(oryUser)
	assert.NoError(t, json.Unmarshal([]byte(payload), user))
	assert.True(t, user.IsActive())
	assert.Equal(t, uid, user.Id())
}

func TestInactiveOrySession(t *testing.T) {
	payload := fmt.Sprintf(whoAmiPattern, false, uuid.NewString(), true)
	user := new(oryUser)
	assert.NoError(t, json.Unmarshal([]byte(payload), user))
	assert.False(t, user.IsActive())
}

type OryBackendSuite struct {
	suite.Suite
	echo         *echo.Echo
	expectP      *mocks.MockPageRendererMockRecorder
	kratosServer *httptest.Server
	user         string
}

func (s *OryBackendSuite) SetupTest() {
	c := gomock.NewController(s.T())
	pm := mocks.NewMockPageRenderer(c)
	s.expectP = pm.EXPECT()

	s.kratosServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		fmt.Println(r.URL.Path)
		switch r.URL.Path {
		case "/self-service/logout/browser":
			_, err = fmt.Fprint(w, `{"logout_url":"targetURL"}`)
		case "/sessions/whoami":
			cookie, _ := r.Cookie("ory_session_verified")
			_, err = fmt.Fprintf(w, whoAmiPattern, true, s.user, cookie.Value == "true")
		case "/self-service/login/flows":
			switch r.URL.RawQuery {
			case "id=123456":
				_, err = fmt.Fprint(w, `{"ui":{"a": "1", "b": "2"},"return_to":"https://target"}`)
			case "id=666":
				_, err = fmt.Fprint(w, `{"invalid": true}`)
			}
		default:
			http.NotFound(w, r)
		}
		s.NoError(err)
	}))

	s.user = uuid.New().String()

	ory := NewOryBackend(s.kratosServer.URL, pm, &statsd.NoOpClient{})
	s.echo = echo.New()
	s.echo.Use(ory.BuildMiddleware())
	s.echo.Any("/", func(c echo.Context) error {
		return c.String(200, GetUserId(c))
	})
	ory.RegisterRoutes(s.echo)
}

func (s *OryBackendSuite) TearDownTest() {
	s.kratosServer.Close()
}

func assertOk(t *testing.T, rec *httptest.ResponseRecorder) {
	assert.Equal(t, http.StatusOK, rec.Code, rec.Body)
}

func (s *OryBackendSuite) runRequest(req *http.Request, checks func(*testing.T, *httptest.ResponseRecorder)) {
	s.T().Helper()
	rec := httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	checks(s.T(), rec)
}

func (s *OryBackendSuite) TestGet_Anonymous() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 200, rec.Code)
		assert.Empty(t, rec.Body)
	})
}

func (s *OryBackendSuite) TestGet_Logged_HasCookie() { // Don't set the has_account cookie if found
	s.echo.GET("/check", func(c echo.Context) error {
		assert.True(s.T(), HasAccount(c))
		assert.True(s.T(), HasAuth(c))
		assert.Equal(s.T(), s.user, GetUserId(c))
		return nil
	})
	req := httptest.NewRequest(http.MethodGet, "/check", nil)
	req.AddCookie(verifiedCookie)
	req.AddCookie(&http.Cookie{
		Name:  "has_account",
		Value: "true",
	})
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 200, rec.Code)
		assert.Empty(t, rec.Header().Get("Set-Cookie"))
	})
}

func (s *OryBackendSuite) TestGet_Logged_SetCookie() { // Set the has_account cookie on first login
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(verifiedCookie)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, s.user, rec.Body.String())
		cookies := rec.Result().Cookies()
		require.Len(t, cookies, 1)
		assert.Equal(t, "has_account", cookies[0].Name)
		assert.Equal(t, "true", cookies[0].Value)
		assert.Equal(t, "/", cookies[0].Path)
		assert.WithinDuration(t, time.Now().AddDate(10, 0, 0), cookies[0].Expires, time.Second)
		assert.True(t, cookies[0].HttpOnly)
		assert.Equal(t, http.SameSiteStrictMode, cookies[0].SameSite)
	})
}

func (s *OryBackendSuite) TestGet_BadUUID() { // Request goes through unauthenticated
	s.user = "invalid-uuid"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(verifiedCookie)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 200, rec.Code)
		assert.Empty(t, rec.Body)
	})
}

func (s *OryBackendSuite) TestGet_CachedUser() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(verifiedCookie)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, s.user, rec.Body.String())
	})

	// Shutdown Kratos, we can still authenticate from cache
	s.kratosServer.Close()
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, s.user, rec.Body.String())
	})
}

func (s *OryBackendSuite) TestGet_KratosDown() { // Request goes through unauthenticated
	s.kratosServer.Close()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(verifiedCookie)
	s.runRequest(req, func(t *testing.T, rec *httptest.ResponseRecorder) {
		assert.Equal(t, 200, rec.Code)
		assert.Empty(t, rec.Body)
	})
}

func (s *OryBackendSuite) TestRenderKratosForm_OK() {
	req := httptest.NewRequest(http.MethodGet, "/user/forms/login?flow=123456", nil)
	s.expectP.BuildPageContext(gomock.Any(), gomock.Any()).
		Return(&pages.Context{OfficialInstance: true})
	s.expectP.Render(gomock.Any(), "kratos-form", gomock.Any()).
		DoAndReturn(func(_ echo.Context, _ string, c *pages.Context) error {
			assert.True(s.T(), c.OfficialInstance)
			assert.EqualValues(s.T(), pages.ContextData{
				"type": "login",
				"ui": map[string]interface{}{
					"a": "1",
					"b": "2",
				},
				"return_to": "https://target",
				"settings":  supportedForms["login"],
			}, c.Data)
			return nil
		})
	s.runRequest(req, assertOk)
}

func (s *OryBackendSuite) TestRenderKratosForm_ErrFormType() {
	req := httptest.NewRequest(http.MethodGet, "/user/forms/unknown?flow=123456", nil)
	s.expectP.Redirect(gomock.Any(), 302, s.kratosServer.URL+"/ui/unknown?flow=123456")
	s.runRequest(req, assertOk)
}

func (s *OryBackendSuite) TestRenderKratosForm_ErrBadFlow() {
	req := httptest.NewRequest(http.MethodGet, "/user/forms/login?flow=666", nil)
	s.expectP.Redirect(gomock.Any(), 302, s.kratosServer.URL+"/ui/login?flow=666")
	s.runRequest(req, assertOk)
}

func (s *OryBackendSuite) TestStartKratosFlow_Settings() {
	req := httptest.NewRequest(http.MethodPost, "/user/action/settings", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/settings/browser")
	s.runRequest(req, assertOk)
}

func (s *OryBackendSuite) TestStartKratosFlow_LoginOrRegistration_Register() {
	req := httptest.NewRequest(http.MethodPost, "/user/action/loginOrRegistration", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/registration/browser")
	s.runRequest(req, assertOk)
}

func (s *OryBackendSuite) TestStartKratosFlow_LoginOrRegistration_Login() {
	req := httptest.NewRequest(http.MethodPost, "/user/action/loginOrRegistration", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(&http.Cookie{
		Name:  "has_account",
		Value: "true",
	})
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/login/browser")
	s.runRequest(req, assertOk)
}

func (s *OryBackendSuite) TestStartKratosFlow_Login() {
	req := httptest.NewRequest(http.MethodPost, "/user/action/login", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/login/browser")
	s.runRequest(req, assertOk)
}

func (s *OryBackendSuite) TestStartKratosFlow_Login_ReturnToFromForm() {
	form := make(url.Values)
	form.Add("return_to", "https://myserver/page")
	req := httptest.NewRequest(http.MethodPost, "https://myserver/user/action/login",
		strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set("Referer", "https://myserver/ignore")
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/login/browser?return_to=https://myserver/page")
	s.runRequest(req, assertOk)
}

func (s *OryBackendSuite) TestStartKratosFlow_Login_ReturnToFromReferer() {
	req := httptest.NewRequest(http.MethodPost, "https://myserver/user/action/login", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set("Referer", "https://myserver/page")
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/login/browser?return_to=https://myserver/page")
	s.runRequest(req, assertOk)
}

func (s *OryBackendSuite) TestStartKratosFlow_Login_ReturnToNotInDomain() {
	req := httptest.NewRequest(http.MethodPost, "https://myserver/user/action/login", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set("Referer", "https://anotherserver/page")
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/login/browser")
	s.runRequest(req, assertOk)
}

func (s *OryBackendSuite) TestStartKratosFlow_Logout() {
	req := httptest.NewRequest(http.MethodPost, "/user/action/logout", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.AddCookie(verifiedCookie)
	s.expectP.Redirect(gomock.Any(), 303, "targetURL")
	s.runRequest(req, assertOk)
}

func (s *OryBackendSuite) TestStartLoginFlow_Register() {
	req := httptest.NewRequest(http.MethodGet, "/user/action/loginOrRegistration", nil)
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/registration/browser")
	s.runRequest(req, assertOk)
}

func (s *OryBackendSuite) TestStartLoginFlow_Login() {
	req := httptest.NewRequest(http.MethodGet, "/user/action/loginOrRegistration", nil)
	req.AddCookie(&http.Cookie{
		Name:  "has_account",
		Value: "true",
	})
	s.expectP.Redirect(gomock.Any(), 303, s.kratosServer.URL+"/self-service/login/browser")
	s.runRequest(req, assertOk)
}

func TestOryBackendSuite(t *testing.T) {
	suite.Run(t, new(OryBackendSuite))
}
