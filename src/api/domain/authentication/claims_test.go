package authentication

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateNewClaimsForOwner(t *testing.T){

	claims := CreateNewClaimsForOwner("session")

	assert.NotNil(t, claims)
	assert.EqualValues(t, true, claims.IsOwner)
	assert.EqualValues(t, true, claims.Permissions.CanDelete)
	assert.EqualValues(t, true, claims.Permissions.CanRead)
	assert.EqualValues(t, true, claims.Permissions.CanEdit)
	assert.EqualValues(t, true, claims.Permissions.CanUpload)
	assert.EqualValues(t, true, claims.Permissions.CanShare)
	assert.EqualValues(t, "session", claims.Session)
}

func TestCreateNewClaims(t *testing.T){

	claims := CreateNewClaims("session", false, false, true, true, false, false)

	assert.NotNil(t, claims)
	assert.EqualValues(t, false, claims.IsOwner)
	assert.EqualValues(t, false, claims.Permissions.CanDelete)
	assert.EqualValues(t, false, claims.Permissions.CanRead)
	assert.EqualValues(t, false, claims.Permissions.CanEdit)
	assert.EqualValues(t, true, claims.Permissions.CanUpload)
	assert.EqualValues(t, true, claims.Permissions.CanShare)
	assert.EqualValues(t, "session", claims.Session)
}

func TestToMapOwnerValid(t *testing.T){

	claims := CreateNewClaimsForOwner("session")

	assert.NotNil(t, claims)
	assert.EqualValues(t, true, claims.IsOwner)
	assert.EqualValues(t, true, claims.Permissions.CanDelete)
	assert.EqualValues(t, true, claims.Permissions.CanRead)
	assert.EqualValues(t, true, claims.Permissions.CanEdit)
	assert.EqualValues(t, true, claims.Permissions.CanUpload)
	assert.EqualValues(t, true, claims.Permissions.CanShare)
	assert.EqualValues(t, "session", claims.Session)

	toMapValues := claims.ToMap()
	assert.NotNil(t, toMapValues)
	assert.EqualValues(t, true, toMapValues["is_owner"])
	assert.NotNil(t, toMapValues["permissions"])

	permissions := toMapValues["permissions"].(map[string]interface{})

	assert.EqualValues(t, true, permissions["can_read"])
	assert.EqualValues(t, true, permissions["can_edit"])
	assert.EqualValues(t, true, permissions["can_delete"])
	assert.EqualValues(t, true, permissions["can_upload"])
	assert.EqualValues(t, true, permissions["can_share"])
	assert.EqualValues(t, "session", toMapValues["session"])
}

//func TestToMapNotOwnerValid(t *testing.T){
//
//	claims := CreateNewClaims("session", false, true, false, false, false, false)
//
//	assert.NotNil(t, claims)
//	assert.EqualValues(t, false, claims.IsOwner)
//	assert.EqualValues(t, false, claims.CanDelete)
//	assert.EqualValues(t, true, claims.CanRead)
//	assert.EqualValues(t, false, claims.CanEdit)
//	assert.EqualValues(t, false, claims.CanUpload)
//	assert.EqualValues(t, false, claims.CanShare)
//	assert.EqualValues(t, "session", claims.Session)
//
//	toMapValues := claims.ToMap()
//	assert.NotNil(t, toMapValues)
//	assert.Nil(t, toMapValues["is_owner"])
//	assert.Nil(t, toMapValues["can_delete"])
//	assert.NotNil(t, toMapValues["can_read"])
//	assert.EqualValues(t, true, toMapValues["can_read"])
//	assert.Nil(t, toMapValues["can_edit"])
//	assert.Nil(t, toMapValues["can_upload"])
//	assert.Nil(t, toMapValues["can_share"])
//	assert.EqualValues(t, "session", toMapValues["session"])
//}