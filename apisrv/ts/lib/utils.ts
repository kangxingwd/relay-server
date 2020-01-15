import { createHash, createHmac } from "crypto"

export function sendOk(res: any, data: any): void {
    let ret = { status: 0, data: data }
    res.send(JSON.stringify(ret))
}

export function sendOkTable(res: any, count: number, num: number, row: any): void {
    let ret = { status: 0, recordsFiltered: count, recordsTotal: num, data: row }
    res.send(JSON.stringify(ret))
}

export function sendError(res: any, errmsg: string): void {
    let ret = { status: 1, msg: errmsg }
    res.send(JSON.stringify(ret))
}

export function md5sum(str: string): string {
    return createHash('md5').update(str).digest("hex")
}