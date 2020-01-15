package handler

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	logger "github.com/cihub/seelog"
	"github.com/valyala/fasthttp"
	"gitlab.com/TenbayMCloud/awesome-raserver/common"
	"gitlab.com/TenbayMCloud/awesome-raserver/model"
)

// TraverseArgs 遍历文件夹时，CallBack回调函数的参数（因为参数较多，因此放到一个Struct里面）
type TraverseArgs struct {
	Root           string
	PreDir         []string
	PreDirDepth    int
	MaxPreDirDepth int
	RelativePath   string
}

// AdRsrcPath 广告资源路径
// Path 这个路径和nginx对应 == nginx/ad
// PreDirDepth 前置文件夹深度, 如ad目录下有: Orico/Nas/wrc10/resource.png, Orico/Nas/wrc10就是前置路径，3层
// CallBack 回调函数，Path路径里面的广告信息生成方法
type AdRsrcPath struct {
	Path        string
	PreDirDepth int
	CallBack    func(ta *TraverseArgs) error
}

var rsrcPaths = []AdRsrcPath{
	{
		Path:        "/var/www/ad/resource", // /var/www/ad: 这里要修改！
		PreDirDepth: 3,
		CallBack:    AdRsrcInfoGenAndAddToDb,
	},
}

// DefaultAdURL: 默认的广告URL，生成广告信息是，使用这个URL
// AdRsrcURI: 广告资源的访问路径，完整路径如：
//	http://192.168.2.162:8080/ad/vendor/product/type/resource.png
const (
	DefaultAdURL = "http://aonas.tenbay.cn:8090/ad/"         // 这里要修改！
	AdRsrcURI    = "http://aonas.tenbay.cn:8090/ad/resource" // 这里要修改！
)

// InfoGetFromArgs 获取参数
func InfoGetFromArgs(ctx *fasthttp.RequestCtx, data []byte,
	info interface{}) error {

	if len(data) <= 0 {
		common.ResponseErrJSON(ctx, "Invalid parameters")
		return fmt.Errorf("Invalid parameters: %+v", data)
	}

	// data里面应该是一个Json，对应AdInfo结构
	// 进行JSON解析
	if err := common.JSONDecode(data, info); err != nil {
		logger.Infof("AdInfoAdd JSONDecode err: %v\n", err)
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}
	logger.Infof("info: %+v\n", info)
	return nil
}

func isLastPreDir(Depth int) bool {
	if Depth == 0 {
		return true
	}
	return false
}

func isPreDir(Depth int) bool {
	if Depth >= 0 {
		return true
	}
	return false
}

// TraversePath 遍历文件夹
// root string: 根路径
// preDir []string: 前置文件夹
// preDirDepth int: 前置文件夹深度
// maxPreDirDepth int: 前置文件夹最大深度
// relativePath string: 相对路径
// f func(ta *TraverseArgs) error: 回调函数
func TraversePath(root string, preDir []string,
	preDirDepth int, maxPreDirDepth int, relativePath string,
	f func(ta *TraverseArgs) error) error {

	// 大于1，说明还是前置目录
	if isPreDir(preDirDepth) {
		preDir = append(preDir, path.Base(root))
	}

	files, err := ioutil.ReadDir(root)
	if err != nil {
		logger.Infof("TraversePath failed: %s: %s", root, err.Error())
		return err
	}
	PthSep := string(os.PathSeparator)
	for _, file := range files {
		var fileName = file.Name()
		newFilePath := root + PthSep + fileName
		relPath := relativePath + PthSep + fileName
		if file.IsDir() {
			TraversePath(newFilePath, preDir,
				preDirDepth-1, maxPreDirDepth, relPath, f)
		} else {
			if !isPreDir(preDirDepth) || isLastPreDir(preDirDepth) {
				traverseArgs := TraverseArgs{
					Root:           root,
					PreDir:         preDir,
					PreDirDepth:    preDirDepth,
					RelativePath:   relPath,
					MaxPreDirDepth: maxPreDirDepth,
				}
				f(&traverseArgs)
			}
		}
	}

	if isPreDir(preDirDepth) {
		preDir = preDir[:len(preDir)-1]
	}

	return nil
}

