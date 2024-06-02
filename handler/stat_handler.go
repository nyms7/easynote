package handler

import (
	"easynote/conf"
	"easynote/data_manager"
	"easynote/utils"
	"net/http"
)

func StatHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token != conf.AdminToken() {
		utils.Response(w, r, http.StatusForbidden, "token auth failed", nil)
		return
	}
	utils.RespondSuccess(w, r, data_manager.GlobalManater)
}
