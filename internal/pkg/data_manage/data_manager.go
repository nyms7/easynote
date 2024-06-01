package data

import (
	"context"
	"easynote/logs"
	"easynote/utils"
	"fmt"
	"sync"
	"sync/atomic"
)

var GlobalManater *NoteManager

type NoteManager struct {
	Seed    string
	LastID  uint64
	CodeMap sync.Map // key: code, value: id
	NoteMap sync.Map // key: id, value: note
}

func (m *NoteManager) Apply(ctx context.Context) string {
	for times := 0; times < 10000; times++ {
		id := atomic.AddUint64(&m.LastID, 1)
		code := utils.HashAndTrim(fmt.Sprintf("%d%s", id, m.Seed))
		if _, ok := m.CodeMap.Load(code); !ok {
			m.CodeMap.Store(code, id)
			return code
		}
	}
	logs.CtxInfo(ctx, "[NoteManager.Apply] apply new code failed, last id: %d, seed: %s", m.LastID, m.Seed)
	return ""
}

type Note struct {
	ID       string
	Data     []byte
	NoteMeta NoteMeta
}

type NoteMeta struct {
	Password  string
	CreatedAt int64
	UpdatedAt int64
}
