package events

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

type ItemRemovedFromWishlistV1 struct {
	WishlistId int
	Item       string
	UserId     string
}

const ItemRemovedFromWishlistTypeV1 = "ItemRemovedFromWishlistV1"
