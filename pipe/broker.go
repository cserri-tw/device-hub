package pipe

import "github.com/thingful/expando"

type Broker interface {
	Channel() (Channel, error)
	Close() error
}

type Channel interface {
	Errors() chan error
	Out() chan expando.Input
}
