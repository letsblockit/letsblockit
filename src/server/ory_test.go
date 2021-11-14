package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/xvello/letsblockit/src/pages"
)

func TestNilOryUser(t *testing.T) {
	var user *oryUser
	assert.False(t, user.IsActive())
	assert.False(t, user.IsVerified())
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
	assert.True(t, user.IsVerified())
	assert.Equal(t, "9a3f8aeb-729a-44cf-bede-f885175344ef", user.Id().String())
}

func TestUnverifiedOryUser(t *testing.T) {
	payload := `{
	  "id": "d631b403-eb29-4a5b-8829-125da6ebdf75",
	  "active": true,
	  "identity": {
		"id": "9a3f8aeb-729a-44cf-bede-f885175344ef"
	  }
	}`
	user := new(oryUser)
	assert.NoError(t, json.Unmarshal([]byte(payload), user))
	assert.True(t, user.IsActive())
	assert.False(t, user.IsVerified())
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
		"settings": supportedForms["login"],
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
