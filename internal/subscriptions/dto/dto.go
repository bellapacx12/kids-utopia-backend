package dto

type CreateSubscriptionRequest struct {
	Plan string `json:"plan" binding:"required"`
}

type SubscriptionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}