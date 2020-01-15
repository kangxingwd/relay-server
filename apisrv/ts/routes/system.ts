import express = require('express')
import { sendOk, sendError, sendOkTable } from "../lib/utils"
export const router: express.Router = express.Router()
var os = require('os-utils')
// const checkDiskSpace = require('check-disk-space')
let cmd = require('node-cmd')

/**********************************************
 * 获取系统参数
 * Author：mjq
 * Date：2019-01-08
 */
router.post('/', async function (req: any, res: any, next: any) {
    //cpu和内存状态
    let system = await new Promise((resolve, reject) => {
        os.cpuUsage(function (cpu: any) {
            let data = {
                'totalmem': parseInt(os.totalmem()),
                'freemem': parseInt(os.freemem()),
                'freememPercentage': (1 - os.freememPercentage().toFixed(2)) * 100,
                'cpu': parseInt((cpu * 100).toFixed())
            }
            resolve(data)
        })
    })

    //tcp连接状态
    let tcpStatus = await new Promise((resolve, reject) => {
        cmd.get(
            "netstat -nat |awk '{print $6}'|sort|uniq -c",
            function (err: any, data: string, stderr: any) {
                let buf = data.split(" ")
                let arrayData = new Array()
                let obj = JSON.parse('{}')

                buf.forEach(v => {
                    if (v) {
                        v = v.replace(/\n/g, "")
                        arrayData.push(v);
                    }
                })

                arrayData.forEach((v, index) => {
                    if (index % 2 == 0) {
                        let key = arrayData[index + 1]
                        if (key) {
                            obj[key] = parseInt(v)
                        }
                    }
                })
                resolve(obj)
            }
        )
    })

    return sendOk(res, { system: system, tcpStatus: tcpStatus })
})