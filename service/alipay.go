package service

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/smartwalle/alipay"
	"github.com/zwh8800/md-blog-gen/dao"
	"github.com/zwh8800/md-blog-gen/model"
	"github.com/zwh8800/md-blog-gen/util"
	"github.com/zwh8800/md-blog-gen/util/generator"
)

type OrderResult struct {
	OrderId string             `json:"orderId"`
	Status  model.ToushiStatus `json:"status"`
}

var (
	orderNotifyMapLock sync.Mutex
	orderNotifyMap     = make(map[string][]chan *OrderResult)
)

func AddNotifyChan(orderId string, ch chan *OrderResult) {
	orderNotifyMapLock.Lock()
	defer orderNotifyMapLock.Unlock()
	list, ok := orderNotifyMap[orderId]
	if !ok {
		list = make([]chan *OrderResult, 0)
	}
	list = append(list, ch)
	orderNotifyMap[orderId] = list
}

func DeleteNotifyChan(orderId string, ch chan *OrderResult) {
	orderNotifyMapLock.Lock()
	defer orderNotifyMapLock.Unlock()
	list, ok := orderNotifyMap[orderId]
	if !ok {
		return
	}
	for i, val := range list {
		if val == ch {
			list = append(list[:i], list[i+1:]...)
		}
	}
	orderNotifyMap[orderId] = list
}

func NotifyAllChan(orderId string, result *OrderResult) {
	orderNotifyMapLock.Lock()
	defer orderNotifyMapLock.Unlock()
	list := orderNotifyMap[orderId]
	glog.Infoln("NotifyAllChan", orderId, result, len(list))

	for _, ch := range list {
		ch <- result
		if result.Status.IsFinish() {
			close(ch)
		}
	}
	if result.Status.IsFinish() {
		delete(orderNotifyMap, orderId)
	}
}

type CreateOrderInput struct {
	Price model.Cent
}

type CreateOrderOutput struct {
	OrderId string
	Url     string
}

