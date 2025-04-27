package handlers

import (
	"cqrs-playground/shopping_cart/add_item"
	"cqrs-playground/shopping_cart/remove_item"
	"encoding/json"
	"net/http"
)

type CartAPI struct {
	AddItemHandler    *add_item.AddItemHandler
	RemoveItemHandler *remove_item.RemoveItemHandler
}

func NewCartApi(
	addItemHandler *add_item.AddItemHandler,
	removeItemHandler *remove_item.RemoveItemHandler,
) *CartAPI {
	return &CartAPI{
		AddItemHandler:    addItemHandler,
		RemoveItemHandler: removeItemHandler,
	}
}

func (api *CartAPI) AddItem(w http.ResponseWriter, r *http.Request) {
	var cmd add_item.AddItemToCartCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := api.AddItemHandler.Handle(cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "item added to cart"})
}

func (api *CartAPI) RemoveItem(w http.ResponseWriter, r *http.Request) {
	var cmd remove_item.RemoveItemFromCartCommand

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := api.RemoveItemHandler.Handle(cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "item removed from cart"})
}
