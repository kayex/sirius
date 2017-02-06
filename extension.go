package sirius

type Extension interface {
	Run(Message) (error, MessageAction)
}

type EID string

type ExtensionLoader interface {
	Load(EID) Extension
}
