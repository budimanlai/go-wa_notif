package main

import (
	"net/url"
	"strconv"
	"sync"
	"time"

	services "github.com/budimanlai/go-cli-service"
	"github.com/eqto/dbm"
)

var (
	baseUrl string
)

func StartService(ctx *services.Service) {
	baseUrl = ctx.Config.GetString("wa.url")
	rate_limit := strconv.Itoa(ctx.Config.GetIntOr("wa.rate_limit", 10))
	rate_time := ctx.Config.GetIntOr("wa.rate_time", 2)

	var duration time.Duration = time.Duration(rate_time) * time.Second
	var wg sync.WaitGroup

	for {
		result, e := ctx.Db.Select("SELECT * FROM wa_messages WHERE status = 'pending' LIMIT " + rate_limit)
		if e != nil {
			ctx.Log(e.Error())
		}

		if len(result) > 0 {
			for _, item := range result {
				wg.Add(1)
				go doSend(&wg, ctx, item)
			}
			wg.Wait()
		}
		ctx.Log("Sleep...")
		time.Sleep(duration)

		if ctx.IsStopped {
			ctx.Log("Exit from loop StartService")
			break
		}
	}
}

func StopService(ctx *services.Service) {
	ctx.Log("Stop Service")
	ctx.IsStopped = true
}

func doSend(wg *sync.WaitGroup, ctx *services.Service, item dbm.Resultset) {
	defer wg.Done()

	ctx.Log("Process ID:", item.Int("id"), ", message:", item.String("content"))

	v := url.Values{}
	v.Add("PhoneNumber", item.String("to_phone"))
	v.Add("Text", item.String("content"))

	resp, e := Get(baseUrl, v.Encode())
	if e != nil {
		ctx.Log("Error:", e)
		_, e1 := ctx.Db.Exec("UPDATE wa_messages SET status = 'error', sended_at = now(), response_log = ? WHERE id = ?",
			e.Error(), item.Int("id"))
		if e1 != nil {
			ctx.Log(e1.Error())
		}
		return
	}

	ctx.Log(string(resp.Body()))

	_, e1 := ctx.Db.Exec("UPDATE wa_messages SET status = 'done', sended_at = now(), response_log = ? WHERE id = ?",
		resp.Body(), item.Int("id"))
	if e1 != nil {
		ctx.Log(e1.Error())
	}
}
