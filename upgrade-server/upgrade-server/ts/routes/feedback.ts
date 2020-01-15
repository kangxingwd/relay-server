import express = require('express');
import libfb = require("../lib/feedback");
import util = require("../lib/util")


export const router: express.Router = express.Router()

router.post('/handle', async function (req: express.Request, res: express.Response, next: express.NextFunction) {
    try{
        let {id, status, handle_way, handle_desc} = req.body
        if (!(id && status && handle_way && handle_desc)){
            throw Error(`[${util.ErrCode.ParamErr}] 缺少参数`)
        }
        await libfb.handle(id, status, handle_way, handle_desc)
        util.sendOk(res, {})
    }catch (e) {
        util.sendError(res, e.message)
    }
})

router.post('/', async function (req: express.Request, res: express.Response, next: express.NextFunction) {
    try{
        let {type, content, contact_info, device_id, product_class, ac_version, account, app_version, pictures} = req.body
        if (!(type && content && contact_info && device_id && product_class && ac_version && account && app_version && pictures))
        {
            throw Error(`[${util.ErrCode.ParamErr}] 缺少参数`)
        }
        if (isNaN(type))
            throw Error(`[${util.ErrCode.ParamErr}] 参数类型错误`)
        if (Number(type)<1 || Number(type)>3 || !Number.isInteger(Number(type)))
            throw Error(`[${util.ErrCode.ParamErr}] 超出范围`)
        let commitId = await libfb.commitFeedback(Number(type), content, contact_info, device_id, product_class, ac_version, account, app_version, pictures)
        util.sendOk(res, {id: commitId})
    }catch (e) {
        util.sendError(res, e.message)
    }
})

router.get('/pictures/:fbId', async function (req: express.Request, res: express.Response, next: express.NextFunction) {
    try{
        // 参数检查
        let {fbId} = req.params
        if (isNaN(fbId))
            throw Error(`[${util.ErrCode.ParamErr}] 参数类型错误`)
        let ret: any = await libfb.getPics(fbId)
        util.sendOk(res, ret)
    }catch (e) {
        util.sendError(res, e.message)
    }
})

router.get('/', async function (req: express.Request, res: express.Response, next: express.NextFunction) {
    try{
        // 参数检查
        if( (req.query.pageSize && ((isNaN(req.query.pageSize) || Number(req.query.pageSize) < 1 || !Number.isInteger(Number(req.query.pageSize)) ))) ||
            (req.query.pageNumber && ((isNaN(req.query.pageNumber) || Number(req.query.pageNumber) < 1 || !Number.isInteger(Number(req.query.pageNumber))) ))){
            throw new Error(`[${util.ErrCode.ParamErr}]  pageSize或pageNumber 不正确`);
        }
        let pageSize: number = (req.query.pageSize)?Math.floor(Number(req.query.pageSize)):10;
        let pageNumber: number = (req.query.pageNumber)?Math.floor(Number(req.query.pageNumber)):1;
        let query: any = (req.query.query)? (req.query.query) : '{}'
        let {type, product_class, ac_version, app_version, status, handle_way} = JSON.parse(query)
        let ret = await libfb.getFeedbacks(pageNumber, pageSize, type, product_class, ac_version, app_version, status, handle_way)
        util.sendOk(res, ret)
    }catch (e) {
        util.sendError(res, e.message)
    }
})
