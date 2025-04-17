package add_item

type AddItemToWishlistCommandV1 struct {
	WishlistId int
	Item       string
}

type AddItemToWishlistCommandV2 struct {
	WishlistId int
	Item       string
	UserId     string
}
