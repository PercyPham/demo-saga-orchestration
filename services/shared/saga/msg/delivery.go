package msg

type Delivery struct {
	Message Message
	Ack     func()
	Nack    func()
}
