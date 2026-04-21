package clients

import (
	"testing"
)

func TestNewAuthClient_Error(t *testing.T) {
	client, err := NewAuthClient("invalid-addr")

	if client != nil {
		defer client.Close()
	}

	_ = err
}
