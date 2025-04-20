package remove_item

type RemoveItemFromWishlistCommandV1 struct {
	WishlistId int
	Item       string
	UserId     string
}
