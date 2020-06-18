package authentication

import (
	"encoding/json"
	"github.com/lelinu/api_utils/utils/random_utils"
)

type Claims struct {
	Id          string       `json:"id"`
	Session     string       `json:"session"`
	IsOwner     bool         `json:"is_owner,omitempty"`
	Permissions *Permissions `json:"permissions,omitempty"`
}

type Permissions struct {
	CanRead   bool `json:"can_read,omitempty"`
	CanUpload bool `json:"can_upload,omitempty"`
	CanShare  bool `json:"can_share,omitempty"`
	CanEdit   bool `json:"can_edit,omitempty"`
	CanDelete bool `json:"can_delete,omitempty"`
}

func CreateNewClaimsForOwner(session string) *Claims {

	id, _ := random_utils.NewUUID()

	return &Claims{
		Id:      id,
		Session: session,
		IsOwner: true,
		Permissions: &Permissions{
			CanRead:   true,
			CanUpload: true,
			CanShare:  true,
			CanEdit:   true,
			CanDelete: true,
		},
	}
}

func CreateNewClaims(session string, isOwner bool, canRead bool, canUpload bool,
	canShare bool, canEdit bool, canDelete bool) *Claims {

	id, _ := random_utils.NewUUID()

	return &Claims{
		Id:      id,
		Session: session,
		IsOwner: isOwner,
		Permissions: &Permissions{
			CanRead:   canRead,
			CanUpload: canUpload,
			CanShare:  canShare,
			CanEdit:   canEdit,
			CanDelete: canDelete,
		},
	}
}

//ToMap will convert the c struct to map
func (c *Claims) ToMap() map[string]interface{} {

	var result map[string]interface{}

	jsonData, _ := json.Marshal(c)
	json.Unmarshal(jsonData, &result)

	return result
}
