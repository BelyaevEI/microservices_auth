package access

import "github.com/BelyaevEI/microservices_auth/internal/model"

var (
	accessibleRoles = map[string]map[model.Role]struct{}{
		"/chat_v1.ChatV1/Create": {
			model.RoleAdmin: {},
		},
		"/chat_v1.ChatV1/Delete": {
			model.RoleAdmin: {},
		},
		"/chat_v1.ChatV1/SendMessage": {
			model.RoleUser:  {},
			model.RoleAdmin: {},
		},
	}
)
