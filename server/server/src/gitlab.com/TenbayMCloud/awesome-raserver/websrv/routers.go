package websrv

import (
	logger "github.com/cihub/seelog"
	"github.com/valyala/fasthttp"
	"gitlab.com/TenbayMCloud/awesome-raserver/common"
	"gitlab.com/TenbayMCloud/awesome-raserver/websrv/handler"
)

func login(ctx *fasthttp.RequestCtx) {
	handler.LoginHandlerFunc(ctx)
}

func relay(ctx *fasthttp.RequestCtx) {
	var ret common.RetInfo
	var err error

	logger.Infof("ctx.Request.Header\n %s", ctx.Request.Header.String())
	logger.Infof("ctx.PostBody:\n %s\n\n", string(ctx.PostBody()))

	// token 认证
	retsult, err := handler.TokenAuth()
	if err != nil {
		ret.Result = common.RetServerErr //	服务器内部错误
		common.RetJson(ctx, &ret)
		logger.Errorf("TokenAuth error: %s\n", err.Error())
		return
	}
	if retsult == false {
		ret.Result = common.RetTokenAuthFailed // token 认证失败
		common.RetJson(ctx, &ret)
		return
	}

	// 解析Json
	var msg common.CliMsg
	if err := common.JsonUnmarshal(ctx, &msg); err != nil {
		ret.Result = common.RetJsonFormatErr
		common.RetJson(ctx, &ret)
		return
	}
	logger.Infof("POST JSON is:\n %+v\n", msg)

	switch ctx.UserValue("option") {
	case common.OpNetJoin:
		handler.RelayNetJoin(ctx, &msg)
	case common.OpNetStop:
		handler.RelayNetStop(ctx, &msg)
	case common.OpRequestDomain:
		handler.RequestDomain(ctx, &msg)
	case common.OpAddMcloud:
		handler.AddMcloud(ctx, &msg)
	case common.OpHeartBeat:
		handler.RelayHeartBeat(ctx, &msg)
	default:
		ret.Result = common.RetRequestOptionErr
		common.RetJson(ctx, &ret)
	}
}

func mcloud(ctx *fasthttp.RequestCtx) {
	var msg common.CliMsg
	var ret common.RetInfo
	var err error

	logger.Infof("ctx.Request.Header\n %s", ctx.Request.Header.String())
	logger.Infof("ctx.PostBody:\n %s\n\n", string(ctx.PostBody()))

	// token 认证
	retsult, err := handler.TokenAuth()
	if err != nil {
		ret.Result = common.RetServerErr //	服务器内部错误
		common.RetJson(ctx, &ret)
		logger.Errorf("TokenAuth error: %s\n", err.Error())
		return
	}
	if retsult == false {
		ret.Result = common.RetTokenAuthFailed // token 认证失败
		common.RetJson(ctx, &ret)
		return
	}

	// 解析Json
	if err := common.JsonUnmarshal(ctx, &msg); err != nil {
		ret.Result = common.RetJsonFormatErr
		common.RetJson(ctx, &ret)
		return
	}
	logger.Infof("POST JSON is:\n %+v\n", msg)

	switch ctx.UserValue("option") {
	case common.OpNetJoin:
		handler.McloudNetJoin(ctx, &msg)
	case common.OpNetStop:
		handler.McloudNetStop(ctx, &msg)
	case common.OpTunnelState:
		handler.McloudTunnelState(ctx, &msg)
	case common.OpHeartBeat:
		handler.McloudHeartBeat(ctx, &msg)
	default:
		ret.Result = common.RetRequestOptionErr
		common.RetJson(ctx, &ret)
	}
}

func ad(ctx *fasthttp.RequestCtx) {
	var ret common.RetInfo

	switch ctx.UserValue("option") {
	case "ad_info_get":
		handler.AdInfoGet(ctx)
	case "ad_info_add":
		handler.AdInfoAdd(ctx)
	case "ad_info_update":
		handler.AdInfoUpdate(ctx)
	case "ad_info_del":
		handler.AdInfoDel(ctx)
	case "ad_info_gets":
		handler.AdInfoGetAll(ctx)
	case "ad_rsrc_info_update":
		handler.AdRsrcInfoUpdate(ctx)
	case "ad_rsrc_info_set":
		handler.AdRsrcInfoSet(ctx)
	case "ad_rsrc_info_get":
		handler.AdRsrcInfoGet(ctx)
	case "ad_rsrc_info_del":
		handler.AdRsrcInfoDel(ctx)
	default:
		ret.Result = common.RetRequestOptionErr
		common.RetJson(ctx, &ret)
	}
}
