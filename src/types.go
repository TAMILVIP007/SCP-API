package src

type Envs struct {
	DbUrl     string
	Token     string
	Encyptkey string
	Port      string
}

type BannedInfo struct {
	UserId   string `bson:"_id"`
	Reason   string `bson:"reason"`
	BannedBy string `bson:"banned_by"`
	Date     string `bson:"date"`
}

type TokensInfo struct {
	UserId string   `bson:"_id"`
	Token  string   `bson:"token"`
	Time   string   `bson:"time"`
	Rights []string `bson:"rights"`
}

type IndexResponse struct {
	Message string `json:"message"`
	Version string `json:"version"`
}

type BanRequest struct {
	Token    string `json:"token"`
	ID       string `json:"id"`
	Reason   string `json:"reason"`
	BannedBy string `json:"bannedBy"`
}

type BanResponse struct {
	Message string `json:"message"`
}

type TokenRequest struct {
	ID     string `json:"id"`
	Rights string `json:"rights"`
}

type TokenResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
