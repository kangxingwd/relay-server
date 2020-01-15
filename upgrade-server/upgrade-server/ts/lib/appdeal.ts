import fs = require("fs");
import validator = require("validator");
import {param as commonParam} from "../config/common";

import {filePath} from "./updatelib";

export function getAppVersion(): any{
    let path: string = filePath.appVersionPath
    if (!path)
        path = commonParam.appMsgDefaultPath
    if (!fs.existsSync(path)){
        console.log(`File '${path}' is not exist`)
        throw Error(`File '${path}' is not exist`)
    }
    let versionStr: string = fs.readFileSync(path).toString()
    let versionMesg: any = validator.isJSON(versionStr)? JSON.parse(versionStr) : {}
    if (!(versionMesg.url && versionMesg.version && versionMesg.changelog)){
        console.log(`File '${path}' is wrong`)
        throw Error(`File '${path}' is wrong`)
    }
    return versionMesg
}