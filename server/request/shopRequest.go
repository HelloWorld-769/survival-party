package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type BuyStoreRequest struct {
	ProductId string `json:"productId"`
	Token     string `json:"token"`
}

func (a BuyStoreRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ProductId, validation.Required),
	)
}
