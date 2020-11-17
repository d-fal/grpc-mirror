package broker

// Broker is the struct that passes all the messages
type Broker struct {
	StopCh        chan struct{}
	PublishCh     chan interface{}
	SubscribeCh   chan chan interface{}
	UnsubscribeCh chan chan interface{}
}

// SetupBroker Setup a broker
func SetupBroker() *Broker {
	return &Broker{
		StopCh:        make(chan struct{}),
		PublishCh:     make(chan interface{}, 1),
		SubscribeCh:   make(chan chan interface{}, 1), // blocking channel
		UnsubscribeCh: make(chan chan interface{}, 1), // blocking channel
	}
}

// Start starts the broker listening to subscriptions
func (brkr *Broker) Start() {
	subscribers := map[chan interface{}]struct{}{}

	for {
		select {

		case <-brkr.StopCh:
			return

		case newChannel := <-brkr.SubscribeCh:
			subscribers[newChannel] = struct{}{}

		case message := <-brkr.PublishCh:
			for s := range subscribers {
				select {
				case s <- message:
				default:
				}
			}
		case s := <-brkr.UnsubscribeCh:
			delete(subscribers, s)
		}

	}
}

// Unsubscribe unsubscribes from a channel
func (brkr *Broker) Unsubscribe(s chan interface{}) {
	brkr.UnsubscribeCh <- s
}

// Stop broker
func (brkr *Broker) Stop() {
	close(brkr.StopCh)
}

// Subscribe to channel
func (brkr *Broker) Subscribe() chan interface{} {
	newChannel := make(chan interface{}, 10)
	brkr.SubscribeCh <- newChannel

	return newChannel
}

// Publish publishes messages
func (brkr *Broker) Publish(message interface{}) {
	brkr.PublishCh <- message
}
