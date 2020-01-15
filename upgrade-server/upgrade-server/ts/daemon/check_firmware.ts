
import updatelib = require("../lib/updatelib");
import {param as commonParam} from "../config/common";

export let lock: any = 0

export async function run() {
    console.log(`start interval for 'checkFirmware'`)
    await update()
    setInterval(() => emit("checkFirmware"), commonParam.checkFWTime)
}

export function emit(event: string, param?: any) {
    process.nextTick(async () => {
        try{
           await update()
        } catch (e) {
            console.log(e.toString())
        }
    })
}

export async function update(): Promise<any>{
    while(lock === 1){
        await updatelib.sleep(1000)
    }
    lock = 1
    console.log("Start to check and update version.json")
    let ret: boolean = updatelib.readConfig(commonParam.configPath)
    if(!ret){
        console.log(`File '${commonParam.configPath}' is incomplete`)
        lock = 0
        console.log("Failed to update version.json! ")
        return
    }
    let changeDirs: any = updatelib.getChangeDir()

    let changeFirmware: any = updatelib.getVersionMesg(changeDirs)


    let versionMaxMap: any = updatelib.getMaxVersion(changeFirmware)
    console.log(versionMaxMap)

    await updatelib.updateFile(versionMaxMap)
    lock = 0
    console.log("Finished updating version.json! ")
}

