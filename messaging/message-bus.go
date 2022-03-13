package messaging

type MessageBus interface {
	CreateRabbitMQWorker() MessageBusConsumer
	SendMessageToRabbitMQWorkQueue(msg string)
	//Close()
}

type MessageBusConsumer interface {
	OnMessageReceived(f func(msg string) error)
	StartListening() error
}
