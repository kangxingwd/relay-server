
export const database = {
    host: "mysql",
    user: "root",
    password: "123456",
    connectTimeout: 1000,
    database: "webdb",
    dateStrings: true,
    timezone: "+8:00",
    acquireTimeout: 10000,
    charset: 'utf8mb4'
}

export const session = {
    host: 'redis',
    port: '6379',
    pass: '123456',
    db: 2
}

export const LOG_DIRNAME = '/tmp/apisrv'