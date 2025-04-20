package remove_item

type ItemRemovedFromCartEvent struct {
	CartId int
	Item   string
}
