package models

type RoomBasic struct {
	Identity     string `json:"identity"`
	UserIdentity string `json:"user_identity"`
	Number       string `json:"number"`
	Name         string `json:"name"`
	Info         string `json:"info"`
	CreateAt     int64  `json:"create_at"`
	UpdateAt     int64  `json:"update_at"`
}

func (RoomBasic) CollectionName() string {
	return "room_basic"
}
