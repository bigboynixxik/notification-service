package models

type NotificationEvent struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
}

func (n *NotificationEvent) IsValid() bool {
	return n.UserID != "" && n.Message != ""
}
