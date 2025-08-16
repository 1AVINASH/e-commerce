package products

type GetProductsInputDto struct {
	Id int64 `json:"id"`
}

type CreateProductsInputDto struct {
	Title         *string   `json:"title"`
	Description   *string   `json:"description"`
	OriginalPrice *string   `json:"original_price"`
	Discount      *string   `json:"discount"`
	Thumbnail     *string   `json:"thumbnail"`
	Photos        []*string `json:"photos"`
}

type UpdateProductsInputDto struct {
	ID            *int64    `json:"id"`
	Title         *string   `json:"title"`
	Description   *string   `json:"description"`
	OriginalPrice *string   `json:"original_price"`
	Discount      *string   `json:"discount"`
	Thumbnail     *string   `json:"thumbnail"`
	Photos        []*string `json:"photos"`
}
