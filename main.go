package main

import (
	"context"
	"cqrs-playground/api/rest"
	db2 "cqrs-playground/db"
	"cqrs-playground/shared"
	"cqrs-playground/shopping-cart/add_item"
	"cqrs-playground/shopping-cart/projection"
	"cqrs-playground/shopping-cart/remove_item"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	db := db2.NewDB(ctx)

	eventStore := shared.EventStore{DB: db.Pool}

	cartProjection := projection.NewCartProjection(ctx, &eventStore, db.Pool)
	go cartProjection.Start(100 * time.Millisecond)

	addItemHandler := add_item.NewAddItemHandler(ctx, &eventStore)
	removeItemHandler := remove_item.NewRemoveItemHandler(ctx, &eventStore)
	api := rest.NewCartApi(addItemHandler, removeItemHandler)

	http.HandleFunc("/cart/add-item", api.AddItem)
	http.HandleFunc("/cart/remove-item", api.RemoveItem)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
