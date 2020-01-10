package Store

type NewTweet struct {
	Id        int    `json:"id"`
	Message   string `json:"message"`
	AccountId int    `json:"account_id"`
}
