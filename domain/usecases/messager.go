package usecases

type Messager interface {
	SendMessage(queue, message string) error
}
