package messages

type CreateMealRequest struct {
	Meal     string `json:"meal"`
	Calories int    `json:"calories"`
	Date     string `json:"date"`
	Time     string `json:"time"`
}
