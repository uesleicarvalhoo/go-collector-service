package broker

type MemoryBroker struct {
	Events map[string][]Event
}

func NewMemoryBroker() (*MemoryBroker, error) {
	client := &MemoryBroker{Events: make(map[string][]Event)}

	return client, nil
}

func (mb *MemoryBroker) SendEvent(event Event) error {
	if _, ok := mb.Events[event.Topic]; !ok {
		mb.Events[event.Topic] = []Event{}
	}

	mb.Events[event.Topic] = append(mb.Events[event.Topic], event)

	return nil
}