// AdRsrcInfoGenAndAddToDb 添加到数据库
func AdRsrcInfoGenAndAddToDb(ta *TraverseArgs) error {
	if ta.Root == "" || len(ta.PreDir) < ta.MaxPreDirDepth+1 { // 第一级
		// error!
		return fmt.Errorf("Invalid parameters: %+v", ta)
	}

	adRsrc := model.AdRsrcInfo{
		Path:       ta.RelativePath,
		DevVendor:  ta.PreDir[1],
		DevProduct: ta.PreDir[2],
		DevType:    ta.PreDir[3],
		AdURL:      DefaultAdURL,
		RsrcURL:    AdRsrcURI + ta.RelativePath,
	}

	return adRsrc.Add()
	//return nil
}

// AdRsrcInfoUpdate 把资源放到一个固定的目录，然后调用这个接口，
// 就会生成相应的信息到数据库。（全量，之前的信息会删除！）
// curl "http://127.0.0.1:8008/apiv1/user/ad/ad_rsrc_info_update" -v
func AdRsrcInfoUpdate(ctx *fasthttp.RequestCtx) error {
	// 一些必要的校验
	// 删除数据库所有信息
	var tARI model.AdRsrcInfo
	tARI.CanDel = common.TRUE // 允许删除的才删
	cond := fmt.Sprintf("can_del = %d", tARI.CanDel)

	tARI.DelWhere(cond)

	// 遍历RsrcPath && 生成相关信息 && 保存到数据库
	for _, pathInfo := range rsrcPaths {
		preDir := []string{}
		TraversePath(pathInfo.Path, preDir,
			pathInfo.PreDirDepth,
			pathInfo.PreDirDepth,
			"",
			pathInfo.CallBack)
	}

	return nil
}

// AdRsrcInfoSet 修改广告资源信息
// curl "http://127.0.0.1:8008/apiv1/user/ad/ad_rsrc_info_set" -d -v
//		'data={"path":"/Orico/Nas/wrc10/1","dev_vendor":"Orico",
//		"dev_product":"Nas","dev_type":"wrc10","ad_url":"http://baidu.com",
//		"desc":"wrc10-lalalalal"}'
func AdRsrcInfoSet(ctx *fasthttp.RequestCtx) error {
	var gArgs = []string{} //[]string{"token"}
	var pArgs = []string{"data"}
	var Args = common.ArgsGet(ctx, gArgs, pArgs)
	var adRsrcInfo model.AdRsrcInfo

	if err := InfoGetFromArgs(ctx, Args["data"], &adRsrcInfo); err != nil {
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	if adRsrcInfo.Path == "" {
		err := fmt.Errorf("Invalid parameters: %+v", Args)
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	if err := adRsrcInfo.Update(); err != nil {
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	common.ResponseOkJSON(ctx, model.AdRsrcInfo{})
	return nil
}

// AdRsrcInfoDel 删除广告资源信息
// curl "http://127.0.0.1:8008/apiv1/user/ad/ad_rsrc_info_del" -d 'data={"path":"/Orico/Nas/wrc10/1"}' -v
func AdRsrcInfoDel(ctx *fasthttp.RequestCtx) error {
	var gArgs = []string{} //[]string{"token"}
	var pArgs = []string{"data"}
	var Args = common.ArgsGet(ctx, gArgs, pArgs)
	var adRsrcInfo model.AdRsrcInfo

	if err := InfoGetFromArgs(ctx, Args["data"], &adRsrcInfo); err != nil {
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	// 这里对主键做校验，防止删除所有
	if adRsrcInfo.Path == "" {
		err := fmt.Errorf("Invalid parameters: %+v", Args)
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	if err := adRsrcInfo.Del(); err != nil {
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	return nil
}

// AdRsrcInfoGet 获取广告资源信息
// 全部：curl "http://127.0.0.1:8008/apiv1/user/ad/ad_rsrc_info_get" -d 'data={"path":""}' -v
// 单个：curl "http://127.0.0.1:8008/apiv1/user/ad/ad_rsrc_info_get" -d 'data={"path":"/Orico/Nas/wrc10/1"}' -v
func AdRsrcInfoGet(ctx *fasthttp.RequestCtx) error {
	var gArgs = []string{} //[]string{"token"}
	var pArgs = []string{"data"}
	var Args = common.ArgsGet(ctx, gArgs, pArgs)
	var adRsrcInfo model.AdRsrcInfo

	if err := InfoGetFromArgs(ctx, Args["data"], &adRsrcInfo); err != nil {
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	// 主键为空，获取所有，不为空获取一个
	infos, err := adRsrcInfo.Get()

	if err != nil {
		common.ResponseErrJSON(ctx, err.Error())
		return err
	}

	common.ResponseOkJSON(ctx, infos)
	return nil
}
