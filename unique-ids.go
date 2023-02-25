package main

import (
	"log"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func RunUniqueIDs() {
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {

		id := uuid.NewString()
		resp := map[string]string{
			"type": "generate_ok",
			"id":   id,
		}
		return n.Reply(msg, resp)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
