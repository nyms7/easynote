package data_manager

import (
	"context"
	"easynote/logs"
	"easynote/utils"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var GlobalManater *NoteManager

type NoteManager struct {
	Seed    string
	LastID  uint64
	CodeMap sync.Map // key: code, value: id
	NoteMap sync.Map // key: id, value: note
}

type Note struct {
	ID       uint64    `json:"-"`
	Content  string    `json:"content"`
	NoteMeta *NoteMeta `json:"note_meta"`
}

type NoteMeta struct {
	Token     string `json:"-"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (m *NoteManager) MarshalJSON() ([]byte, error) {
	codeMap := make(map[string]uint64)
	m.CodeMap.Range(func(k, v interface{}) bool {
		key, _ := k.(string)
		value, _ := v.(uint64)
		codeMap[key] = value
		return true
	})
	noteMap := make(map[uint64]*Note)
	m.NoteMap.Range(func(k, v interface{}) bool {
		key, _ := k.(uint64)
		value, _ := v.(*Note)
		noteMap[key] = value
		return true
	})
	return json.Marshal(&struct {
		Seed    string            `json:"-"`
		LastID  uint64            `json:"last_id"`
		CodeMap map[string]uint64 `json:"code_map"`
		NoteMap map[uint64]*Note  `json:"note_map"`
	}{
		Seed:    m.Seed,
		LastID:  m.LastID,
		CodeMap: codeMap,
		NoteMap: noteMap,
	})
}

func Apply(ctx context.Context) string {
	return GlobalManater.Apply(ctx)
}

func Load(ctx context.Context, code string) *Note {
	return GlobalManater.Load(ctx, code)
}

func Upsert(ctx context.Context, code, token, content string) error {
	return GlobalManater.Upsert(ctx, code, token, content)
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
	logs.CtxWarn(ctx, "[NoteManager.Apply] apply new code failed, last id: %d, seed: %s", m.LastID, m.Seed)
	return ""
}

func (m *NoteManager) Load(ctx context.Context, code string) *Note {
	if val, ok := m.CodeMap.Load(code); ok {
		id := val.(uint64)
		if noteVal, ok := m.NoteMap.Load(id); ok {
			return noteVal.(*Note)
		}
	}
	return nil
}

func (m *NoteManager) Upsert(ctx context.Context, code, token, content string) error {
	origin := m.Load(ctx, code)
	if origin == nil {
		id := atomic.AddUint64(&m.LastID, 1)
		m.CodeMap.Store(code, id)
		m.NoteMap.Store(id, &Note{
			ID:      id,
			Content: content,
			NoteMeta: &NoteMeta{
				Token:     token,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			},
		})
		return nil
	}
	// first time to set token
	if origin.NoteMeta.Token == "" {
		origin.NoteMeta.Token = token
	}
	if origin.NoteMeta.Token != token {
		return errors.New("token auth failed")
	}
	origin.Content = content
	origin.NoteMeta.UpdatedAt = time.Now().Unix()
	return nil
}

func (m *NoteManager) Reset() {
	m.LastID = 0
	clearSyncMap(&m.NoteMap)
	clearSyncMap(&m.CodeMap)
}

func clearSyncMap(m *sync.Map) {
	m.Range(func(key, value interface{}) bool {
		m.Delete(key)
		return true
	})
}
