package server

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNilOryUser(t *testing.T) {
	var user *oryUser
	assert.False(t, user.IsActive())
	assert.False(t, user.IsVerified())
	assert.Empty(t, user.Id())
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
	assert.Equal(t, "9a3f8aeb-729a-44cf-bede-f885175344ef", user.Id())
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
	assert.Equal(t, "9a3f8aeb-729a-44cf-bede-f885175344ef", user.Id())
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
