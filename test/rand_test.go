package test

import (
	"context"
	"easynote/data_manager"
	"easynote/utils"
	"testing"
)

func TestSecureRandString(t *testing.T) {
	str, err := utils.SecureRandString(16)
	if err != nil {
		t.Fatalf("[TestSecureRandString] err: %+v", err)
	}
	if str == "" {
		t.Fatalf("[TestSecureRandString] failed")
	}
	t.Logf("[TestSecureRandString] result: %v", str)
}

func TestNoteManagerApply(t *testing.T) {
	seed, _ := utils.SecureRandString(16)
	t.Logf("seed: %s", seed)
	if seed == "" {
		t.Fatalf("[TestNote] generate seed failed")
	}
	m := &data_manager.NoteManager{
		Seed: seed,
	}
	for i := 0; i < 1000; i++ {
		t.Logf("applyed code: " + m.Apply(context.Background()))
	}
}
