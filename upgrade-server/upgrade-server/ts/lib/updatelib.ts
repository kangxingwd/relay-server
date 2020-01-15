import fs = require("fs");
import validator = require("validator");
import crypto = require("crypto");



export let filePath: any = {
    preDownloadPath: undefined,
    onlineUpgradePath: undefined,
    versionDownloadMesgPath: undefined,
    dirChangeTimePath: undefined,
    appVersionPath: undefined
}

export function readConfig(configPath: string): boolean {
    if (!fs.existsSync(configPath)){
        console.log(`File '${configPath}' is not exist`)
        return false
    }
    let configStr: string = fs.readFileSync(configPath).toString()
    let config: any = JSON.parse(configStr)
    filePath.onlineUpgradePath = config.online_upgrade_path
    filePath.versionDownloadMesgPath = config.version_download_mesg_path
    filePath.dirChangeTimePath = config.dir_change_time_path
    filePath.preDownloadPath = config.pre_download_path
    filePath.appVersionPath = config.app_version_path
    if (!(filePath.onlineUpgradePath && filePath.versionDownloadMesgPath && filePath.dirChangeTimePath && filePath.preDownloadPath))
        return false
    return true
}

export function getDir(): string[] {
    if (!fs.existsSync(filePath.onlineUpgradePath)){
        console.log(`Directory '${filePath.onlineUpgradePath}' is not exist`)
        return []
    }
    let items: string[] = fs.readdirSync(filePath.onlineUpgradePath)
    let dirs: string[] = items.filter(item => {
        return fs.statSync(filePath.onlineUpgradePath + item).isDirectory()
    })
    return dirs
}

export function getChangeDir(): any {
    let dirsFirmware: string[] = getDir()
    let changeDirsFW: any = {}
    let lastTimeMap: any = {}
    if (fs.existsSync(filePath.dirChangeTimePath)){
        let dirChangeTime: string = fs.readFileSync(filePath.dirChangeTimePath).toString()
        lastTimeMap = validator.isJSON(dirChangeTime)? JSON.parse(dirChangeTime) : {}
    }
    for (let dir of dirsFirmware) {
        let changeTime: string = fs.statSync(filePath.onlineUpgradePath + dir).ctime.toString()
        if (lastTimeMap[dir] !== changeTime){
            lastTimeMap[dir] = changeTime
            changeDirsFW[dir] = {changeStatus: 1}
        }else{
            changeDirsFW[dir] = {changeStatus: 0}
        }
    }
    fs.writeFileSync(filePath.dirChangeTimePath, JSON.stringify(lastTimeMap))
    return changeDirsFW
}

export function getVersionMesg(dirs: any): any {
    let changeFirmware: any = dirs
    for (let dir in changeFirmware) {
        changeFirmware[dir].firmware = []
        let firmwares: any = fs.readdirSync(filePath.onlineUpgradePath + dir)
        if (firmwares.length === 0){
            continue
        }
        for (let firmware of firmwares) {
            let img: any = firmware.match(/.(img)$/)? firmware.match(/.(img)$/)[0] : null
            let version: string = firmware.match(/V\d+(\.\d+)+-\d+/)? firmware.match(/V\d+(\.\d+)+-\d+/)[0] : null
            let preVersion: string = version? (version.match(/V(\d+(\.\d+)+)-/)?version.match(/V(\d+(\.\d+)+)-/)[1]:null) : null
            let versionTime: string = version? (version.match(/-(\d+)$/)?version.match(/-(\d+)$/)[1]:null) : null
            if (!(img && version && preVersion && versionTime)) {
                console.log(`File '${firmware}' what we put is wrong`)
                continue
            }
            let data: any = {
                version: version,
                preVersion: preVersion,
                versionTime: versionTime,
                firmware: firmware
            }
            changeFirmware[dir].firmware.push(data)
        }
    }
    return changeFirmware
}

