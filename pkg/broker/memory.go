package broker

type MemoryBroker struct {
	Events map[string][]Event
}

func NewMemoryBroker() *MemoryBroker {
	return &MemoryBroker{Events: make(map[string][]Event)}
}

func (mb *MemoryBroker) SendEvent(event Event) error {
	mb.Events[event.Topic] = append(mb.Events[event.Topic], event)

	return nil
}

func (mb *MemoryBroker) DeclareTopic(payload CreateTopicInput) error {
	mb.Events[payload.Name] = []Event{}

	return nil
}
