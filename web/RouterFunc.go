package web

import (
	"strings"
)

//将uri和handle组合成一个路由
func (router *Router) ComposRouter(uri string, handle Handle) {
	router.staff.Router.Path = uri
	router.staff.Router.Handle = handle
	router.staff.Router.Routers[uri] = handle
	//fmt.Println(uri,"22222222222,ComposRouter")
}

//解析开发者路由
func (router *Router) JudgeDynamic(uri string) {
	newUri:=SimplyUri(uri)

	uris := strings.Split(uri, "/")
	nodes := make(map[int]string)
	params := make(map[string]string)
	//paramsNameMap:=*router.ParamsNameMap
	a:=0
	for i, _ := range uris {
		if i == 0 {
			nodes[i] = newUri
		} else if strings.Contains(uris[i], ":") {
			nodes[i] = strings.Split(uris[i], ":")[1]
			params[nodes[i]]=""
			//fmt.Println(len(router.ParamsNameMap[newUri]),"----------")
			router.ParamsNameMap[newUri]=append(router.ParamsNameMap[newUri],nodes[i])
			//fmt.Println(router.ParamsNameMap[newUri][a],"99999999")

			a++
		} else {
			nodes[i] = uris[i]
		}
	}
	//最多存10个
	router.NodesMap[router.index] = nodes
	router.ParamsMap[newUri]=params
	router.index++
}

