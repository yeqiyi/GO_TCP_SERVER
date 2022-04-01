package conn

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"net"
	"sync"
	"tcp_server/body"
	"time"
)

type IConn interface {
	Close() error
}

type Conn struct {
	addr    string
	tcp     *net.TCPConn
	ctx     context.Context
	writer  *bufio.Writer
	cnlFunc context.CancelFunc //通知ctx结束
	retChan *sync.Map          //存放通道结果集合的map，属于统一连接
	err     error
}

type Option struct {
	Addr        string
	Size        int //缓冲区大小
	IsKeepAlive bool
	ReadTimeout time.Duration
	DialTimeout time.Duration
	KeepAlive   time.Duration
}

func (c *Conn) Close() (err error) {
	if c.cnlFunc != nil {
		c.cnlFunc()
	}

	if c.tcp != nil {
		err = c.tcp.Close()
	}

	if c.retChan != nil {
		c.retChan.Range(func(key, value interface{}) bool {
			if ch, ok := value.(chan string); ok {
				close(ch)
			}
			return true
		})
	}
	return
}

func (c *Conn) Send(ctx context.Context,msg *body.Message)(ch chan string,err error){
	ch=make(chan string)
	c.retChan.Store(msg.Uid,ch)
	js,_:=json.Marshal(msg)
	_,err=c.writer.Write(js)
	if err!=nil{
		return
	}
	err=c.writer.Flush()
	return
}

func NewConn(opt *Option) (c *Conn, err error) {
	c = &Conn{
		addr:    opt.Addr,
		retChan: new(sync.Map),
	}
	var conn net.Conn
	if conn, err = net.DialTimeout("tcp", opt.Addr, opt.DialTimeout); err != nil {
		return
	} else {
		c.tcp = conn.(*net.TCPConn)
	}

	c.writer = bufio.NewWriter(c.tcp)

	if opt.IsKeepAlive {
		if err = c.tcp.SetKeepAlive(true); err != nil {
			return
		}
		if err = c.tcp.SetKeepAlivePeriod(opt.KeepAlive); err != nil {
			return
		}
	} else {
		if err = c.tcp.SetKeepAlive(false); err != nil {
			return
		}
	}
	//丢弃未发送或未确认数据包
	if err = c.tcp.SetLinger(0); err != nil {
		return
	}

	c.ctx, c.cnlFunc = context.WithCancel(context.Background())
	go receiveResp(c)
	return
}

func receiveResp(c *Conn){
	scanner:=bufio.NewScanner(c.tcp)
	for{
		select{
		case<-c.ctx.Done():
			return
		default:
			if scanner.Scan(){
				resp:=new(body.Resp)
				if err:=json.Unmarshal(scanner.Bytes(),resp);err!=nil{
					return
				}
				//log.Println("resp=",resp)
				uid:=resp.Uid
				if load,ok:=c.retChan.Load(uid);ok{
					c.retChan.Delete(load)
					
					if ch,ok:=load.(chan string);ok{
						ch<-resp.Val
						close(ch)
					}
				}
			}else{
				if scanner.Err()!=nil{
					c.err=scanner.Err()
				}else{
					c.err=errors.New("scanner done")
				}
				c.Close()
				return
			}
		}
	}
}
