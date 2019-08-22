package web

import (
	"fmt"
	"net/http"
)

func (staff *Staff) WebSocket(uri string, config WebSocketConfig) {
	staff.Router.Method = "GET"
	staff.Router.ComposRouter(uri, Enter)
	staff.websocketMap[uri] = config
}

//staff的get请求
func (staff *Staff) Get(uri string, handle Handle) {
	//uri="GET:"+uri
	staff.Router.Method = "GET"
	//简化路由
	staff.Router.JudgeDynamic(uri)
	uri = SimplyUri(uri)
	//这里是开发者开发时，进行的组装
	staff.Router.ComposRouter(uri, handle)
}

func (staff *Staff) Post(uri string, handle Handle) {
	//uri="POST:"+uri
	staff.Router.Method = "POST"
	staff.Router.ComposRouter(uri, handle)
}

//开启端口服务
func (staff *Staff) Run(port string) {
	err := http.ListenAndServe(port, staff)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (staff *Staff) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//从对象池中获得一个新文本     不知道为啥可以高效
	ctx := staff.pool.Get().(*Context)
	//fmt.Println(r.Method,"ServeHTTP")
	ctx.ResponseWriter = w
	ctx.Request = r

	//通过URI从路由中取出方法，并执行
	//通过请求对路由进行解析
	ctx.ParseRouter(w, r)

	//将文本放回对象池
	staff.pool.Put(ctx)
}

//将文本初始化
func (staff *Staff) InitContext() *Context {
	return &Context{
		staff: staff,
		//参数必须初始化，该map为nil对象
		Params: map[string]string{},
	}
}
