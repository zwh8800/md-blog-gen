package model

import (
	"encoding/json"
	"time"

	"github.com/gocraft/dbr"
)

const ToushiTableName = "toushi"

type Toushi struct {
	Id                   int64        `db:"id" json:"id"`
	Status               ToushiStatus `db:"status" json:"status"`
	UUID                 string       `db:"uuid" json:"uuid"`
	PriceInCent          Cent         `db:"price_in_cent" json:"priceInCent"`
	AlipayQRCode         string       `db:"alipay_qrcode" json:"alipayQRCode"`
	AlipayTradeNo        string       `db:"alipay_trade_no" json:"alipayTradeNo"`
	AlipayBuyerLogonId   string       `db:"alipay_buyer_logon_id" json:"alipayBuyerLogonId"`
	AlipayTradeStatus    string       `db:"alipay_trade_status" json:"alipayTradeStatus"`
	AlipayTotalAmount    Cent         `db:"alipay_total_amount" json:"alipayTotalAmount"`
	AlipayReceiptAmount  Cent         `db:"alipay_receipt_amount" json:"alipayReceiptAmount"`
	AlipayBuyerPayAmount Cent         `db:"alipay_buyer_pay_amount" json:"alipayBuyerPayAmount"`
	AlipayPointAmount    Cent         `db:"alipay_point_amount" json:"alipayPointAmount"`
	AlipayInvoiceAmount  Cent         `db:"alipay_invoice_amount" json:"alipayInvoiceAmount"`
	AlipaySendPayDate    dbr.NullTime `db:"alipay_send_pay_date" json:"alipaySendPayDate"`
	AlipayFundBillList   string       `db:"alipay_fund_bill_list" json:"alipayFundBillList"`
	AlipayBuyerUserId    string       `db:"alipay_buyer_user_id" json:"alipayBuyerUserId"`
	CreatedAt            time.Time    `db:"created_at" json:"createdAt"`
	AlipayCreatedAt      dbr.NullTime `db:"alipay_created_at" json:"alipayCreatedAt"`
	FinishAt             dbr.NullTime `db:"finish_at" json:"finishAt"`
}

type ToushiStatus int

const (
	ToushiStatusPending ToushiStatus = iota
	ToushiStatusOrderCreated
	ToushiStatusSuccess
	ToushiStatusFailed
)

func (t ToushiStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t ToushiStatus) String() string {
	switch t {
	case ToushiStatusPending:
		return "PENDING"
	case ToushiStatusOrderCreated:
		return "ORDER_CREATED"
	case ToushiStatusSuccess:
		return "SUCCESS"
	case ToushiStatusFailed:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

func (t ToushiStatus) IsPending() bool {
	return t == ToushiStatusPending
}

func (t ToushiStatus) IsOrderCreated() bool {
	return t == ToushiStatusOrderCreated
}

func (t ToushiStatus) IsSuccess() bool {
	return t == ToushiStatusSuccess
}

func (t ToushiStatus) IsFailed() bool {
	return t == ToushiStatusFailed
}

func (t ToushiStatus) IsFinish() bool {
	switch t {
	case ToushiStatusPending:
		fallthrough
	case ToushiStatusOrderCreated:
		return false
	case ToushiStatusSuccess:
		fallthrough
	case ToushiStatusFailed:
		fallthrough
	default:
		return true
	}
}
