package mutexmaps

func (e *Event) Register(groupID int, allocatedChan interface{}) *Subscriber {
	g, ok := (*e)[groupID]
	if !ok {
		g = newGroup()
		(*e)[groupID] = g
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	sub := &Subscriber{
		Ch:    allocatedChan,
		group: g,
		next:  g.Root.next,
		prev:  g.Root,
	}

	if g.Root.next != nil {
		g.Root.next.prev = sub
	}
	g.Root.next = sub

	return sub
}

func (e *EventHandler) Unsubscribe(sub **Subscriber) {
	(*sub).group.mu.Lock()
	defer (*sub).group.mu.Unlock()

	(*sub).prev.next = (*sub).next
	if (*sub).next != nil {
		(*sub).next.prev = (*sub).prev
	}
	*sub = nil // fix - gc moment
}

func (e *Event) getGroup(groupID int) *subGroup {
	g, ok := (*e)[groupID]
	if !ok {
		g = newGroup()
		(*e)[groupID] = g
	}

	return g
}
