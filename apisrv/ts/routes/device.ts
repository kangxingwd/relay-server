import express = require('express')

import db = require('../dao/device')
import { validateCgi } from '../lib/validator'
import { deviceValidator } from "./validator"
import { sendOk, sendError, sendOkTable } from "../lib/utils"
export const router: express.Router = express.Router()

/**********************************************
 * 统计设备数
 * Author：mjq
 * Date：2019-01-10
 */
router.post('/total', async function (req: any, res: any, next: any) {
    try {
        const data = await db.SelectDevicesTotal()
        return sendOk(res, data)
    } catch (e) {
        return sendError(res, e.message)
    }
})


/**********************************************
 * 根据dataType查询设备
 * Author：mjq
 * Date：2019-01-08
 */
router.post('/dataType', async function (req: any, res: any, next: any) {
    const { start, length, dataType } = (req as any).body
    try {
        validateCgi(req.body, deviceValidator.deviceDataType)

        const data = await db.SelectDevicesByDataType({
            start: parseInt(start),
            length: parseInt(length),
            dataType: dataType,
            devId: ''
        })
        return sendOkTable(res, data.recordsFiltered, data.recordsTotal, data.data)
    } catch (e) {
        return sendError(res, e.message)
    }
})

/**********************************************
 * 根据devId查询设备
 * Author：mjq
 * Date：2019-01-08
 */
router.post('/devId', async function (req: any, res: any, next: any) {
    const { start, length, devId } = (req as any).body
    try {
        validateCgi(req.body, deviceValidator.deviceDevId)

        const data = await db.SelectDevicesByDevId({
            start: parseInt(start),
            length: parseInt(length),
            dataType: 'all',
            devId: devId
        })
        return sendOkTable(res, data.recordsFiltered, data.recordsTotal, data.data)
    } catch (e) {
        return sendError(res, e.message)
    }
})

/**********************************************
 * 查询单个设备的role
 * Author：mjq
 * Date：2019-01-08
 */
router.post('/role', async function (req: any, res: any, next: any) {
    const { devId, role } = (req as any).body
    try {
        validateCgi(req.body, deviceValidator.deviceRole)

        const association = await db.selectAssociated(devId, role)

        if ('mcloud' == role) {
            //查询mcloud表
            const data = await db.selectMcloud(devId)
            return sendOk(res, {
                data: data,
                association: association,
            })
        }

        if ('relay' == role) {
            //查询relay表
            const data = await db.selectRelay(devId)
            return sendOk(res, {
                data: data,
                association: association,
            })
        }

        return sendError(res, 'role错误')
    } catch (e) {
        return sendError(res, e.message)
    }
})