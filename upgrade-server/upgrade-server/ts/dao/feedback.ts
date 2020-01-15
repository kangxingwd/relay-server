import {getConnectionAsync} from "../lib/mysqlpool";

export async function addFeedback(devId: string, productClass: string, acVersion: string, account: string, appVersion: string,
                                  type: number, content: string, contactInfo: string, commitTime: string): Promise<any> {
    let sql: string = `insert into feedback(device_id,product_class,ac_version,account,app_version,type,content,contact_info,commit_time,status) values('${devId}','${productClass}','${acVersion}','${account}','${appVersion}',${type},'${content}','${contactInfo}','${commitTime}',1)`
    let ret = await getConnectionAsync(async conn => await conn.queryAsync(sql))
    if(ret && ret.insertId) {
        return ret.insertId
    }
    return -1
}

export async function addFileInfo(fid: number, type: string, filename: string, path: string, size: number): Promise<any> {
    let sql: string = `insert into fb_attachment(fid, type, filename, path, size) values(${fid},'${type}','${filename}','${path}',${size})`
    let ret = await getConnectionAsync(async conn => await conn.queryAsync(sql))
    if(ret && ret.insertId) {
        return ret.insertId
    }
    return -1
}

export async function countFeedback(type?: number, productClass?: string, acVersion?: string,
                                    appVersion?: string, status?: number, handleWay?: number): Promise<any> {
    let sql: string = `select count(*) as total from feedback where 1=1`
    sql = type? (sql+` and type=${type}`) : sql
    sql = productClass? (sql+` and product_class='${productClass}'`) : sql
    sql = acVersion? (sql+` and ac_version='${acVersion}'`) : sql
    sql = appVersion? (sql+` and app_version='${appVersion}'`) : sql
    sql = status? (sql+` and status=${status}`) : sql
    sql = handleWay? (sql+` and handle_way=${handleWay}`) : sql
    let rows = await getConnectionAsync(async conn => await conn.queryAsync(sql))
    return rows[0].total
}

export async function getFeedbacks(start: number, limit: number, type?: number, productClass?: string, acVersion?: string,
                                    appVersion?: string, status?: number, handleWay?: number): Promise<any> {
    let sql: string = `select * from feedback where 1=1`
    sql = type? (sql+` and type=${type}`) : sql
    sql = productClass? (sql+` and product_class='${productClass}'`) : sql
    sql = acVersion? (sql+` and ac_version='${acVersion}'`) : sql
    sql = appVersion? (sql+` and app_version='${appVersion}'`) : sql
    sql = status? (sql+` and status=${status}`) : sql
    sql = handleWay? (sql+` and handle_way=${handleWay}`) : sql
    sql = sql + ` limit ${start},${limit}`
    let rows = await getConnectionAsync(async conn => await conn.queryAsync(sql))
    return rows
}

export async function isExistFB(id: number): Promise<boolean> {
    let sql: string = `select count(*) as total from feedback where id = ${id}`
    let rows = await getConnectionAsync(async conn => await conn.queryAsync(sql))
    return rows[0].total>0? true : false
}

export async function handleUpdate(id: number, status: number, handleWay: number, handleDesc: string, handleTime: string): Promise<number> {
    let sql: string = `update feedback set status=${status},handle_way=${handleWay},handle_desc='${handleDesc}',handle_time='${handleTime}' where id = ${id}`
    let ret = await getConnectionAsync(async conn => await conn.queryAsync(sql))
    return ret.affectedRows
}

export async function getFileInfo(fbid: number, filetype: string): Promise<any> {
    let sql: string = `select * from fb_attachment where fid = ${fbid} and type = '${filetype}'`
    let rows = await getConnectionAsync(async conn => await conn.queryAsync(sql))
    return rows
}