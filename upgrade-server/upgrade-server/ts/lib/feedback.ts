import daofb = require("../dao/feedback")
import formidable = require('formidable')
import path = require("path")
import fs = require("fs")
import util = require("../lib/util")
import {addFileInfo} from "../dao/feedback";

export async function commitFeedback(type: number, content: string, contact_info: string, device_id: string, product_class: string,
                                     ac_version: string, account: string, app_version: string, pictures: string): Promise<any> {
    let curTime: string = util.getTime("yyyy-mm-dd HH:MM:ss")
    let picStruct: any
    try{
        picStruct = JSON.parse(pictures)
    }catch (e) {
        console.log(`字符串'pictures'格式错误,pictures= '${pictures}'`)
        throw Error(`字符串'pictures'格式错误`)
    }
    let commitId: number = await daofb.addFeedback(device_id, product_class, ac_version, account, app_version, type, content, contact_info, curTime)
    let picDir: string = __dirname + '/../../data/feedback/'
    for (let picNum in picStruct) {
        if (!(picStruct[picNum].content)) {
            console.log(`无法从'${picNum}'中解析出content`)
            continue
        }
        let timeStr: string = util.getTime("yyyymmddHHMMss")
        let randomStr: string = String(Math.floor(Math.random()*8999 + 1000))
        let extname = picStruct[picNum].extname? ('.'+picStruct[picNum].extname) : ''
        let filename: string = timeStr + randomStr + extname;
        let file_path: string = picDir + filename
        let pic_content: any
        if (/;base64,/.test(picStruct[picNum].content)){
            pic_content = new Buffer(picStruct[picNum].content.split(';base64,')[1], 'base64')
        }
        else{
            pic_content = new Buffer(picStruct[picNum].content,'base64')
        }
        let size: number = pic_content.length
        fs.writeFileSync(file_path, pic_content)
        await addFileInfo(commitId, 'image', filename, file_path, size)
    }
    return commitId
}

export async function getFeedbacks(pageNumber: number, pageSize: number, type?: number, productClass?: string,
                                   acVersion?: string, appVersion?: string, status?: number, handleWay?: number): Promise<any> {
    let fbs: any[] = []
    let total = await daofb.countFeedback(type, productClass, acVersion, appVersion, status, handleWay)
    let rows: any[] = await daofb.getFeedbacks((pageNumber-1)*pageSize, pageSize, type, productClass, acVersion, appVersion, status, handleWay)
    for (let row of rows){
        let fb: any = {}
        fb.id = row.id
        fb.device_id = row.device_id
        fb.product_class = row.product_class
        fb.ac_version = row.ac_version
        fb.account = row.account
        fb.app_version = row.app_version
        fb.type = row.type
        fb.content = row.content
        fb.contact_info = row.contact_info
        fb.commit_time = row.commit_time
        fb.status = row.status
        if (row.handle_way) fb.handle_way = row.handle_way
        fb.handle_desc = row.handle_desc?row.handle_desc:""
        fb.handle_time = row.handle_time?row.handle_time:""
        fbs.push(fb)
    } 
    return {totalPages: Math.floor((total+pageSize-1)/pageSize), curPage: pageNumber, feedbacks: fbs}
}

export async function handle(id: number, status: number, handleWay: number, handleDesc: string): Promise<any> {
    if (status > 3 || status < 2 || !Number.isInteger(status))
        throw Error(`[${util.ErrCode.ParamErr}] 'status'只能从(2,3)中选取`)
    if (handleWay > 2 || handleWay < 1 || !Number.isInteger(handleWay))
        throw Error(`[${util.ErrCode.ParamErr}] 'handle_way'只能从(1,2)中选取`)
    if (status == 3 && handleWay == 2)
        throw Error(`[${util.ErrCode.ParamErr}] '已忽略'的反馈不能标记为'已解决'`)
    if (!await daofb.isExistFB(id))
        throw Error(`[${util.ErrCode.FeedbackErr}] '不存在此反馈[${id}]'`)
    let curTime: string = util.getTime("yyyy-mm-dd HH:MM:ss")
    let updateNum = await daofb.handleUpdate(id, status, handleWay, handleDesc, curTime)
    return updateNum
}

export async function getPics(fbid: number): Promise<any> {
    if (!await daofb.isExistFB(fbid))
        throw Error(`[${util.ErrCode.FeedbackErr}] '不存在此反馈[${fbid}]'`)
    let rows: any[] = await daofb.getFileInfo(fbid, "image")
    let files: any[] = []
    for (let row of rows) {
        let picContent
        try{
            let pic = fs.readFileSync(row.path)
            picContent = new Buffer(pic).toString('base64')
        }catch (e) {
            continue
        }
        let file: any = {}
        file.filename = row.filename
        file.content = picContent
        file.filesize = row.size
        files.push(file)
    }
    return {picnum: files.length, pictures: files}
}