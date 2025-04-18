package main

import (
	"context"
	"cqrs-playground/api/rest"
	db2 "cqrs-playground/db"
	"cqrs-playground/shared"
	cart_add_item "cqrs-playground/shopping-cart/add_item"
	cart_projection "cqrs-playground/shopping-cart/projection"
	cart_remove_item "cqrs-playground/shopping-cart/remove_item"
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
	//publisher := kafkaService.NewSyncProducer()

	/*go func() {
		partition := kafkaService.NewConsumerOffsetNewest("test")

		for message := range partition.Messages() {
			log.Printf("Nachricht empfangen: Partition=%d, Offset=%d, Key=%s, Value=%s\n",
				message.Partition, message.Offset, string(message.Key), string(message.Value))
		}
	}()

	err := kafkaService.SendEvent(publisher, "test", "Hallo oida")
	if err != nil {
		panic(err)
	}

	err = kafkaService.SendEvent(publisher, "test", "Test servas servas")
	if err != nil {
		panic(err)
	}*/

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
