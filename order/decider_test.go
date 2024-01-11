package order

import "testing"

func given(events []Event, command Command) []Event {
	return Decide(command, Fold(events))
}

func TestOrderDecider(t *testing.T) {
	createOrder := CreateOrder{}
	orderCreated := OrderCreated{}

	t.Run("it creates an order", func(t *testing.T) {
		events := given([]Event{}, createOrder)
		if len(events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(events))
		}
		event := events[0]
		switch event.(type) {
		case OrderCreated:
			// pass
		default:
			t.Errorf("Expected OrderCreated, got %T", event)

		}
	})

	t.Run("creating an order is idempotent", func(t *testing.T) {
		events := given([]Event{orderCreated}, createOrder)
		if len(events) != 0 {
			t.Errorf("Expected 0 events, got %d", len(events))
		}
	})

	addFirstItemToOrder := AddItemToOrder{
		itemId:    1,
		productId: 2,
		cost:      10,
	}
	firstItemAddedToOrder := ItemAddedToOrder{
		itemId:    1,
		productId: 2,
		cost:      10,
	}

	t.Run("cannot add an item before creating an order", func(t *testing.T) {
		events := given([]Event{}, addFirstItemToOrder)
		if len(events) != 0 {
			t.Errorf("Expected 0 events, got %d", len(events))
		}
	})

	t.Run("adding an item to an order", func(t *testing.T) {
		events := given([]Event{orderCreated}, addFirstItemToOrder)
		if len(events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(events))
		}
		event := events[0]
		if event != firstItemAddedToOrder {
			t.Errorf("Expected ItemAddedToOrder, got %T", event)
		}
	})

	t.Run("adding the same item to an order is idempotent", func(t *testing.T) {
		events := given([]Event{orderCreated, firstItemAddedToOrder}, addFirstItemToOrder)
		if len(events) != 0 {
			t.Errorf("Expected 0 events, got %d", len(events))
		}
	})

	addSecondItemToOrder := AddItemToOrder{
		itemId:    2,
		productId: 3,
		cost:      20,
	}
	secondItemAddedToOrder := ItemAddedToOrder{
		itemId:    2,
		productId: 3,
		cost:      20,
	}
	t.Run("adding a second item to an order", func(t *testing.T) {
		events := given([]Event{orderCreated, firstItemAddedToOrder}, addSecondItemToOrder)
		if len(events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(events))
		}
		event := events[0]
		if event != secondItemAddedToOrder {
			t.Errorf("Expected ItemAddedToOrder, got %T", event)
		}
	})

	finaliseOrder := FinaliseOrder{}
	orderFinalised := OrderFinalised{}
	t.Run("finalising the order", func(t *testing.T) {
		events := given([]Event{orderCreated, firstItemAddedToOrder}, finaliseOrder)
		if len(events) != 1 {
			t.Errorf("Expected 1 event, got %d", len(events))
		}
		event := events[0]
		if event != orderFinalised {
			t.Errorf("Expected OrderFinalised, got %T", event)
		}
	})

	t.Run("finalising the order is idempotent", func(t *testing.T) {
		events := given([]Event{orderCreated, firstItemAddedToOrder, orderFinalised}, finaliseOrder)
		if len(events) != 0 {
			t.Errorf("Expected 0 events, got %d", len(events))
		}
	})

	t.Run("cannot add an item after finalising an order", func(t *testing.T) {
		events := given([]Event{orderCreated, firstItemAddedToOrder, orderFinalised}, addSecondItemToOrder)
		if len(events) != 0 {
			t.Errorf("Expected 0 events, got %d", len(events))
		}
	})
}
