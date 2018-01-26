package service

import "encoding/json"

type MessageFrame interface {
	ToByte() []byte
}

type Frame struct {
	Cmd       string      `json:"cmd"`
	Sessionid []string    `json:"sesstion_id"`
	Data      interface{} `json:"data"`
	Errno     string      `json:"errno"`
}

func (f *Frame) ToByte() []byte {
	msg, _ := json.Marshal(f)
	return msg
}

func (f *Frame) Read(p []byte) (int, error) {
	msg, err := json.Marshal(f)
	if err != nil {
		return 0, err
	}
	copy(p, msg)
	return len(msg), nil
}

func ToFrame(bytes []byte) *Frame {
	var frame Frame
	json.Unmarshal(bytes, &frame)
	return &frame
}

type ClientMessage struct {
	Cmd   string      `json:"cmd"`
	Data  interface{} `json:"data"`
	Errno string      `json:"errno"`
}

func (c *ClientMessage) ToByte() []byte {
	msg, _ := json.Marshal(c)
	return msg
}

func ToClientMessage(bytes []byte) *ClientMessage {
	var frame ClientMessage
	json.Unmarshal(bytes, &frame)
	return &frame
}

func ClientToFrame(c *ClientMessage, id string) *Frame {
	return &Frame{Cmd: c.Cmd, Data: c.Data, Sessionid: []string{id}, Errno: c.Errno}
}

func FrameToClient(f *Frame) *ClientMessage {
	return &ClientMessage{Cmd: f.Cmd, Data: f.Data, Errno: f.Errno}
}
