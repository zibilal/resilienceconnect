package resilienceconnect

type Connector interface {
	Connect(request Requestor, options ConnectorOption, output interface{}) (Responder, error)
}

type Responder interface {
	StatusCode() int
	StatusMessage() string
}

type Requestor interface {
	Request(interface{}) error
}

type ConnectorOption map[string]interface{}

func (o ConnectorOption) Get(key string) interface{} {
	return o[key]
}

func (o ConnectorOption) Put(key string, value interface{}) {
	o[key] = value
}

type ConnectionFunc func(request Requestor, output interface{}) (Responder, error)