package domain

type Product struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	Code_value   string  `json:"code_value"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration"`
	Price        float64 `json:"price"`
}

type Request struct {
	Name         string  `json:"name" validate:"required"`
	Quantity     int     `json:"quantity" validate:"required"`
	Code_value   string  `json:"code_value" validate:"required"`
	Is_published bool    `json:"is_published"`
	Expiration   string  `json:"expiration" validate:"required"`
	Price        float64 `json:"price" validate:"required"`
}

type PatchRequest struct {
	Name         *string  `json:"name"`
	Quantity     *int     `json:"quantity"`
	Code_value   *string  `json:"code_value"`
	Is_published *bool    `json:"is_published"`
	Expiration   *string  `json:"expiration"`
	Price        *float64 `json:"price"`
}
