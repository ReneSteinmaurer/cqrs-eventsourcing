create table cart_items(
    cart_id integer not null,
    item text not null,
    quantity integer not null default 0,
    primary key (cart_id, item)
)