package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/narqo/fly-dist-sys/internal"
)

func main() {
	svc := &UniqueIDSvc{}
	err := internal.Run(svc)
	if err != nil {
		log.Fatal(err)
	}
}

type GenerateOKMessage struct {
	ID string `json:"id"`
}

func (m GenerateOKMessage) MessageType() string {
	return "generate_ok"
}

type UniqueIDSvc struct {
	nodeID string
	epoch  int64
	msgID  uint
}

func (svc *UniqueIDSvc) Init(init map[string]any) error {
	svc.nodeID = init["node_id"].(string)
	svc.epoch = time.Now().UnixNano()
	svc.msgID = 1
	return nil
}

func (svc *UniqueIDSvc) Handle(req internal.RPCMessage, w *json.Encoder) error {
	typ := req.MessageType()
	switch typ {
	case "generate":
		msg := GenerateOKMessage{
			ID: fmt.Sprintf("%v-%d-%d", svc.nodeID, svc.epoch, svc.msgID),
		}
		svc.msgID++
		reply := internal.ReplyTo(req, svc.msgID, msg)
		return w.Encode(reply)
	case "generate_ok":
		return nil
	default:
		return fmt.Errorf("unexpected message type %q", typ)
	}
}
