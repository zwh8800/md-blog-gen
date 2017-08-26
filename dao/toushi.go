package dao

import (
	"github.com/gocraft/dbr"
	"github.com/zwh8800/md-blog-gen/model"
)

var insertToushiOmitTable = []string{"id"}

func InsertToushi(sr dbr.SessionRunner, toushi *model.Toushi) (*model.Toushi, error) {
	obj, err := commonInsert(sr, model.ToushiTableName, toushi, insertToushiOmitTable)
	if err != nil {
		return nil, err
	}
	return obj.(*model.Toushi), nil
}

func UpdateToushiQRCode(sr dbr.SessionRunner, toushiId int64, qrcode string) error {
	_, err := sr.Update(model.ToushiTableName).Set("alipay_qrcode", qrcode).Where("id = ?", toushiId).Exec()
	return err
}

func FindToushiForUpdate(sr dbr.SessionRunner, orderId string) (*model.Toushi, error) {
	var toushi model.Toushi
	err := sr.SelectBySql(`SELECT * FROM `+model.ToushiTableName+` WHERE uuid = ? FOR UPDATE`, orderId).LoadStruct(&toushi)
	return &toushi, err
}

func FindToushi(sr dbr.SessionRunner, orderId string) (*model.Toushi, error) {
	var toushi model.Toushi
	err := sr.Select("*").From(model.ToushiTableName).Where("uuid = ?", orderId).LoadStruct(&toushi)
	return &toushi, err
}

func UpdateToushi(sr dbr.SessionRunner, toushi *model.Toushi) error {
	builder := sr.Update(model.ToushiTableName)
	if toushi.AlipayCreatedAt.Valid {
		builder.Set("alipay_created_at", toushi.AlipayCreatedAt)
	}
	_, err := builder.Set("status", toushi.Status).
		Set("alipay_trade_no", toushi.AlipayTradeNo).
		Set("alipay_buyer_logon_id", toushi.AlipayBuyerLogonId).
		Set("alipay_trade_status", toushi.AlipayTradeStatus).
		Set("alipay_total_amount", toushi.AlipayTotalAmount).
		Set("alipay_receipt_amount", toushi.AlipayReceiptAmount).
		Set("alipay_buyer_pay_amount", toushi.AlipayBuyerPayAmount).
		Set("alipay_point_amount", toushi.AlipayPointAmount).
		Set("alipay_invoice_amount", toushi.AlipayInvoiceAmount).
		Set("alipay_send_pay_date", toushi.AlipaySendPayDate).
		Set("alipay_fund_bill_list", toushi.AlipayFundBillList).
		Set("alipay_buyer_user_id", toushi.AlipayBuyerUserId).
		Set("finish_at", toushi.FinishAt).
		Where("id = ?", toushi.Id).
		Exec()
	return err
}
