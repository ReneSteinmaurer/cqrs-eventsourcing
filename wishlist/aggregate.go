package wishlist

import (
	"cqrs-playground/shared"
	"cqrs-playground/wishlist/events"
	"encoding/json"
	"errors"
)

var MaxSize = 2

type WishlistAggregate struct {
	Id     int
	UserId string
	Items  map[string]bool
}

func NewWishlistAggregateFrom(events []shared.Event) *WishlistAggregate {
	w := &WishlistAggregate{
		Items: make(map[string]bool),
	}
	for _, event := range events {
		w.Apply(event)
	}
	return w
}

func (w *WishlistAggregate) Apply(event shared.Event) {
	switch event.Type {
	case events.ItemAddedToWishlistEventTypeV2:
		var e events.ItemAddedToWishlistEventV2
		_ = json.Unmarshal(event.Payload, &e)
		w.Id = e.WishlistId
		w.UserId = e.UserId
		w.Items[e.Item] = true
	}
}

func (w *WishlistAggregate) HandleAddItem(cmd events.ItemAddedToWishlistEventV2) error {
	if !w.canAdd(cmd.Item) {
		return errors.New("item already exists in wishlist")
	}
	if w.hasReachedMaxItems() {
		return errors.New("wishlist has exceeded maximum possible items")
	}
	return nil
}

func (w *WishlistAggregate) canAdd(item string) bool {
	return !w.Items[item]
}

func (w *WishlistAggregate) hasReachedMaxItems() bool {
	return len(w.Items) >= MaxSize
}