func CreateOrder(input *CreateOrderInput) (*CreateOrderOutput, error) {
	orderId := strconv.FormatUint(generator.Default.Next(), 10)

	sess := newSession()
	tx, err := sess.Begin()
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	defer tx.RollbackUnlessCommitted()
	toushi, err := dao.InsertToushi(tx, &model.Toushi{
		Status:      model.ToushiStatusPending,
		UUID:        orderId,
		PriceInCent: input.Price,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	client := newAlipayClient()
	request := alipay.AliPayTradePreCreate{
		OutTradeNo:     orderId,
		Subject:        "向一只废喵投食",
		StoreId:        "废喵一号店",
		TotalAmount:    model.Cent(input.Price).CurrencyWithoutComma(),
		TimeoutExpress: "90m",
		NotifyURL:      "https://lengzzz.com/alipay/notify",
	}

	glog.Infoln("request:", util.JsonStringify(request, true))
	response, err := client.TradePreCreate(request)
	if err != nil {
		glog.Error(err)
		return nil, fmt.Errorf("Code: %s, Msg: %s, SubCode: %s, SubMsg: %s",
			response.AliPayPreCreateResponse.Code, response.AliPayPreCreateResponse.Msg,
			response.AliPayPreCreateResponse.SubCode, response.AliPayPreCreateResponse.SubMsg)
	}
	glog.Infoln("response:", util.JsonStringify(response, true))

	if err := dao.UpdateToushiQRCode(tx, toushi.Id, response.AliPayPreCreateResponse.QRCode); err != nil {
		glog.Error(err)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		glog.Error(err)
		return nil, err
	}

	go PollAlipay(orderId)

	return &CreateOrderOutput{
		OrderId: orderId,
		Url:     response.AliPayPreCreateResponse.QRCode,
	}, nil
}

func PollAlipay(orderId string) {
	for i := 0; i < 180; i++ {
		time.Sleep(5 * time.Second)

		client := newAlipayClient()
		response, err := client.TradeQuery(alipay.AliPayTradeQuery{
			OutTradeNo: orderId,
		})
		if err != nil {
			glog.Error(err)
			continue
		}
		glog.Infoln("response:", util.JsonStringify(response, true))
		if response.IsSuccess() {
			if HandleAlipayResult(&AlipayResultInput{
				TradeNo:        response.AliPayTradeQuery.TradeNo,
				OutTradeNo:     response.AliPayTradeQuery.OutTradeNo,
				BuyerLogonId:   response.AliPayTradeQuery.BuyerLogonId,
				TradeStatus:    response.AliPayTradeQuery.TradeStatus,
				TotalAmount:    response.AliPayTradeQuery.TotalAmount,
				ReceiptAmount:  response.AliPayTradeQuery.ReceiptAmount,
				BuyerPayAmount: response.AliPayTradeQuery.BuyerPayAmount,
				PointAmount:    response.AliPayTradeQuery.PointAmount,
				InvoiceAmount:  response.AliPayTradeQuery.InvoiceAmount,
				SendPayDate:    response.AliPayTradeQuery.SendPayDate,
				FundBillList:   util.JsonStringify(response.AliPayTradeQuery.FundBillList, false),
				BuyerUserId:    response.AliPayTradeQuery.BuyerUserId,
			}) {
				return
			}
		}
	}
}

func HandleAlipayNotification(req *http.Request) {
	client := newAlipayClient()
	notification, err := client.GetTradeNotification(req)
	if err != nil {
		glog.Error(err)
		return
	}
	totalAmount, _ := strconv.ParseFloat(notification.TotalAmount, 64)
	receiptAmount, _ := strconv.ParseFloat(notification.ReceiptAmount, 64)
	buyerPayAmount, _ := strconv.ParseFloat(notification.BuyerPayAmount, 64)
	pointAmount, _ := strconv.ParseFloat(notification.PointAmount, 64)
	invoiceAmount, _ := strconv.ParseFloat(notification.InvoiceAmount, 64)

	HandleAlipayResult(&AlipayResultInput{
		TradeNo:        notification.TradeNo,
		OutTradeNo:     notification.OutTradeNo,
		BuyerLogonId:   notification.BuyerLogonId,
		TradeStatus:    notification.TradeStatus,
		TotalAmount:    totalAmount,
		ReceiptAmount:  receiptAmount,
		BuyerPayAmount: buyerPayAmount,
		PointAmount:    pointAmount,
		InvoiceAmount:  invoiceAmount,
		FundBillList:   util.JsonStringify(notification.FundBillList, false),
	})
}

type AlipayResultInput struct {
	TradeNo        string
	OutTradeNo     string
	BuyerLogonId   string
	TradeStatus    string
	TotalAmount    float64
	ReceiptAmount  float64
	BuyerPayAmount float64
	PointAmount    float64
	InvoiceAmount  float64
	SendPayDate    string
	FundBillList   string
	BuyerUserId    string
}

func HandleAlipayResult(input *AlipayResultInput) bool {
	targetTime := time.Now()
	sess := newSession()
	tx, err := sess.Begin()
	if err != nil {
		glog.Error(err)
		return false
	}
	defer tx.RollbackUnlessCommitted()
	order, err := dao.FindToushiForUpdate(tx, input.OutTradeNo)
	if err != nil {
		glog.Error(err)
		return false
	}

	order.AlipayTradeNo = input.TradeNo
	order.AlipayBuyerLogonId = input.BuyerUserId
	order.AlipayTradeStatus = input.TradeStatus
	order.AlipayTotalAmount = model.Cent(0).ParseFloat(input.TotalAmount)
	order.AlipayReceiptAmount = model.Cent(0).ParseFloat(input.ReceiptAmount)
	order.AlipayBuyerPayAmount = model.Cent(0).ParseFloat(input.BuyerPayAmount)
	order.AlipayPointAmount = model.Cent(0).ParseFloat(input.PointAmount)
	order.AlipayInvoiceAmount = model.Cent(0).ParseFloat(input.InvoiceAmount)
	if sendPayDate, err := time.Parse("2006-01-02 15:04:05", input.SendPayDate); err == nil {
		order.AlipaySendPayDate.Valid = true
		order.AlipaySendPayDate.Time = sendPayDate
	}
	order.AlipayFundBillList = input.FundBillList
	order.AlipayBuyerUserId = input.BuyerLogonId

	finished := false

	switch input.TradeStatus {
	case "WAIT_BUYER_PAY":
		order.Status = model.ToushiStatusOrderCreated
		order.AlipayCreatedAt.Valid = true
		order.AlipayCreatedAt.Time = targetTime

	case "TRADE_SUCCESS":
		fallthrough
	case "TRADE_FINISHED":
		order.Status = model.ToushiStatusSuccess
		order.FinishAt.Valid = true
		order.FinishAt.Time = targetTime
		finished = true

	case "TRADE_CLOSED":
		fallthrough
	default:
		order.Status = model.ToushiStatusFailed
		order.FinishAt.Valid = true
		order.FinishAt.Time = targetTime
		finished = true
	}

	if err := dao.UpdateToushi(tx, order); err != nil {
		glog.Error(err)
		return false
	}

	if err := tx.Commit(); err != nil {
		glog.Error(err)
		return false
	}

	NotifyAllChan(input.OutTradeNo, &OrderResult{
		OrderId: order.UUID,
		Status:  order.Status,
	})
	return finished
}

func WaitOrderResult(orderId string) *OrderResult {
	ch := make(chan *OrderResult)
	AddNotifyChan(orderId, ch)
	defer DeleteNotifyChan(orderId, ch)
	select {
	case data := <-ch:
		return data
	case <-time.After(15 * time.Minute):
		return nil
	}
}

func GetOrderResult(orderId string) (*OrderResult, error) {
	sess := newSession()
	order, err := dao.FindToushi(sess, orderId)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return &OrderResult{
		OrderId: order.UUID,
		Status:  order.Status,
	}, nil
}
