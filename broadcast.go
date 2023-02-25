package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var messages = struct {
	mu     sync.Mutex
	values []int
}{}

var body = struct {
	Type    string `json:"type"`
	Message int    `json:"message"`
	MsgID   int    `json:"msg_id"`
}{}

func RunBroadcast() {
	n := maelstrom.NewNode()

	n.Handle("topology", func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]any{"type": "topology_ok"})
	})

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		messages.mu.Lock()
		messages.values = append(messages.values, body.Message)
		messages.mu.Unlock()

		return n.Send(msg.Src, map[string]any{"type": "broadcast_ok", "in_reply_to": body.MsgID})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		messages.mu.Lock()
		defer messages.mu.Unlock()

		return n.Reply(msg, map[string]any{
			"type":     "read_ok",
			"messages": messages.values,
		})
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
