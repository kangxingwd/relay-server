import { getConnectionAsync, transactionAsync } from '../lib/mysqlpool'


export async function selectAccount(name: string): Promise<any> {
    let sql = `select * from account where name = '${name}' limit 1`
    console.log(sql)
    let res = await getConnectionAsync(async conn => await conn.queryAsync(sql))
    if (res.length == 0) {
        throw new Error('无效的账号')
    }
    return res[0]
}


export async function updateAccountPassword(name: string, password: string): Promise<any> {
    let sql = `update account set password = '${password}' where name = '${name}'`
    console.log(sql)
    await getConnectionAsync(async conn => await conn.queryAsync(sql))
}