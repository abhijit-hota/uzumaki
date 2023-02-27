package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type MessageBody[T comparable] struct {
	Type    string `json:"type"`
	Message T      `json:"message"`
	MsgID   *int   `json:"msg_id"`
}

func RunMultiBroadcast() {
	messages := NewSafeSet[int]()

	neighbours := make([]string, 0)

	n := maelstrom.NewNode()

	n.Handle("topology", func(msg maelstrom.Message) error {
		topologyBody := struct {
			Topology map[string][]string `json:"topology"`
		}{}

		if err := json.Unmarshal(msg.Body, &topologyBody); err != nil {
			return err
		}

		neighbours = topologyBody.Topology[n.ID()]

		return n.Reply(msg, map[string]any{"type": "topology_ok"})
	})

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body MessageBody[int]
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		if newAdded := messages.Add(body.Message); newAdded {
			for _, id := range neighbours {
				// I don't understand why n.Send() doesn't work.
				n.RPC(id, map[string]any{
					"type":    "broadcast",
					"message": body.Message,
				}, func(msg maelstrom.Message) error {
					return nil
				})
			}
		}

		if body.MsgID == nil {
			return nil
		}

		return n.Reply(msg, map[string]any{"type": "broadcast_ok"})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		return n.Reply(msg, map[string]any{
			"type":     "read_ok",
			"messages": messages.Values(),
		})
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
