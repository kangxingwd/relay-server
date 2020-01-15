package handler

import (
	"crypto/md5"
	"fmt"
	logger "github.com/cihub/seelog"
	"github.com/valyala/fasthttp"
	"gitlab.com/TenbayMCloud/awesome-raserver/common"
	"io"
)

func LoginHandlerFunc(ctx *fasthttp.RequestCtx) {
	var ret common.RetInfo

	// 参数校验
	if 	(false == ctx.QueryArgs().Has("dev_id")) ||
		(false == ctx.QueryArgs().Has("purchaser")) ||
		(false == ctx.QueryArgs().Has("passwd")) {
		ret.Result = common.RetUrlErr
		common.RetJson(ctx, &ret)
		logger.Debug(common.ErrInfoMap[ret.Result], ctx.QueryArgs().String())
		return
	}

	dev_id := ctx.QueryArgs().Peek("dev_id")
	purchaser := ctx.QueryArgs().Peek("purchaser")
	passwd :=  ctx.QueryArgs().Peek("passwd")

	if true == checkPasswd(string(dev_id), string(purchaser),string(passwd)) {
		ret.Result = common.RetSuccess
		ret.Data.Token = makeToken(string(dev_id), string(passwd))
		saveToken(ret.Data.Token)
	}else {
		ret.Result = common.RetPasswdAuthFailed
	}

	common.RetJson(ctx, &ret)
}

// token auth
func TokenAuth() (ret bool, err error) {

	return true,nil
}

// passwd chaeck
func checkPasswd(dev_id string, purchaser string, passwd string) bool {
	password := myMd5(dev_id, purchaser, "999")
	logger.Debugf("dev_id(%s) success password = %s\n", dev_id, password)
	return passwd == password
}

// make token by md5
func makeToken(dev_id string, passwd string) string {
	return myMd5(dev_id, passwd,  "999")
}

// make md5
func myMd5(strs ...string) string {
	h := md5.New()
	for _,v := range strs{
		io.WriteString(h, v)
	}
	ret := fmt.Sprintf("%x", h.Sum(nil))
	return ret
}

// save token to redis
func saveToken(token string) (errr error) {

	return nil
}

