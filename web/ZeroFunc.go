package web

import (
	"math"
	"net"
	"strings"
	"unsafe"
)

//这里面的func    都没有对象



//返回一个staff对象  用于调用最表面的方法
func Default() *Staff {
	staff := Init()
	return staff
}

//初始化staff    给出默认值
func Init() *Staff {
	staff := &Staff{
		websocketMap: map[string]WebSocketConfig{},
		Router: Router{
			//至少得有一个"/"    http.HandleFunc才能运行
			Path:    "/",
			Handle:  nil,
			Routers: map[string]Handle{},
			Method:  "",
			NodesMap: map[int]Nodes{
				0: map[int]string{},
				1: map[int]string{},
				2: map[int]string{},
				3: map[int]string{},
				4: map[int]string{},
				5: map[int]string{},
				6: map[int]string{},
				7: map[int]string{},
				8: map[int]string{},
				9: map[int]string{},
			},
			index: 0,
			ParamsMap: map[string]Params{
				"0": map[string]string{},
				"1": map[string]string{},
				"2": map[string]string{},
				"3": map[string]string{},
				"4": map[string]string{},
				"5": map[string]string{},
				"6": map[string]string{},
				"7": map[string]string{},
				"8": map[string]string{},
				"9": map[string]string{},
			},
			ParamsNameMap: map[string]ParamsName{
				"2":nil,
			},
		},
	}
	//staff.Router.InitMap()
	staff.Router.staff = staff
	staff.pool.New = func() interface{} {
		return staff.InitContext()
	}
	return staff
}

//对uri进行简化保存
func SimplyUri(uri string) string {
	uris := strings.Split(uri, "/")
	newUri := ""
	for i, _ := range uris {
		if i > 0 {
			newUri += "/"
			if strings.Contains(uris[i], ":") {
				newUri += "-"
			} else {
				newUri += uris[i]
			}
		}
	}
	return newUri
}


func wsClose(conn net.Conn, onclose func(ctx *Context), ctx *Context) {
	onclose(ctx)
	conn.Close()
}

//websocket解帧
func parseBinToInt(b []bool) (res int) {
	pos := 0
	for i := len(b) - 1; i >= 0; i-- {
		if b[i] {
			res += int(math.Pow(float64(2), float64(pos)))
		}
		pos++
	}
	return
}

func parseIntToBin(i int) (b []bool) {
	b = make([]bool, 8)
	pos := 7
	for i != 0 {
		add := false
		tmp := i % 2
		i = i / 2
		if tmp == 1 {
			add = true
		}
		b[pos] = add
		pos--
	}
	return
}

func QuickBytesToString(b []byte) (s string) {
	return *(*string)(unsafe.Pointer(&b))
}

func QuickStringToBytes(s string) (b []byte) {
	return *(*[]byte)(unsafe.Pointer(&s))
}


//将需要发送的数据组装成数据帧    这里  只对125及其以下的做了组装
func Compose(data string) []byte {
	frame := make([]byte, 2)
	//前8个二进制
	frame[0] = byte(int(129))
	//后8个二进制    理论上小于126   否则mask会变成1  然后自动中断通讯
	if len(data)>125 {
		return []byte("")
	}
	frame[1] = byte(len(data))
	//将数据内容变成byte，利用循环加到桢里面
	dataByte:=QuickStringToBytes(data)
	for i:=0;i< len(dataByte);i++  {
		frame=append(frame, dataByte[i])
	}
	return frame
}
