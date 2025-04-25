package main

import (
	"context"
	"cqrs-playground/api/rest"
	"cqrs-playground/api/ws"
	"cqrs-playground/bibliothek/medien/ausleihen"
	"cqrs-playground/bibliothek/medien/erwerben"
	"cqrs-playground/bibliothek/medien/katalogisieren"
	"cqrs-playground/bibliothek/medien/projections/bestand"
	"cqrs-playground/bibliothek/medien/projections/isbn_index"
	"cqrs-playground/bibliothek/medien/projections/verliehen"
	"cqrs-playground/bibliothek/medien/rueckgeben"
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
	webSocketHub := shared.NewWebSocketHub(ctx)
	webSocketApi := ws.NewNotificationWsAPI(ctx, webSocketHub)
	notificationService := shared.NewNotificationService(ctx, webSocketHub)

	eventStore := shared.EventStore{DB: db.Pool}

	projectionUpdater := shared.NewProjectionStateUpdater(ctx, db.Pool)

	cartProjection := cart_projection.NewCartProjection(ctx, &eventStore, db.Pool, projectionUpdater, kafkaService)
	go cartProjection.Start(100 * time.Millisecond)

	wishlistProjection := wishlist_projection.NewWishlistProjection(ctx, &eventStore, db.Pool, projectionUpdater, kafkaService)
	wishlistProjection.Start()

	isbnMediumIdProjection := isbn_index.NewISBNIndexProjection(ctx, &eventStore, db.Pool, projectionUpdater, kafkaService)
	isbnMediumIdProjection.Start()

	mediumBestandProjection := bestand.NewMediumBestandProjection(ctx, &eventStore, db.Pool, kafkaService, notificationService)
	mediumBestandProjection.Start()

	mediumVerliehenProjection := verliehen.NewMediumVerliehenProjection(ctx, &eventStore, db.Pool, kafkaService)
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

	mediumBestandAPI := rest.NewMediumBestandAPI(db.Pool)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", webSocketApi.Handle)

	mux.HandleFunc("/cart/add-item", cartApi.AddItem)
	mux.HandleFunc("/cart/remove-item", cartApi.RemoveItem)
	mux.HandleFunc("/wishlist/add-item", wishlistApi.AddItem)
	mux.HandleFunc("/wishlist/remove-item", wishlistApi.RemoveItem)

	mux.HandleFunc("/bibliothek/registriere-nutzer", nutzerRegistrierungApi.RegistriereNutzer)
	mux.HandleFunc("/bibliothek/erwerbe-medium", erwerbeMediumAPI.ErwerbeMedium)
	mux.HandleFunc("/bibliothek/katalogisiere-medium", erwerbeMediumAPI.KatalogisiereMedium)
	mux.HandleFunc("/bibliothek/verleihe-medium", erwerbeMediumAPI.VerleiheMedium)
	mux.HandleFunc("/bibliothek/gebe-medium-zurueck", erwerbeMediumAPI.GebeMediumZurueck)

	mux.HandleFunc("/bibliothek/bestand", mediumBestandAPI.GetAll)
	corsWrapped := withCORS(mux)

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", corsWrapped)
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
