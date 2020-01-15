import { run as checkFirmwareRun, emit as checkFirmwareEmit } from "./check_firmware"

let eventMap: Map<string, any>;
export function notify(event: string, param?: any) {
    let emit = eventMap.get(event);
    if (!emit) {
        throw new Error("invalid event " + event)
    }

    emit(event, param)
}

export async function run() {
    eventMap = new Map<string, any>();
    eventMap.set("checkFirmware", checkFirmwareEmit);
    checkFirmwareRun();
}

