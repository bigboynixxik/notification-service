package models

import "testing"

func TestNotificationEvent_IsValid(t *testing.T) {
	e := NotificationEvent{UserID: "1", Message: "hi"}
	if !e.IsValid() {
		t.Error("Should be valid")
	}

	e2 := NotificationEvent{}
	if e2.IsValid() {
		t.Error("Should be invalid")
	}
}
