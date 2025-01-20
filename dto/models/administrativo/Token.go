package administrativo

import "github.com/google/uuid"

type Token struct {
	Token    string    "redis:token"
	Refresh  string    "redis:refresh_id"
	AccessID uuid.UUID "redis:access_id"
}
