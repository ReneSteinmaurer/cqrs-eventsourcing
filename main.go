package main

import (
	"context"
	"cqrs-playground/api/rest"
	db2 "cqrs-playground/db"
	"cqrs-playground/shared"
	cart_add_item "cqrs-playground/shopping_cart/add_item"
	cart_projection "cqrs-playground/shopping_cart/projection"
	cart_remove_item "cqrs-playground/shopping_cart/remove_item"
	wishlist_add_item "cqrs-playground/wishlist/add_item"
	"cqrs-playground/wishlist/wishlist_projection"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	db := db2.NewDB(ctx)
	kafkaService := shared.NewKafkaService()

	eventStore := shared.EventStore{DB: db.Pool}

	projectionUpdater := shared.NewProjectionStateUpdater(ctx, db.Pool)

	cartProjection := cart_projection.NewCartProjection(ctx, &eventStore, db.Pool, projectionUpdater, kafkaService)
	go cartProjection.Start(100 * time.Millisecond)

	wishlistProjection := wishlist_projection.NewWishlistProjection(ctx, &eventStore, db.Pool, projectionUpdater, kafkaService)
	wishlistProjection.Start()

	addItemHandlerCart := cart_add_item.NewAddItemHandler(ctx, &eventStore)
	removeItemHandlerCart := cart_remove_item.NewRemoveItemHandler(ctx, &eventStore)
	cartApi := rest.NewCartApi(addItemHandlerCart, removeItemHandlerCart)

	addItemHandlerWishlist := wishlist_add_item.NewAddItemHandler(ctx, &eventStore, kafkaService)
	wishlistApi := rest.NewWishlistApi(addItemHandlerWishlist)

	http.HandleFunc("/cart/add-item", cartApi.AddItem)
	http.HandleFunc("/cart/remove-item", cartApi.RemoveItem)
	http.HandleFunc("/wishlist/add-item", wishlistApi.AddItem)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
