package ports

type Messager interface {
	SendMessage(queue, message string) error
}
