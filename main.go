package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"tcp_server/body"
	"time"
)

const TAG = "server: hello, "

func transfer(conn net.Conn) {
    defer func() {
        remoteAddr := conn.RemoteAddr().String()
        log.Print("discard remove add:", remoteAddr)
        conn.Close()
    }()

    // 设置1秒关闭连接
	conn.SetDeadline(time.Now().Add(10*time.Second))
    for {
        var msg body.Message

        if err := json.NewDecoder(conn).Decode(&msg); err != nil && err != io.EOF {
            log.Printf("Decode from client err: %v", err)
            // todo... 仿照redis协议写入err前缀符号`-`，通知client错误处理
            return
        }else{
			log.Println("receive message from ",conn.RemoteAddr().String()," : ",msg)
		}

        if msg.Uid != "" || msg.Val != "" {
            //conn.Write([]byte(msg.Val))
            var rsp body.Resp
            rsp.Uid = msg.Uid
            rsp.Val = "Hello,"+msg.Uid
            ser, _ := json.Marshal(rsp)

            conn.Write(append(ser, '\n'))
        }
    }
}

func ListenAndServer() {
    log.Print("Start server...")
    // 启动监听本地tcp端口3000
    listen, err := net.Listen("tcp", "0.0.0.0:3000")
    if err != nil {
        log.Fatal("Listen failed. msg: ", err)
        return
    }
    for {
		conn, err := listen.Accept()
		//resp:=make([]byte,1024)
		//_,err=conn.Read(resp)
		//log.Println("accep from ",conn.RemoteAddr().String()," : ",string(resp))
        if err != nil {
            log.Printf("accept failed, err: %v", err)
            continue
        }
        go transfer(conn)
    }
}

func main(){
	ListenAndServer()
}


