package images

type Image struct {
	ID        *int64  `json:"id"`
	ProductId *int64  `json:"product_id"`
	Path      *string `json:"path"`
}
