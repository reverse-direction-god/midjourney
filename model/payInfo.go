package model

import "gorm.io/gorm"

type PayInfo struct {
	gorm.Model
	PayerTotal    int    `json:"payerTotal"`
	OutTradeNo    string `json:"out_trade_no"`
	UserId        int    `json:"userId"`
	TradeState    string `json:"trade_state"`
	TransactionId string `json:"transaction_id"`
}
