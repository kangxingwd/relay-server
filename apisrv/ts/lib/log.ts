let winston = require("winston")
require('winston-daily-rotate-file')
import dateformat = require("dateformat")
import { LOG_DIRNAME } from '../config/config'
const { combine, timestamp, label, printf } = winston.format;

let log_debug = new winston.transports.DailyRotateFile({
    level: 'debug',
    filename: `debug-%DATE%.log`,
    datePattern: 'YYYY-MM-DD',
    dirname: LOG_DIRNAME,
    maxSize: '20m',
    maxFiles: '14d',
    // colorize: true
})

let log_error = new winston.transports.DailyRotateFile({
    level: 'error',
    filename: `error-%DATE%.log`,
    datePattern: 'YYYY-MM-DD',
    dirname: LOG_DIRNAME,
    maxSize: '20m',
    maxFiles: '14d',
    colorize: true
})

const myFormat = printf((info: any) => {
    return `${dateformat(info.timestamp, "yyyy-mm-dd HH:MM:ss")} [${info.label}] ${info.level}: ${info.message}`;
})

let log_dev_default = winston.createLogger({
    format: combine(
        label({ label: 'dev' }),
        timestamp(),
        myFormat
    ),
    transports: [log_debug, log_error]
});

let log_h5_default = winston.createLogger({
    format: combine(
        label({ label: 'h5' }),
        timestamp(),
        myFormat
    ),
    transports: [log_debug, log_error]
});

let log_pc_default = winston.createLogger({
    format: combine(
        label({ label: 'pc' }),
        timestamp(),
        myFormat
    ),
    transports: [log_debug, log_error]
});

export let log_dev = log_dev_default
export let log_h5 = log_h5_default
export let log_pc = log_pc_default

export function logControl(dev: any, pc: any, h5: any) {
    if (dev && pc && h5) {
        log_dev = log_dev_default
        log_h5 = log_h5_default
        log_pc = log_pc_default
    }
    else {
        if (!dev) {
            log_dev = {
                debug: console.log,
                error: console.log
            }
        }
        if (!pc) {
            log_pc = {
                debug: console.log,
                error: console.log
            }
        }
        if (!h5) {
            log_h5 = {
                debug: console.log,
                error: console.log
            }
        }
    }
}