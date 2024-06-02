package handler

import (
	"context"
	"easynote/conf"
	"easynote/data_manager"
	"easynote/logs"
	"easynote/utils"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"
)

//go:embed static/template/index.html
var templateFS embed.FS

type NotePayload struct {
	Token   string `json:"token"`
	Content string `json:"content"`
}

func NoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(context.Background(), logs.LogID("logID"), fmt.Sprintf("%d", time.Now().UnixNano()))
	path := r.URL.Path

	if path == "/" {
		if code := data_manager.Apply(ctx); code == "" {
			utils.Response(w, r, http.StatusInternalServerError, "apply code failed", nil)
		} else {
			w.Header().Set("Cache-Control", "no-store")
			http.Redirect(w, r, "/"+code, http.StatusFound)
		}
		return
	}

	code := path[1:]

	if r.Method == http.MethodGet {
		// redirect long code
		if len(code) > conf.MaxCodeSize() {
			if code := data_manager.Apply(ctx); code == "" {
				utils.Response(w, r, http.StatusInternalServerError, "apply code failed", nil)
			} else {
				w.Header().Set("Cache-Control", "no-store")
				http.Redirect(w, r, "/"+code, http.StatusFound)
			}
			return
		}
		note := data_manager.Load(ctx, code)
		if note != nil {
			if note.NoteMeta.Token != "" {
				token := r.URL.Query().Get("token")
				if token != note.NoteMeta.Token {
					utils.Response(w, r, http.StatusForbidden, "token auth failed", nil)
					return
				}
			}
			renderNotePage(ctx, w, r, note)
			return
		}
		renderNotePage(ctx, w, r, note)
		return
	}

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.Response(w, r, http.StatusBadRequest, "invalid params", nil)
			return
		}
		defer r.Body.Close()

		var payload NotePayload
		if err := json.Unmarshal(body, &payload); err != nil {
			utils.Response(w, r, http.StatusBadRequest, "invalid too long", nil)
			return
		}

		logs.CtxInfo(ctx, "[NoteHandler] post payload: %+v", payload)

		// check content size
		if len(payload.Content) > conf.MaxContentSize() {
			utils.Response(w, r, http.StatusBadRequest, "payload too long", nil)
			return
		}

		// check code size
		if len(code) > conf.MaxCodeSize() {
			utils.Response(w, r, http.StatusBadRequest, "code too long", nil)
			return
		}

		// check token size
		if len(payload.Token) > conf.MaxTokenSize() {
			utils.Response(w, r, http.StatusBadRequest, "token too long", nil)
			return
		}

		if err := data_manager.Upsert(ctx, code, payload.Token, payload.Content); err != nil {
			utils.Response(w, r, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		utils.RespondSuccess(w, r, nil)
		return
	}

	utils.Response(w, r, http.StatusMethodNotAllowed, "method not allowed", nil)
}

func renderNotePage(ctx context.Context, w http.ResponseWriter, r *http.Request, note *data_manager.Note) {
	tmplData, err := templateFS.ReadFile("static/template/index.html")
	if err != nil {
		logs.CtxError(ctx, "[renderPage] read template file err: %v", err)
		utils.Response(w, r, http.StatusInternalServerError, "read template file error", nil)
		return
	}

	tmpl, err := template.New("template.html").Parse(string(tmplData))
	if err != nil {
		logs.CtxError(ctx, "[renderPage] parse template file err: %v", err)
		utils.Response(w, r, http.StatusInternalServerError, "parse template file error", nil)
		return
	}

	if note == nil {
		note = &data_manager.Note{}
	}

	err = tmpl.Execute(w, note)
	if err != nil {
		logs.CtxError(ctx, "[renderPage] render template file err: %v", err)
		utils.Response(w, r, http.StatusInternalServerError, "render template file error", nil)
		return
	}
}
