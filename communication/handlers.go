package communication

import (
	"net"

	"github.com/superc03/colirc/data"
	"github.com/superc03/colirc/types"
)

// Handler represents a function meant to handle its corresponding command
type Handler func(net.Conn, *data.Client)

// Map of Corresponding
var Handlers = map[string]Handler{
	types.CommandJOIN: func(c net.Conn) {
		c.Write([]byte("Hello True Believer\n"))
	},
}

func BadRequestError(conn net.Conn) {
	msg, _ := data.MarshalMessage(&data.Message{
		Command: types.ErrorUnknownCommand,
		Params:  []string{":Bad Request"},
	})
	conn.Write([]byte(msg))
}

func UnknownCommandError(conn net.Conn, command string) {
	msg, _ := data.MarshalMessage(&data.Message{
		Command: types.ErrorUnknownCommand,
		Params:  []string{command, ":Unknown command"},
	})
	conn.Write([]byte(msg))
}
