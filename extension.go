package sirius

type Extension interface {
	Run(Message) (error, MessageAction)
}
