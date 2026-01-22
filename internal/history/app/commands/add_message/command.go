package addmessage

type AddMessageCommand struct {
	UserID  string
	Message map[string]interface{}
	Limit   uint8
}
