package handlers

import (
	"cqrs-playground/wishlist/add_item"
	"cqrs-playground/wishlist/remove_item"
	"encoding/json"
	"net/http"
)

type WishlistAPI struct {
	AddItemHandler    *add_item.AddItemHandler
	RemoveItemHandler *remove_item.RemoveItemHandler
}

func NewWishlistApi(
	addItemHandler *add_item.AddItemHandler,
	removeItemHandler *remove_item.RemoveItemHandler,
) *WishlistAPI {
	return &WishlistAPI{
		AddItemHandler:    addItemHandler,
		RemoveItemHandler: removeItemHandler,
	}
}

func (api *WishlistAPI) AddItem(w http.ResponseWriter, r *http.Request) {
	var cmd add_item.AddItemToWishlistCommandV2

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := api.AddItemHandler.HandleV2(cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "item added to wishlist"})
}

func (api *WishlistAPI) RemoveItem(w http.ResponseWriter, r *http.Request) {
	var cmd remove_item.RemoveItemFromWishlistCommandV1

	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := api.RemoveItemHandler.HandleV1(cmd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "item removed from wishlist"})
}
