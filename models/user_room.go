package models

type UserRoom struct {
	UserIdentity     string `json:"user_identity"`
	RoomIdentity     string `json:"room_identity"`
	Message_identity string `json:"message___identity"`
	CreateAt         int64  `json:"create_at"`
	UpdateAt         int64  `json:"update_at"`
}

func (UserRoom) CollectionName() string {
	return "user_room"
}
