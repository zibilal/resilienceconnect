package resilienceconnect

type Connector interface {
	Connect(resource string, options ConnectorOption, input interface{}, output interface{}) error
}

type Requestor interface {
	Request(interface{}) error
}

type ConnectorOption map[string]interface{}
func(o ConnectorOption) Get(key string) interface{} {
	return o[key]
}

func(o ConnectorOption) Put(key string, value interface{}) {
	o[key] = value
}

type ConnectionFunc func (request Requestor, output interface{}) error