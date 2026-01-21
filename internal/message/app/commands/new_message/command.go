package newmessage

import "io"

type NewMessageCommand struct {
	Stream    io.Reader
	NodeID    string
	MediaType string
}
