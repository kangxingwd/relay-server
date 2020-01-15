import { getConnectionAsync, transactionAsync } from '../lib/mysqlpool'

interface paramType {
    start: number           //开始位置
    length: number          //长度
    dataType?: string        //全部，在线设备，离线设备，mcloud设备，relay设备
    devId: string           //设备ID
}


export async function SelectDevicesTotal(): Promise<any> {
    let sqlAll = `select ifnull(count(*), 0) num from devices`
    let sqlConn = `select conn_state, ifnull(count(dev_id), 0) num from (select dev_id, sum(conn_state) conn_state from (select dev_id, conn_state from mclouds union select dev_id, conn_state from relays) acstate group by dev_id) s GROUP BY conn_state`
    let sqlMcloud = 'select ifnull(count(*), 0) num from mclouds'
    let sqlRelay = 'select ifnull(count(*), 0) num from relays'
    let rowsAll = await getConnectionAsync(async conn => await conn.queryAsync(sqlAll))
    let rowsConn = await getConnectionAsync(async conn => await conn.queryAsync(sqlConn)) as Array<any>
    let rowsMcloud = await getConnectionAsync(async conn => await conn.queryAsync(sqlMcloud))
    let rowsRelay = await getConnectionAsync(async conn => await conn.queryAsync(sqlRelay))
    let onlineNum: number = 0
    let OfflineNum: number = 0
    rowsConn.forEach(element => {
        if (1 == element.conn_state) {
            onlineNum = element.num
        } else {
            OfflineNum = element.num
        }
    })

    return {
        allNum: rowsAll[0].num,
        onlineNum: onlineNum,
        OfflineNum: OfflineNum,
        mcloudNum: rowsMcloud[0].num,
        relayNum: rowsRelay[0].num
    }
}

export async function SelectDevicesByDataType(opt: paramType): Promise<any> {
    const { start, length, dataType } = opt

    let sqlFrom: string = ''
    switch (dataType) {
        case 'all':
            //查询全部设备
            sqlFrom = `from devices d`
            break
        case 'online':
            //查询在线设备, mcloud和relay有一个在线，为在线
            //union: 合并时去重， conn_state不会为2
            sqlFrom = `from devices d, (select dev_id, sum(conn_state) conn_state from (select dev_id, conn_state from mclouds union select dev_id, conn_state from relays) acstate group by dev_id) s where d.dev_id = s.dev_id and s.conn_state = 1`
            break
        case 'offline':
            //查询离线设备，mcloud和relay全部离线，为离线
            sqlFrom = `from devices d, (select dev_id, sum(conn_state) conn_state from (select dev_id, conn_state from mclouds union select dev_id, conn_state from relays) acstate group by dev_id) s where d.dev_id = s.dev_id and s.conn_state = 0`
            break
        case 'mcloud':
            //查询mcloud设备
            sqlFrom = `from devices d right join mclouds m on m.dev_id = d.dev_id`
            break
        case 'relay':
            //查询relay设备
            sqlFrom = `from devices d right join relays r on r.dev_id = d.dev_id`
            break
        default:
            throw new Error('查询错误')
    }

    let sqlData = `select d.* ${sqlFrom} order by d.join_time desc limit ${start}, ${length}`
    let sqlNum = `select ifnull(count(*), 0) num ${sqlFrom}`
    console.log(sqlData)
    let rowsNum = await getConnectionAsync(async conn => await conn.queryAsync(sqlNum))
    let rows = await getConnectionAsync(async conn => await conn.queryAsync(sqlData)) as Array<any>
    if (rows.length == 0) {
        return {
            recordsFiltered: 0,
            recordsTotal: 0,
            data: new Array()
        }
    }

    rows.forEach(v => {
        if (v.ext) {
            try {
                v.ext = JSON.parse(v.ext)
                if (!v.ext.firmware) v.ext.firmware = '--'
                if (!v.ext.software) v.ext.software = '--'
                if (!v.ext.release) v.ext.release = '--'
            } catch (e) {
                v.ext = { firmware: '--', software: '--', release: '--', err: e.message }
            }
        } else {
            v.ext = { firmware: '--', software: '--', release: '--' }
        }
        v.hostname = v.hostname ? v.hostname : '--'
        v.vendor = v.vendor ? v.vendor : '--'
    })

    return {
        recordsFiltered: rowsNum[0].num,
        recordsTotal: rowsNum[0].num,
        data: rows
    }
}

function replaceStr(str: string) {
    function replace(str: string, rplStr: string) {
        //正常替换
        let valistr = str.replace(/'/g, "''").replace(/\\/g, "\\\\")
        if (rplStr === '+') {
            valistr = valistr.replace(/\+/g, `${rplStr}+`).replace(/\//g, `${rplStr}/`)
        } else if (rplStr === '/') {
            valistr = valistr.replace(/\+/g, `${rplStr}+`)
        }
        valistr = valistr.replace(/_/g, `${rplStr}_`).replace(/%/g, `${rplStr}%`)
        return {
            valistr: valistr,
            rplStr: rplStr,
        }
    }

    return str ? (str.includes('/') ? replace(str, '+') : replace(str, '/')) : { valistr: '', rplStr: '/' }
}

export async function SelectDevicesByDevId(opt: paramType): Promise<any> {
    let { start, length, devId } = opt

    let devMap = replaceStr(devId)
    let sql = `select * from devices where dev_id like '%${devMap.valistr}%' ESCAPE '${devMap.rplStr}' order by join_time desc limit ${start}, ${length}`
    let sqlNum = `select ifnull(count(*), 0) num from devices where dev_id like '%${devMap.valistr}%' ESCAPE '${devMap.rplStr}'`
    console.log(sql)
    let rowsNum = await getConnectionAsync(async conn => await conn.queryAsync(sqlNum))
    let rows = await getConnectionAsync(async conn => await conn.queryAsync(sql))
    if (rows.length == 0) {
        return {
            recordsFiltered: 0,
            recordsTotal: 0,
            data: new Array()
        }
    }
    return {
        recordsFiltered: rowsNum[0].num,
        recordsTotal: rowsNum[0].num,
        data: rows
    }
}

/**********************************************
 * 查询关联设备
 * Author：mjq
 * Date：2019-01-08
 */
export async function selectAssociated(devId: string, dataType: string) {
    let sqlWhere: string = ''
    switch (dataType) {
        case 'mcloud':
            sqlWhere = `where mcloud_id = '${devId}'`
            break
        case 'relay':
            sqlWhere = `where relay_id = '${devId}'`
            break
        default:
            throw new Error('类型错误')
    }
    let sql = `select * from rcmaps ${sqlWhere}`
    let rows = await getConnectionAsync(async conn => await conn.queryAsync(sql))
    if (rows.length == 0) {
        return new Array()
    }
    return rows
}

/**********************************************
 * 查询单个mcloud信息
 * Author：mjq
 * Date：2019-01-08
 */
export async function selectMcloud(devId: string) {
    let sql = `select * from mclouds where dev_id = '${devId}' limit 1`
    let rows = await getConnectionAsync(async conn => await conn.queryAsync(sql))
    if (rows.length == 0) {
        return new Array()
    }
    return rows[0]
}

/**********************************************
 * 查询单个relay信息
 * Author：mjq
 * Date：2019-01-08
 */
export async function selectRelay(devId: string) {
    let sql = `select * from relays where dev_id = '${devId}' limit 1`
    let rows = await getConnectionAsync(async conn => await conn.queryAsync(sql))
    if (rows.length == 0) {
        return new Array()
    }
    return rows[0]
}

