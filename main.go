package main

import (
	"context"
	"cqrs-playground/api/rest"
	"cqrs-playground/bibliothek/medien/ausleihen"
	"cqrs-playground/bibliothek/medien/bestand_projection"
	"cqrs-playground/bibliothek/medien/erwerben"
	"cqrs-playground/bibliothek/medien/isbn_index_projection"
	"cqrs-playground/bibliothek/medien/katalogisieren"
	"cqrs-playground/bibliothek/medien/rueckgeben"
	"cqrs-playground/bibliothek/medien/verliehen_projection"
	"cqrs-playground/bibliothek/nutzer/registrierung"
	db2 "cqrs-playground/db"
	"cqrs-playground/shared"
	cart_add_item "cqrs-playground/shopping_cart/add_item"
	cart_projection "cqrs-playground/shopping_cart/projection"
	cart_remove_item "cqrs-playground/shopping_cart/remove_item"
	wishlist_add_item "cqrs-playground/wishlist/add_item"
	wishlist_remove_item "cqrs-playground/wishlist/remove_item"
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

	isbnMediumIdProjection := isbn_index_projection.NewISBNIndexProjection(ctx, &eventStore, db.Pool, projectionUpdater, kafkaService)
	isbnMediumIdProjection.Start()

	mediumBestandProjection := bestand_projection.NewMediumBestandProjection(ctx, &eventStore, db.Pool, kafkaService)
	mediumBestandProjection.Start()

	mediumVerliehenProjection := verliehen_projection.NewMediumVerliehenProjection(ctx, &eventStore, db.Pool, kafkaService)
	mediumVerliehenProjection.Start()

	addItemHandlerCart := cart_add_item.NewAddItemHandler(ctx, &eventStore)
	removeItemHandlerCart := cart_remove_item.NewRemoveItemHandler(ctx, &eventStore)
	cartApi := rest.NewCartApi(addItemHandlerCart, removeItemHandlerCart)

	addItemHandlerWishlist := wishlist_add_item.NewAddItemHandler(ctx, &eventStore, kafkaService)
	removeItemHandlerWishlist := wishlist_remove_item.NewRemoveItemHandler(ctx, &eventStore, kafkaService)
	wishlistApi := rest.NewWishlistApi(addItemHandlerWishlist, removeItemHandlerWishlist)

	registriereNutzerHandler := registrierung.NewNutzerRegistrierungHandler(ctx, &eventStore, kafkaService)
	nutzerRegistrierungApi := rest.NewNutzerRegistrierungAPI(registriereNutzerHandler)

	erwerbeMediumHandler := erwerben.NewErwerbeMediumHandler(ctx, &eventStore, kafkaService, db.Pool)
	katalogisiereMediumHandler := katalogisieren.NewKatalogisiereMediumHandler(ctx, &eventStore, kafkaService)
	verleiheMediumHandler := ausleihen.NewVerleiheMediumHandler(ctx, &eventStore, kafkaService, db.Pool)
	rueckgabeMediumHandler := rueckgeben.NewMediumRueckgabeHandler(ctx, &eventStore, kafkaService, db.Pool)
	erwerbeMediumAPI := rest.NewErwerbeMediumAPI(
		erwerbeMediumHandler,
		katalogisiereMediumHandler,
		verleiheMediumHandler,
		rueckgabeMediumHandler)

	http.HandleFunc("/cart/add-item", cartApi.AddItem)
	http.HandleFunc("/cart/remove-item", cartApi.RemoveItem)
	http.HandleFunc("/wishlist/add-item", wishlistApi.AddItem)
	http.HandleFunc("/wishlist/remove-item", wishlistApi.RemoveItem)

	http.HandleFunc("/bibliothek/registriere-nutzer", nutzerRegistrierungApi.RegistriereNutzer)
	http.HandleFunc("/bibliothek/erwerbe-medium", erwerbeMediumAPI.ErwerbeMedium)
	http.HandleFunc("/bibliothek/katalogisiere-medium", erwerbeMediumAPI.KatalogisiereMedium)
	http.HandleFunc("/bibliothek/verleihe-medium", erwerbeMediumAPI.VerleiheMedium)
	http.HandleFunc("/bibliothek/gebe-medium-zurueck", erwerbeMediumAPI.GebeMediumZurueck)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
