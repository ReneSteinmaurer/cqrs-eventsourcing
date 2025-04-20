package add_item

type ItemAddedToCartEvent struct {
	CartId   int
	Item     string
	Quantity int
}
