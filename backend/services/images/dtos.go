package images

type GetProductPhotosInputDto struct {
	ProductId *int64 `json:"product_id"`
}

type CreateProductPhotoInputDto struct {
	ProductId *int64 `json:"product_id"`
}
