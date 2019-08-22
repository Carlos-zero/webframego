package main

import (
	"fmt"
	"webframego/web"
)

func main() {
	stall := web.Default()
	stall.Post("/testPost", func(ctx *web.Context) {
		ctx.ResString("测试！")
		ctx.Test()
		test := ctx.Params["test"]
		fmt.Println(test)
	})
	stall.Get("/testGet", func(ctx *web.Context) {
		ctx.ResString("测试！")
		ctx.Test()
		test := ctx.Params["test"]
		fmt.Println(test)
	})

	stall.Post("/testPost11", func(ctx *web.Context) {
		ctx.ResString("测试！")
		ctx.Test()
		test := ctx.Params["test"]
		fmt.Println(test)
	})

	stall.Get("/user/:username/post/:post_id", func(ctx *web.Context) {
		ctx.ResString("测试！")
		ctx.Test()
		username := ctx.Params["username"]
		fmt.Println("username：", username)
		post_id := ctx.Params["post_id"]
		fmt.Println("post_id:", post_id)
	})

	stall.Get("/:user/:username/post/:post_id", func(ctx *web.Context) {
		ctx.ResString("测试！")
		ctx.Test()
		user := ctx.Params["user"]
		fmt.Println("user", user)
		username := ctx.Params["username"]
		fmt.Println("username:", username)
		post_id := ctx.Params["post_id"]
		fmt.Println("post_id:", post_id)
	})

	stall.WebSocket("/ws", web.WebSocketConfig{
		OnOpen: func(ctx *web.Context) {
			fmt.Println("OnOpen!")
		},
		OnClose: func(ctx *web.Context) {
			fmt.Println("OnClose!")
		},
		OnMessage: func(ctx *web.Context, s string) {
			fmt.Println(ctx.Conn.RemoteAddr(), "send message:", s)
			//ctx.ResString(s)
		},
		OnError: nil,
	})

	stall.Run(":8089")
}
