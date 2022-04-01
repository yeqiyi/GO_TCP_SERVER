package body

type Message struct{
	Uid string `json:"Uid"`
	Val string `json:"Val"`
}

type Resp struct {
    Uid string `json:"Uid"`
	Val string `json:"Val"`
    Ts string `json:"TS"`
}