export function getMaxVersion(versionMesgMap: any): any {
    let versionMaxMesgMap: any = {}
    for (let dir in versionMesgMap) {
        versionMaxMesgMap[dir] = {changeStatus: versionMesgMap[dir].changeStatus}
        if (versionMesgMap[dir].firmware.length === 0){
            continue
        }
        versionMaxMesgMap[dir].firmware = versionMesgMap[dir].firmware[0]
        for (let i=1; i<versionMesgMap[dir].firmware.length; i++){
            if (Number(versionMesgMap[dir].firmware[i].versionTime) > Number(versionMaxMesgMap[dir].firmware.versionTime))
                versionMaxMesgMap[dir].firmware = versionMesgMap[dir].firmware[i]
        }
        console.log(`new version for ${dir} >> ${versionMaxMesgMap[dir].firmware.firmware}`)
    }
    return versionMaxMesgMap
}

function getMD5sum(pathFirmware: string): string{
    let buf = fs.readFileSync(pathFirmware)
    let fsHash = crypto.createHash('md5')
    fsHash.update(buf)
    let md5: string = fsHash.digest('hex')
    return md5
}

export let sleep = function (delay_ms: number) {
    return new Promise((resolve, reject) => {
        setTimeout(() => {
            try {
                resolve(1)
            } catch (e) {
                reject(0)
            }
        }, delay_ms);
    })
}

async function checkMD5sum(md5sum: string, pathFirmware: string): Promise<any>{
    let equalTimes: number = 3
    let md5sumNew: string = ""
    let sleepTimeMs: number = 2000
    while (true){
        await sleep(sleepTimeMs)
        md5sumNew = getMD5sum(pathFirmware)
        console.log(`equalTimes=>${equalTimes} md5sum=>${md5sum} md5sumNew=>${md5sumNew}`)
        if (md5sumNew === md5sum) {
            equalTimes--
            sleepTimeMs = 1000
        } else {
            md5sum = md5sumNew
            sleepTimeMs = 2000
        }
        if (equalTimes <= 0)
            return md5sum
    }
    return null
}

export async function updateFile(versionMaxMesgMap: any): Promise<any> {
    let filesData: any[] = []
    if (fs.existsSync(filePath.versionDownloadMesgPath)){
        let versionData: string = fs.readFileSync(filePath.versionDownloadMesgPath).toString()
        filesData = validator.isJSON(versionData)? JSON.parse(versionData) : []
    }
    //清除目录不存在或目录内没有固件的型号信息
    let rmdir: number[] = []
    for (let i = 0; i < filesData.length; i++){
        if(versionMaxMesgMap[filesData[i].devmodel] && versionMaxMesgMap[filesData[i].devmodel].firmware)
            continue
        rmdir.push(i)
    }
    for (let j = rmdir.length-1; j >= 0; j--) {
        console.log(`Clear empty :: ${filesData[rmdir[j]]}`)
        filesData.splice(rmdir[j], 1)
    }

    for (let dir in versionMaxMesgMap) {
        if (versionMaxMesgMap[dir].changeStatus === 0 || !versionMaxMesgMap[dir].firmware)
            continue
        let md5sum: string = getMD5sum(filePath.onlineUpgradePath + dir + "/" + versionMaxMesgMap[dir].firmware.firmware)
        let i: number = 0
        for(i = 0; i < filesData.length; i++){
            if (dir === filesData[i].devmodel)
                break
        }
        let finalMD5sum: any = await checkMD5sum(md5sum, filePath.onlineUpgradePath + dir + "/" + versionMaxMesgMap[dir].firmware.firmware)
        if (i < filesData.length){
            filesData[i].download_path = filePath.preDownloadPath + dir + "/" + versionMaxMesgMap[dir].firmware.firmware
            filesData[i].md5sum = finalMD5sum
            filesData[i].path_firmware = filePath.onlineUpgradePath + dir + "/" + versionMaxMesgMap[dir].firmware.firmware
            filesData[i].version = versionMaxMesgMap[dir].firmware.version
        } else {
        filesData.push({
            devmodel: dir,
            download_path: filePath.preDownloadPath + dir + "/" + versionMaxMesgMap[dir].firmware.firmware,
            md5sum: finalMD5sum,
            path_firmware: filePath.onlineUpgradePath + dir + "/" + versionMaxMesgMap[dir].firmware.firmware,
            version: versionMaxMesgMap[dir].firmware.version
        })
        console.log(`add or update '${dir}' by '${filePath.onlineUpgradePath + dir + "/" + versionMaxMesgMap[dir].firmware}'`)
        }
    }
    fs.writeFileSync(filePath.versionDownloadMesgPath, JSON.stringify(filesData))
    return ""
}


