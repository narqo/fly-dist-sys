package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/narqo/fly-dist-sys/internal"
)

func main() {
	svc := &EchoSvc{}
	err := internal.Run(svc)
	if err != nil {
		log.Fatal(err)
	}
}

type EchoOKMessage struct {
	Echo string `json:"echo"`
}

func (m EchoOKMessage) MessageType() string {
	return "echo_ok"
}

type EchoSvc struct {
	id uint
}

func (svc *EchoSvc) Init(map[string]any) error {
	svc.id = 1
	return nil
}

func (svc *EchoSvc) Handle(req internal.RPCMessage, w *json.Encoder) error {
	typ := req.MessageType()
	switch typ {
	case "echo":
		msg := EchoOKMessage{
			Echo: req.Body["echo"].(string),
		}
		reply := internal.ReplyTo(req, svc.id, msg)
		svc.id++
		return w.Encode(reply)
	case "echo_ok":
		return nil
	default:
		return fmt.Errorf("unexpected message type %q", typ)
	}
}
