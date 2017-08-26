package index

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/zwh8800/md-blog-gen/conf"
	"github.com/zwh8800/md-blog-gen/model"
	"github.com/zwh8800/md-blog-gen/render"
	"github.com/zwh8800/md-blog-gen/service"
	"github.com/zwh8800/md-blog-gen/util"
)

func AlipayIndex(c *gin.Context) {
	c.Render(http.StatusOK, render.NewRender("alipay.html", gin.H{
		"site":   conf.Conf.Site,
		"social": conf.Conf.Social,
		"prod":   conf.Conf.Env.Prod,
		"haha":   util.HahaGenarate(),
	}))
}

func AlipayDo(c *gin.Context) {
	price, err := strconv.ParseInt(c.PostForm("price"), 10, 64)
	if err != nil {
		glog.Error(err)
		c.Redirect(http.StatusFound, "/alipay")
	}

	output, err := service.CreateOrder(&service.CreateOrderInput{
		Price: model.Cent(price),
	})
	if err != nil {
		glog.Error(err)
		return
	}

	qrcodeDataUrl, err := util.GenerateQrcodePngDataUrl(output.Url)
	if err != nil {
		glog.Error(err)
		c.Redirect(http.StatusFound, "/alipay")
		return
	}

	cookie, err := util.AesEncrypt(output.OrderId, []byte(conf.Conf.Crypto.AesKey))
	if err != nil {
		glog.Error(err)
		c.Redirect(http.StatusFound, "/alipay")
		return
	}
	c.SetCookie("tradeNo", cookie, 0, "", "", conf.Conf.Env.Prod, true)

	c.Render(http.StatusOK, render.NewRender("alipay_do.html", gin.H{
		"orderId":       output.OrderId,
		"qrcodeDataUrl": template.URL(qrcodeDataUrl),
		"site":          conf.Conf.Site,
		"social":        conf.Conf.Social,
		"prod":          conf.Conf.Env.Prod,
		"haha":          util.HahaGenarate(),
	}))
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		EnableCompression: true,
	}
)

func AlipayWs(c *gin.Context) {
	cookie, err := c.Cookie("tradeNo")
	if err != nil {
		return
	}

	orderId, err := util.AesDecrypt(cookie, []byte(conf.Conf.Crypto.AesKey))
	if err != nil {
		glog.Error(err)
		return
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	defer ws.Close()

	for {
		output := service.WaitOrderResult(orderId)
		glog.Infoln("websocket output:", output)
		if output == nil {
			return
		}
		if err := ws.WriteJSON(output); err != nil {
			glog.Error(err)
			return
		}
	}
}

func AlipayQuery(c *gin.Context) {
	cookie, err := c.Cookie("tradeNo")
	if err != nil {
		return
	}

	orderId, err := util.AesDecrypt(cookie, []byte(conf.Conf.Crypto.AesKey))
	if err != nil {
		glog.Error(err)
		return
	}

	output, err := service.GetOrderResult(orderId)
	if err != nil {
		glog.Error(err)
		c.JSON(http.StatusOK, nil)
		return
	}
	c.JSON(http.StatusOK, output)
}

func AlipayNotify(c *gin.Context) {
	service.HandleAlipayNotification(c.Request)
	c.String(http.StatusOK, "OK")
}

func AlipayStatus(c *gin.Context) {
	c.SetCookie("tradeNo", "", -1, "", "", false, true)
	cookie, err := c.Cookie("tradeNo")
	if err != nil {
		c.Redirect(http.StatusFound, "/alipay")
		return
	}

	orderId, err := util.AesDecrypt(cookie, []byte(conf.Conf.Crypto.AesKey))
	if err != nil {
		glog.Error(err)
		c.Redirect(http.StatusFound, "/alipay")
		return
	}

	output, err := service.GetOrderResult(orderId)
	if err != nil {
		glog.Error(err)
		c.Redirect(http.StatusFound, "/alipay")
		return
	}

	c.Render(http.StatusOK, render.NewRender("alipay_status.html", gin.H{
		"status": output,
		"site":   conf.Conf.Site,
		"social": conf.Conf.Social,
		"prod":   conf.Conf.Env.Prod,
		"haha":   util.HahaGenarate(),
	}))
}

func AlipayRefund(c *gin.Context) {
	//outTradeNo := c.PostForm("no")
	//client := alipay.New(conf.Conf.Alipay.AppId,
	//	conf.Conf.Alipay.PartnerId,
	//	[]byte(conf.Conf.Alipay.PublicKey),
	//	[]byte(conf.Conf.Alipay.PrivateKey),
	//	conf.Conf.Alipay.Prod)
	//
	//request := alipay.AliPayTradeCancel{
	//	OutTradeNo: outTradeNo,
	//}
	//
	//glog.Infoln("request:", util.JsonStringify(request, true))
	//response, err := client.TradeCancel(request)
	//if err != nil {
	//	glog.Error(err)
	//	return
	//}
	//glog.Infoln("response:", util.JsonStringify(response, true))
}
