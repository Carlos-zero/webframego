package web

import (
	"bufio"
	"net"
	"net/http"
	"sync"
)

//这里是结构体或者interface

//相当于gin的engine		作为最外层的
type Staff struct {
	//路由
	Router Router

	//临时对象池  用于处理context
	pool sync.Pool

	//websocket功能能
	websocketMap map[string]WebSocketConfig
}

/**/
type WebSocketConfig struct {
	OnOpen    func(*Context)
	OnClose   func(*Context)
	OnMessage func(*Context, string)
	OnError   func(*Context)
}


//动态路由匹配机制
type Nodes map[int]string

//节点组
type NodesMap map[int]Nodes

//函数操作
type Handle func(*Context)

//函数操作的集合
type Routers map[string]Handle

//参数
type Params map[string]string

//参数map
type ParamsMap map[string]Params

//参数名map
//type ParamsName map[int]string
type ParamsName []string

//参数名mapmap
//type ParamsNameMap map[string]ParamsName
type ParamsNameMap map[string]ParamsName

//路由
type Router struct {
	//该路由的方法   post   or   get
	Method string
	//开发者给的URI
	Path    string
	Handle  Handle
	Routers Routers
	staff   *Staff
	Params  Params

	//尝试放到这里面   然后再调出来
	//Nodes  Nodes
	//计数器   第几个nodes
	index         int
	NodesMap      NodesMap
	ParamsMap     ParamsMap
	ParamsNameMap ParamsNameMap
}

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Handle         Handle
	Params         Params
	Uri            string
	staff          *Staff

	//websocket的内容
	Conn net.Conn
	buf  *bufio.ReadWriter
}

//各种各样的功能      =可以在这里写一些小功能   使用ctx实现
type Various interface {
	//返回string
	ResString(str string)
}
