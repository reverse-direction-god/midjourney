package service

import (
	"context"
	"crypto/rsa"
	"log"
	"mj/model"
	"strconv"
	"time"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

type wechat struct{}

var Wechat wechat

var (
	mchPrivateKey              *rsa.PrivateKey
	appId                      string = "*"
	mchID                      string = "*" // 商户号
	mchCertificateSerialNumber string = "*" // 商户证书序列号
	mchAPIv3Key                string = "*" // 商户APIv3密钥
	requestUrl                        = "https://api.mch.weixin.qq.com/v3/pay/transactions/native"
	method                            = "POST"
)

func init() {
	mchPrivateKey, err = utils.LoadPrivateKeyWithPath("./wechatpay/apiclient_key.pem")
	if err != nil {
		log.Print("load merchant private key error")
	}

}

func preferential(userId string) bool {
	var mod []model.Coupons
	DB.Model(&model.Coupons{}).Where("user_id=?", userId).Find(&mod)
	if len(mod) == 0 {
		return false
	}
	DB.Model(&model.Coupons{}).Where("id=?", mod[0].ID).Delete(nil)
	return true
}
func (s *wechat) Pay(userId string, Description string, number int64) (string, string) {

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
		return "", ""
	}
	//优惠
	ok := preferential(userId)
	var GoodsTag string
	if ok {
		switch number {
		case 2990:
			number = number - 1000
			GoodsTag = "优惠券立减10元"
		case 19900:
			number = number - 2000
			GoodsTag = "优惠券立减20元"
		case 39900:
			number = number - 3000
			GoodsTag = "优惠券立减30元"
		}
	}
	GoodsTag = "无优惠券"
	now := time.Now()
	nodString := strconv.Itoa(int(now.Unix()))
	svc := native.NativeApiService{Client: client}
	note := nodString + userId + "99085ylsx"
	resp, result, err := svc.Prepay(ctx,
		native.PrepayRequest{
			Appid:       core.String(appId),
			Mchid:       core.String(mchID),
			Description: core.String(Description),
			OutTradeNo:  core.String(note),
			TimeExpire:  core.Time(now),

			NotifyUrl: core.String("*"),
			// Attach:      core.String("充值测试数据说明"),
			GoodsTag: core.String(GoodsTag), //没有 非必填
			// LimitPay:      []string{"LimitPay_example"},
			SupportFapiao: core.Bool(false),
			Amount: &native.Amount{
				Currency: core.String("CNY"),
				Total:    core.Int64(number),
			},
		},
	)

	if err != nil {
		// 处理错误
		log.Printf("call Prepay err:%s", err)
		return "", ""
	} else {
		// 处理返回结果
		log.Printf("status=%d resp=%s", result.Response.StatusCode, resp)
		return *resp.CodeUrl, note
	}

}
