package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

func Run(svc Service) error {
	in := json.NewDecoder(os.Stdin)
	out := json.NewEncoder(os.Stdout)

	var init RPCMessage
	err := in.Decode(&init)
	if err != nil {
		return fmt.Errorf("decode init message: %w", err)
	}
	if typ := init.MessageType(); typ != "init" {
		return fmt.Errorf("expected init message, got %q", typ)
	}

	err = svc.Init(init.Body)
	if err != nil {
		return fmt.Errorf("init service: %w", err)
	}

	reply := ReplyTo(init, 0, InitOKMessage{})
	err = out.Encode(reply)
	if err != nil {
		return fmt.Errorf("handle init message: %w", err)
	}

	for {
		var req RPCMessage
		err := in.Decode(&req)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}

		err = svc.Handle(req, out)
		if err != nil {
			return fmt.Errorf("handle request: %s -> %s: %w", req.Src, req.Dst, err)
		}
	}
}

type Service interface {
	Init(init map[string]any) error
	Handle(req RPCMessage, w *json.Encoder) error
}

type Message interface {
	MessageType() string
}

type RPCMessage struct {
	Src  string         `json:"src"`
	Dst  string         `json:"dest"`
	Body map[string]any `json:"body"`
}

func (m RPCMessage) MessageType() string {
	typ, _ := m.Body["type"].(string)
	return typ
}

func ReplyTo(req RPCMessage, msgID uint, msg Message) RPCMessage {
	body := map[string]any{
		"type":        msg.MessageType(),
		"msg_id":      msgID,
		"in_reply_to": req.Body["msg_id"],
	}
	reply := RPCMessage{
		Src:  req.Dst,
		Dst:  req.Src,
		Body: body,
	}

	buf, _ := json.Marshal(msg)
	v := make(map[string]any, 0)
	_ = json.Unmarshal(buf, &v)
	for k, vv := range v {
		reply.Body[k] = vv
	}

	return reply
}

type InitOKMessage struct{}

func (m InitOKMessage) MessageType() string {
	return "init_ok"
}
