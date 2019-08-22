package web

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

//ctx可调用的方法

//向浏览器返回string字符串   通过ctx调用
func (ctx *Context) ResString(str string) {
	_, err := ctx.ResponseWriter.Write([]byte(str))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (ctx *Context) Test() {
	_, err := ctx.ResponseWriter.Write([]byte("TEST,111111111111111111"))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (ctx *Context) ParseRouter(w http.ResponseWriter, r *http.Request) {
	//上面这种method    在开多个api的时候   会出错   method没办法找对
	//method := ctx.staff.Router.Method
	method := r.Method
	//fmt.Println(method, "ParseRouter")
	if method == "GET" {
		//4，验证是不是websocket
		//fmt.Println(method, "ParseRouter进入GET")
		upgrade := r.Header.Get("Upgrade")
		if upgrade == "websocket" {
			//fmt.Println(ctx.Uri, r.RequestURI, "ParseRouter里面的URI测试")
			handle := ctx.staff.Router.Routers[r.RequestURI]

			handle(ctx)
			return
		} else {
			//对URL进行解析   并将解析内容装入ctx
			//fmt.Println(ctx.Uri, "ParseRouter3333333333333333", r.RequestURI)
			ctx.ParseGet(r.RequestURI)
			//通过ctx中存放的uri来获取方法

			handle := ctx.staff.Router.Routers[ctx.Uri]

			//fmt.Println(ctx.Uri, "ParseRouter1111111111111111")
			if handle == nil {
				ctx.ResponseWriter.Write(QuickStringToBytes("404 not found"))
				return
			}
			handle(ctx)
		}
	} else if method == "POST" {
		//对请求进行解析   将参数装入ctx
		ctx.ParsePost(r)
		ctx.Uri = r.RequestURI
		handle := ctx.staff.Router.Routers[ctx.Uri]
		if handle==nil {
			ctx.ResponseWriter.Write(QuickStringToBytes("404 NOT FOUND!"))
			return
		}
		handle(ctx)
	}
}

//将带参数的get请求的url进行解析
func (ctx *Context) ParseGet(uri string) {
	//如果是带参数的get请求
	if strings.Contains(uri, "?") {
		//分离uri和参数部分
		uriArr := strings.Split(uri, "?")
		//uri
		ctx.Uri = uriArr[0]
		//参数集
		var params []string
		//有可能get请求没有参数   所以需要进行判断
		if len(uriArr) >= 2 {
			params = strings.Split(uriArr[1], "&")
			for _, str := range params {
				param := strings.Split(str, "=")
				key := param[0]
				value := param[1]
				//将参数装入ctx
				ctx.Params[key] = value
			}
		}
	} else {
		//者部分是将用户的uri分解组装到map中    方便匹配
		uris := strings.Split(uri, "/")
		nodes := make(map[int]string)
		for i, _ := range uris {
			if i == 0 {
				//这里暂时不做优化  等会弄
				//这里是用户输入的uri
				nodes[i] = uri
			} else {
				nodes[i] = uris[i]
				//fmt.Println(nodes[i],"ninainai")
			}
		}
		//for key, value := range nodes {
		//	fmt.Println(key, "=", value, "这是uri组装部分")
		//}
		//上面是组装过程
		//下面是匹配过程
		//当前匹配度最高的路由号码
		mapNum := 0
		//当前最大的权值
		weightMax := 0

		var index []int
		//从路由号0开始匹配
		for a := 0; a <= len(ctx.staff.Router.NodesMap); a++ {
			nodeMap := ctx.staff.Router.NodesMap[a]
			//权值   匹配得上的越多  权值越大
			weight := 0

			//对各各节点进行匹配
			for b := 0; b <= len(nodeMap); b++ {
				if nodes[b] == nodeMap[b]&&nodes[b]!="" {
					//fmt.Println(nodes[b],"==",nodeMap[b],"a=",a,"b=",b,"nodeMap[b]=",len(nodeMap))
					weight++
				}
			}
			//如果匹配度大于等于刚才的匹配度，则重新赋值
			if weight >= weightMax {
				weightMax = weight
				mapNum = a
			}
		}

		node := ctx.staff.Router.NodesMap[mapNum]
		for c := 1; c < len(node); c++ {
			if nodes[c] != node[c] {
				//fmt.Println(nodes[c], c, "nodes[c]")
				//fmt.Println(node[c], c, "node[c]")
				index = append(index, c)
			}
		}

		//匹配上的uri
		newUri := ctx.staff.Router.NodesMap[mapNum][0]
		//fmt.Println( ctx.staff.Router.NodesMap[1][0],"。。。。。")
		params := ctx.staff.Router.ParamsMap[newUri]
		//fmt.Println(newUri,"。。。。。")
		//paramsNameMap:=*ctx.staff.Router.ParamsNameMap


		for d, i := range index {
			//fmt.Println(len(params), "切片长度", d)
			//fmt.Println(i, "?????")
			if d>= len(ctx.staff.Router.ParamsNameMap[newUri]){
				return
			}
			params[ctx.staff.Router.ParamsNameMap[newUri][d]] = nodes[i]

		}

		ctx.Uri = newUri
		//fmt.Println(newUri, "44444444444444444444444444444444")
		ctx.Params = params
	}
}

//对post请求进行解析
//由于post请求的参数不在url中，所以对传参为  请求
func (ctx *Context) ParsePost(r *http.Request) {
	contentType := strings.Split(r.Header["Content-Type"][0], "/")[1]
	//x-www-form-urlencoded传入的header的Content-Type中只有  x-www-form-urlencoded
	if contentType == "x-www-form-urlencoded" {
		//r.ParseForm()   是解析"x-www-form-urlencoded"类型文本的必要步骤
		_ = r.ParseForm()
		//将postForm表单以map形式赋值给params
		params := r.PostForm
		//遍历params
		for key, value := range params {
			ctx.Params[key] = value[0]
		}
		//form-data传入的header的Content-Type带有其他内容   form-data; boundary=--------------------------323814193179985924299079
	} else if strings.Contains(contentType, "form-data") {
		//整个请求体最大的内存限制
		_ = r.ParseMultipartForm(32 << 20)
		//同上  赋值遍历
		params := r.MultipartForm.Value
		for key, value := range params {
			ctx.Params[key] = value[0]
		}
	}
}

//这个方法   在调用的 时候   需要传入   的是     调用者的request    response   config
func Enter(ctx *Context) {
	r := ctx.Request
	w := ctx.ResponseWriter
	//通过请求拿到config
	config := ctx.staff.websocketMap[r.RequestURI]

	key := r.Header.Get("Sec-WebSocket-Key")

	s := sha1.New()
	s.Write(QuickStringToBytes(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	b := s.Sum(nil)
	accept := base64.StdEncoding.EncodeToString(b)

	hijack := w.(http.Hijacker)
	con, buf, _ := hijack.Hijack()

	ctx.Conn = con
	ctx.buf = buf
	//context := &WebSocketContext{
	//	Conn: con,
	//	buf:  buf,
	//}

	upgrade := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + accept + "\r\n\r\n"
	buf.Write(QuickStringToBytes(upgrade))
	buf.Flush()

	//config.OnOpen(context)
	config.OnOpen(ctx)

	go func() {
		for true {
			data := make([]byte, 2)
			_, err := buf.Read(data)
			//1，以byte[]的形式读出来，反过来，以byte[]的形式放回去
			if err != nil {
				//wsClose(con, config.OnClose, context)
				wsClose(con, config.OnClose, ctx)
				break
			}

			//fmt.Println(QuickBytesToString(data), "data...........")
			//fmt.Println(int(data[0]), "data[0]-----------")
			//fmt.Println(int(data[1]), "data[0]-----------")
			//2，将byte转int再转bool      反过来  将bool转int再转byte
			bin1 := parseIntToBin(int(data[0]))
			bin2 := parseIntToBin(int(data[1]))

			// bin1
			// 0    1    2     3   4 5 6 7
			// FIN RSV1 RSV2 RSV3 opcode(4)

			// bin2
			// 8     9 10 11 12 13 14 15
			// MASK      PayloadLen

			// RSV   这三个其中一个为1（true）  就关闭
			if bin1[1] || bin1[2] || bin1[3] {
				//wsClose(con, config.OnClose, context)
				wsClose(con, config.OnClose, ctx)

				break
			}

			//MASK   必须是1
			if !bin2[0] {
				//wsClose(con, config.OnClose, context)
				wsClose(con, config.OnClose, ctx)

				break
			}

			//操作代码，Opcode的值决定了应该如何解析后续的数据载荷（data payload）。如果操作代码是不认识的，那么接收端应该断开连接（fail the connection）。可选的操作代码如下：
			//
			//%x0：表示一个延续帧。当Opcode为0时，表示本次数据传输采用了数据分片，当前收到的数据帧为其中一个数据分片。
			//%x1：表示这是一个文本帧（frame）
			//%x2：表示这是一个二进制帧（frame）
			//%x3-7：保留的操作代码，用于后续定义的非控制帧。
			//%x8：表示连接断开。
			//%x9：表示这是一个ping操作。
			//%xA：表示这是一个pong操作。
			//%xB-F：保留的操作代码，用于后续定义的控制帧。
			opcode := parseBinToInt(bin1[4:])

			//数据长度
			payloadLen := parseBinToInt(bin2[1:])

			switch opcode {
			case 1:
				maskingKey := make([]byte, 4)
				buf.Read(maskingKey)
				//3，以byte[]读取
				payload := make([]byte, payloadLen)
				buf.Read(payload)

				data := make([]byte, payloadLen)
				for i := 0; i < payloadLen; i++ {
					data[i] = payload[i] ^ maskingKey[i%4]
				}
				fmt.Println(data, "data")
				//config.OnMessage(context, QuickBytesToString(data))
				//4，转string
				config.OnMessage(ctx, QuickBytesToString(data))

				strData:=QuickBytesToString(data)
				strData+="我太难了！"
				frame := Compose(strData)
				//ctx.ResponseWriter.Write(frame)
				buf.Write(frame)
				buf.Flush()
				//for i := range frame {
				//	fmt.Println(i, "kankan")
				//}

				//te := parseIntToBin(int(frame[8]))

				//for i := 0; i < len(te); i++ {
				//	fmt.Println(te[i], "bool")
				//}
				//buf.Write(frame)
				//buf.Writer.Write(frame)
				//w.Write(frame)
				//fmt.Println(len(data), "data.len")

			default:
				//wsClose(con, config.OnClose, context)
				wsClose(con, config.OnClose, ctx)
				break
			}
		}
	}()
	//return ""
}

//一个byte8个字节
