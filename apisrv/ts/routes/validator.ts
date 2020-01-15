
const name = {
    isLength: {
        errmsg: "账号长度错误！",
        param: [4, 64]
    }
}

const password = {
    isLength: {
        errmsg: "密码长度错误！",
        param: [4, 64]
    }
}

const start = {
    tansform: function (v: number) { return "" + v },
    isInt: {
        errmsg: "start字段错误！",
        param: { min: 0, max: 1000000000 }
    }
}

const dataLength = {
    tansform: function (v: number) { return "" + v },
    isInt: {
        errmsg: "查询长度错误！",
        param: { min: 1, max: 1000000000 }
    }
}

const devId = {
    isLength: {
        errmsg: "devId长度错误！",
        param: [0, 128]
    }
}

const dataType = {
    isIn: {
        errmsg: "查询条件错误",
        param: [['all', 'online', 'offline', 'mcloud', 'relay']]
    }
}

const role = {
    isIn: {
        errmsg: "所属角色错误！",
        param: [['mcloud', 'relay']]
    }
}

export const loginValidator = {
    loginLogin: {
        name: name,
        password: password
    },

    loginUpdate: {
        name: name,
        oldPassword: password,
        newPassword: password
    },
}

export const deviceValidator = {
    deviceDataType: {
        start: start,
        length: dataLength,
        dataType: dataType
    },

    deviceDevId: {
        start: start,
        length: dataLength,
        devId: devId
    },

    deviceRole: {
        devId: devId,
        role: role
    },
}