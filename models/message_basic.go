package models

type MessageBasic struct {
	Identity     string `json:"identity"`
	UserIdentity string `json:"user_identity"`
	RootIdentity string `json:"root_identity"`
	Data         string `json:"data"`
	CreateAt     int64  `json:"create_at"`
	UpdateAt     int64  `json:"update_at"`
}

func (MessageBasic) CollectionName() string {
	return "message_basic"
}
