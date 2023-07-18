package src

// Envs

type Envs struct {
	DbUrl     string
	Token     string
	Encyptkey string
	LogChat   int64
	Port      string
}

type IndexResponse struct {
	Message string `json:"message"`
	Version string `json:"version"`
}

type BanRequest struct {
	Token        string `json:"token"`
	UserId       string `json:"user_id"`
	Reason       string `json:"reason"`
	From         string `json:"from"`
	BanClass     string `json:"ban_class"`
	EvidenceLink string `json:"evidence_link,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type UnbanRequest struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
	Reason string `json:"reason"`
	From   string `json:"from"`
}

type GeneralResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type TokenRequest struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
	Role   string `json:"role"`
}

type TokenResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type BannedInfo struct {
	BanRequest
	TimeStamp string `bson:"date"`
	BanId     string `bson:"ban_id"`
}

type TokensInfo struct {
	UserId     string `bson:"_id"`
	Token      string `bson:"token"`
	Time       string `bson:"time"`
	Role       string `bson:"role"`
	AssignedBy string `bson:"assigned_by"`
}
