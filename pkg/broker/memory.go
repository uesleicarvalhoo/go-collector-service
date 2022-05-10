package broker

type MemoryBroker struct {
	Events map[string][]Event
}

func NewMemoryBroker(topics ...CreateTopicInput) (*MemoryBroker, error) {
	client := &MemoryBroker{Events: make(map[string][]Event)}

	for _, topic := range topics {
		if err := client.DeclareTopic(topic); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (mb *MemoryBroker) SendEvent(event Event) error {
	mb.Events[event.Topic] = append(mb.Events[event.Topic], event)

	return nil
}

func (mb *MemoryBroker) DeclareTopic(payload CreateTopicInput) error {
	mb.Events[payload.Name] = []Event{}

	return nil
}
