package dtos

type GetBalanceRequest struct {
	Address string `json:"address" validate:"uuid4,required"`
}

type GetBalanceResponse struct {
	Balance string `json:"balance"`
}
