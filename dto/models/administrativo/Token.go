package administrativo

type Token struct {
	Token    string "redis:token"
	Refresh  string "redis:refresh_id"
	AccessID uint64 "redis:access_id"
}
