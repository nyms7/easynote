package main

import (
	"context"
	"easynote/conf"
	"easynote/data_manager"
	"easynote/handler"
	"easynote/logs"
	"easynote/utils"
	"flag"
	"fmt"
	"math"
	"net/http"
	"time"
)

func main() {
	var port int
	var token string

	flag.IntVar(&port, "p", 9600, "easynote server listen port")
	flag.StringVar(&token, "t", "", "easynote admin token")
	flag.Parse()

	initService(token)

	serve(port)
}

func initService(token string) {
	logs.Infof("start init service")
	initConf(token)
	initNoteManager()
	initCleaner()
	logs.Infof("start service success")
}

func initConf(token string) {
	if token == "" {
		token, _ = utils.SecureRandString(16)
	}
	conf.GlobalConf = &conf.NoteConf{
		MaxCodeSize:    256,
		MaxContentSize: int(math.Pow(2, 20)), // 1MB
		MaxTokenSize:   32,
		AdminToken:     token,
	}
	logs.Infof("init conf success, admin token: %s", token)
}

func initNoteManager() {
	seed, err := utils.SecureRandString(16)
	if err != nil {
		logs.CtxError(context.Background(), "[init] SecureRandString err: %+v", err)
		panic("generate rand seed err")
	}

	data_manager.GlobalManater = &data_manager.NoteManager{
		Seed: seed,
	}
	logs.Infof("init note data manager success, seed: %s", seed)
}

func initCleaner() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Errorf("Recovered from panic: %v", r)
			}
		}()
		for {
			now := time.Now()
			// 凌晨四点清除所有数据
			next := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}
			duration := next.Sub(now)
			logs.Infof("time until next clean task: %v", duration)

			timer := time.NewTimer(duration)
			<-timer.C

			data_manager.GlobalManater.Reset()
		}
	}()
	logs.Infof("init cleaner success")
}

func serve(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.NoteHandler)
	mux.HandleFunc("/api/stat", handler.StatHandler)

	logs.Infof("start service on http://localhost:%d", port)
	_ = http.ListenAndServe(":"+fmt.Sprintf("%d", port), utils.RateLimit(mux))
}
