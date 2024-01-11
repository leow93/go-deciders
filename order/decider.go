package order

type Item struct {
	itemId    int
	productId int
	cost      int
}

type CreateOrder struct{}
type AddItemToOrder = Item
type FinaliseOrder = struct{}
type Command interface{}

type OrderCreated struct{}
type ItemAddedToOrder = Item
type OrderFinalised struct{}
type Event interface{}

type State struct {
	lineItems []Item
	status    string // 'initial' | 'created' | 'finalised'. I would love a union type but...
}

func initialState() State {
	return State{
		lineItems: make([]Item, 0),
		status:    "initial",
	}
}

func Evolve(state State, event interface{}) State {
	switch event.(type) {
	case OrderCreated:
		return State{
			lineItems: make([]Item, 0),
			status:    "created",
		}

	case ItemAddedToOrder:
		return State{
			lineItems: append(state.lineItems, event.(ItemAddedToOrder)),
			status:    state.status,
		}
	case OrderFinalised:
		return State{
			lineItems: state.lineItems,
			status:    "finalised",
		}
	default:
		return state
	}
}

func containsItem(items []Item, item Item) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

func Decide(command Command, state State) []Event {
	events := make([]Event, 0)
	switch command.(type) {
	case CreateOrder:
		{
			if state.status == "created" {
				return events
			}
			return append(events, OrderCreated{})
		}
	case AddItemToOrder:
		{
			if state.status != "created" || containsItem(state.lineItems, command.(AddItemToOrder)) {
				return events
			}
			return append(events, ItemAddedToOrder{
				itemId:    command.(AddItemToOrder).itemId,
				productId: command.(AddItemToOrder).productId,
				cost:      command.(AddItemToOrder).cost,
			})
		}
	case FinaliseOrder:
		if state.status != "created" {
			return events
		}
		return append(events, OrderFinalised{})
	default:
		return events
	}
}

func reduce(state State, events []Event) State {
	for _, event := range events {
		state = Evolve(state, event)
	}
	return state
}

func Fold(events []Event) State {
	return reduce(initialState(), events)
}
