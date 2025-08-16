package products

type Product struct {
	ID            *int64    `json:"id"`
	Title         *string   `json:"title"`
	Description   *string   `json:"description"`
	OriginalPrice *string   `json:"original_price"`
	Discount      *string   `json:"discount"`
	Thumbnail     *string   `json:"thumbnail"`
}
