create table wishlist_items
(
    wishlist_id  integer not null,
    item     text    not null,
    primary key (wishlist_id, item)
)