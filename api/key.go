package api

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"github.com/provider-go/pkg/encryption"
	"github.com/provider-go/pkg/logger"
	"github.com/provider-go/pkg/output"
	"github.com/provider-go/pkg/util"
	"github.com/provider-go/sso/global"
	"github.com/provider-go/sso/middleware"
	"github.com/provider-go/sso/models"
	"strconv"
)

func LoginByKey(ctx *gin.Context) {
	var returnDID string
	json := make(map[string]interface{})
	_ = ctx.BindJSON(&json)
	pubkey := output.ParamToString(json["pubkey"])
	timestamp := output.ParamToString(json["timestamp"])
	sign := output.ParamToString(json["sign"])
	//判断时间戳是与系统时间相差
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		output.ReturnErrorResponse(ctx, 9999, "时间戳不正确~")
		return
	}
	differ := util.CurrentSecond() - ts
	// 验证签名
	ok := encryption.SM2Verify(pubkey, timestamp, sign)
	if !ok && differ >= 0 && differ < 120 {
		output.ReturnErrorResponse(ctx, 9999, "用户签名不正确~")
		return
	}
	// 对比数据库记录,公钥是否存在
	item, err := models.ViewSSOKey(pubkey)
	if err != nil {
		if err.Error() == "ErrRecordNotFound" {
			returnDID = hex.EncodeToString([]byte(pubkey))[0:10]
			err = models.CreateSSOKey(returnDID, pubkey)
			if err != nil {
				logger.Error("LoginByKey", "step", "CreateSSOKey", "err", err)
				output.ReturnErrorResponse(ctx, 9999, "系统错误~")
				return
			}
		} else {
			logger.Error("LoginByKey", "step", "ViewSSOKey", "err", err)
			output.ReturnErrorResponse(ctx, 9999, "系统错误~")
			return
		}
	} else {
		returnDID = item.DID
	}

	// 生成token
	token := middleware.InitJwt(global.SecretKey).GenerateToken(returnDID)
	output.ReturnSuccessResponse(ctx, token)

}
