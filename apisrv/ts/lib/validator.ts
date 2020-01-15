import validator = require("validator")

export function validate(str: string, rules: any): void {
    for (let k in rules) {
        let method = (<any>validator)[k]
        if (typeof method !== "function") {
            throw new Error(`invalid validator method ${k}`)
        }

        let rule = rules[k]
        let [errmsg, param] = [rule.errmsg, rule.param]
        if (Array.isArray(param)) {
            if (!method.call(null, str, ...param))
                throw new Error(errmsg)
        } else if (typeof param == "object") {
            if (!method.call(null, str, param))
                throw new Error(errmsg)
        } else if (!method.call(null, str)) {
            throw new Error(errmsg)
        }
    }
}

export function validateCgi(param: any, cgiConfig: any): void {
    if (!param || !cgiConfig) {
        throw new Error("invalid param")
    }

    for (let k in cgiConfig) {
        let rules = cgiConfig[k]
        let userParam = param[k]
        if (userParam == null || userParam == undefined) {
            if (rules.require === 0) {
                continue
            }
            throw new Error(`缺少参数${k}！`)
        }

        let f = rules.tansform
        if (f) {
            userParam = f(userParam)
            delete rules.tansform
        }

        delete rules.require
        validate(userParam, rules)
    }
}
