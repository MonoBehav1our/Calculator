package model

type ExpressionResult struct {
	Id         int    `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}
