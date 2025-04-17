package add_item

type ItemAddedToWishlistEventV1 struct {
	WishlistId int
	Item       string
}

type ItemAddedToWishlistEventV2 struct {
	WishlistId int
	Item       string
	UserId     string
}

const ItemAddedToWishlistEventTypeV1 = "ItemAddedToWishlistV1"
const ItemAddedToWishlistEventTypeV2 = "ItemAddedToWishlistV2"
