package handler

import (
	"fmt"

	logger "github.com/cihub/seelog"
	"github.com/valyala/fasthttp"
	"gitlab.com/TenbayMCloud/awesome-raserver/common"
	"gitlab.com/TenbayMCloud/awesome-raserver/model"
)

// AdInfoGet 根据设备信息获取广告信息
// curl "http://127.0.0.1:8008/apiv1/user/ad/ad_info_get?token=1123" -d 'isshe=chudai&test2=test2&data={}' -v
// curl "http://127.0.0.1:8008/apiv1/user/ad/ad_info_get?token=1123" -d 'isshe=chudai&test2=test2&data=' -v
// curl "http://127.0.0.1:8008/apiv1/user/ad/ad_info_get?token=1123"
//		-d 'isshe=chudai&test2=test2&data={"dev_vendor":"ORICO","dev_product":"NAS","dev_type":"WRC10"}' -v
func AdInfoGet(ctx *fasthttp.RequestCtx) error {
	var gArgs = []string{} //[]string{"token"}
	var pArgs = []string{"data"}

	var Args = common.ArgsGet(ctx, gArgs, pArgs)
	var data = Args["data"]
	if len(data) <= 0 {
		common.ResponseErrJSON(ctx, "Invalid parameters")
		return fmt.Errorf("Invalid parameters: %+v", Args)
	}

	var devInfo common.DevInfo
	if err := common.JSONDecode(data, &devInfo); err != nil {
		logger.Infof("JSONDecode err: %v\n", err)
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	if len(devInfo.DevVendor) <= 0 ||
		len(devInfo.DevProduct) <= 0 ||
		len(devInfo.DevType) <= 0 {
		common.ResponseErrJSON(ctx, "Invalid device info")
		return fmt.Errorf("Invalid device info: %+v", devInfo)
	}

	var adInfo model.AdInfo
	adInfo.DevVendor = devInfo.DevVendor
	adInfo.DevProduct = devInfo.DevProduct
	adInfo.DevType = devInfo.DevType

	pAdInfo, err := adInfo.Get()
	if err != nil {
		logger.Infof("AdInfoAdd adInfo.Add err: %v\n", err)
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	common.ResponseOkJSON(ctx, *pAdInfo)
	return nil
}

// AdInfoGetFromArgs 获取参数并转换为AdInfo结构
// {"data":{...}} => AdInfo{...}
func AdInfoGetFromArgs(ctx *fasthttp.RequestCtx) (*model.AdInfo, error) {
	var gArgs = []string{} //[]string{"token"}
	var pArgs = []string{"data"}
	var Args = common.ArgsGet(ctx, gArgs, pArgs)

	var data = Args["data"]
	if len(data) <= 0 {
		common.ResponseErrJSON(ctx, "Invalid parameters")
		return nil, fmt.Errorf("Invalid parameters: %+v", Args)
	}

	// data里面应该是一个Json，对应AdInfo结构
	// 进行JSON解析
	var adInfo model.AdInfo
	if err := common.JSONDecode(data, &adInfo); err != nil {
		logger.Infof("AdInfoAdd JSONDecode err: %v\n", err)
		common.ResponseErrJSON(ctx, err.Error())
		return nil, err
	}
	logger.Infof("adInfo: %+v\n", adInfo)
	logger.Infof("adInfo.DevVendor: %+v\n", adInfo.DevVendor)
	return &adInfo, nil
}

// AdInfoAdd 根据AdInfo进行操作
// curl "http://127.0.0.1:8008/apiv1/user/ad/ad_info_add?token=1123"
//		-d 'isshe=chudai&test2=test2&data={"dev_vendor":"ORICO",
//		"dev_product":"NAS","dev_type":"WRC20","ad_url":"http://host:ip/haha",
//		"dev_softver":"v1.0-20190226103540","desc":"wrc10-lalalalal"}' -v
func AdInfoAdd(ctx *fasthttp.RequestCtx) error {
	pAdInfo, err := AdInfoGetFromArgs(ctx)
	if err != nil {
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	if err := pAdInfo.AddNoDup(); err != nil {
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	common.ResponseOkJSON(ctx, *pAdInfo)
	return nil
}

// AdInfoUpdate 根据AdInfo进行操作
// curl "http://127.0.0.1:8008/apiv1/user/ad/ad_info_update?token=1123"
// 		-d 'isshe=chudai&test2=test2&data={"dev_vendor":"ORICO",
//		"dev_product":"NAS","dev_type":"WRC20","ad_url":"http://host:ip/haha",
// 		"dev_softver":"v2.0-20190226103540","desc":"wrc20-lalalalal","id":2}' -v
// id是需要的，如果没有id，就会变成增加
func AdInfoUpdate(ctx *fasthttp.RequestCtx) error {
	pAdInfo, err := AdInfoGetFromArgs(ctx)
	if err != nil {
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	var tInfo model.AdInfo
	tInfo.ID = pAdInfo.ID
	ptInfo, err := tInfo.Get()
	if ptInfo != nil && err == nil {
		pAdInfo.CreateTime = ptInfo.CreateTime
		pAdInfo.LastUpdateTime = ptInfo.LastUpdateTime
		pAdInfo.LastAccessTime = ptInfo.LastAccessTime
	}

	if err := pAdInfo.Update(); err != nil {
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	common.ResponseOkJSON(ctx, *pAdInfo)
	return nil
}

// AdInfoDel 根据AdInfo进行操作
// curl "http://127.0.0.1:8008/apiv1/user/ad/ad_info_del?token=1123" -d 'data={"id":6}'
func AdInfoDel(ctx *fasthttp.RequestCtx) error {
	pAdInfo, err := AdInfoGetFromArgs(ctx)
	if err != nil {
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	if err := pAdInfo.Del(); err != nil {
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	common.ResponseOkJSON(ctx, true)
	return nil
}

// AdInfoGetAll 根据AdInfo进行操作
// curl "http://127.0.0.1:8008/apiv1/user/ad/ad_info_gets?token=1123"
func AdInfoGetAll(ctx *fasthttp.RequestCtx) error {
	Infos, err := model.GetAllAdInfo()
	if err != nil {
		logger.Infof("AdInfoAdd adInfo.Add err: %v\n", err)
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	common.ResponseOkJSON(ctx, *Infos)
	return nil
}
