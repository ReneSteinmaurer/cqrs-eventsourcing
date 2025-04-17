package add_item

type AddItemToCartCommand struct {
	CartId   int
	Item     string
	Quantity int
}
