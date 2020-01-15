import express = require('express')

import db = require('../dao/login')
import { validateCgi } from '../lib/validator'
import { loginValidator } from "./validator"
import { sendOk, sendError, sendOkTable } from "../lib/utils"

export const router: express.Router = express.Router()


/**********************************************
 * 登录
 * Author：mjq
 * Date：2019-01-07
 */
router.post('/login', async function (req: any, res: any, next: any) {
    const { name, password } = (req as any).body

    try {
        validateCgi(req.body, loginValidator.loginLogin)

        let account = await db.selectAccount(name)
        if (password != account.password) {
            return sendError(res, '密码错误')
        }

        req.session.name = name
        return sendOk(res, 'ok')
    } catch (e) {
        return sendError(res, e.message)
    }
})


/**********************************************
 * 登出
 * Author：mjq
 * Date：2019-01-07
 */
router.post('/loginOut', function (req: any, res: any, next: any) {
    //销毁session
    req.session.destroy()
    return sendOk(res, 'ok')
})

/**********************************************
 * 更改密码
 * Author：mjq
 * Date：2019-01-07
 */
router.post('/update', async function (req: any, res: any, next: any) {
    const { oldPassword, newPassword } = (req as any).body
    const name = req.session.name
    try {
        validateCgi({ name: name, oldPassword: oldPassword, newPassword: newPassword }, loginValidator.loginUpdate)

        let account = await db.selectAccount(name)
        if (oldPassword != account.password) {
            return sendError(res, '密码错误')
        }

        await db.updateAccountPassword(name, newPassword)
        return sendOk(res, 'ok')
    } catch (e) {
        return sendError(res, e.message)
    }
})
