package sandbox

import (
	"github.com/pagoda-tech/bastion/models"
	"testing"
)

func TestManagerGetOrCreate(t *testing.T) {
	opt := ManagerOptions{
		ImageName: "pagodatech/sandbox",
		DataDir:   "/tmp/test-bastion",
	}
	m, err := NewManager(opt)
	if err != nil {
		t.Error(err)
	}
	u := models.User{
		Login:      "test",
		PublicKey:  "test",
		PrivateKey: "test",
	}
	u.ID = 1
	s, err := m.GetOrCreate(u)
	if err != nil {
		t.Error(err)
	}
	t.Log(s)
}
