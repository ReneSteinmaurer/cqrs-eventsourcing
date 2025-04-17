package rest

import (
	"cqrs-playground/wishlist/add_item"
	"encoding/json"
	"net/http"
)

type WishlistAPI struct {
	AddItemHandler *add_item.AddItemHandler
}

func NewWishlistApi(
	addItemHandler *add_item.AddItemHandler,
) *WishlistAPI {
	return &WishlistAPI{
		AddItemHandler: addItemHandler,
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
