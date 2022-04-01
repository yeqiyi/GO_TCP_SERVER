package main

import (
	"context"
	"tcp_server/body"
	. "tcp_server/conn"
	"testing"
	"time"
)

var Opt=&Option{
	Addr: "127.0.0.1:3000",
	Size: 3,
	ReadTimeout: 3*time.Second,
	DialTimeout: 3*time.Second,
	IsKeepAlive: true,
	KeepAlive: 1*time.Second,
}

func createConn(opt *Option)*Conn{
	c,err:=NewConn(opt)
	if err!=nil{
		panic(err)
	}
	return c
}


func TestServer(t *testing.T){
	c:=createConn(Opt)
	msg:=&body.Message{Uid: "client-1",Val: "hello!"}
	rec,err:=c.Send(context.Background(),msg)
	if err!=nil{
		t.Error(err)
	}else{
		t.Logf("rec1: %+v",<-rec)
	}
	time.Sleep(2*time.Second)
	msg.Val = "another pig!"
    rec2, err := c.Send(context.Background(), msg)
    if err != nil {
        t.Error(err)
    } else {
        t.Logf("rec2: %+v", <-rec2)
    }
    t.Log("finished")
}