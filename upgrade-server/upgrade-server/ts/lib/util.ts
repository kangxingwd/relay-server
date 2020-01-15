import dateformat = require("dateformat")

export const ErrCode: any = {
    ServerErr: 4000,
    ParamErr: 4001,

    FeedbackErr: 4011,
}

const ErrMap: any = {
    4000: "服务端错误",
    4001: "参数错误",
    4011: "反馈操作错误",
}

export function getTime(formatStr?: string): string {
    let now: Date = new Date(Date.now())
    let mask: string = 'yyyy-mm-dd HH:MM:ss'
    if (formatStr)
        mask = formatStr
    return dateformat(now, mask)
}

export function sendOk(res: any, data: any): void {
    let ret = { status: 0, data: data };
    res.setHeader("Content-Type","application/json");
    res.send(JSON.stringify(ret))
}

export function sendError(res: any, errMsg: string): void {
    let errCode: number = ErrCode.ServerErr;
    let errMsgInfo = errMsg.match(/^\[(\d{4})](.*)/);
    let errDesc = "";
    if(errMsgInfo) {
        errCode = Number(errMsgInfo[1])
        errDesc = errMsgInfo[2]
    }
    console.log("sendError:: " + errMsg)
    let ret = { status: 1, data: { errMsg: `[${ErrMap[errCode]}]` + `${errDesc}`, errCode: errCode} };
    res.setHeader("Content-Type","application/json");
    res.send(JSON.stringify(ret))
}