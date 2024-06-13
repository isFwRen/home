const dayjs = require("dayjs");
const chalk = require("chalk");
const common = require("./common");
const logger = {};

const {DURIAN_LOG_TIME, DURIAN_MODE} = require("./env");

function initClient() {
    const log = require("electron-log");
    const levels = ["debug", "log", "info", "warn", "error"];
    const colors = {
        info: "green",
        warn: "yellow",
        error: "red",
    };
    const LV = {}
    levels.map((lv) => {
        LV[lv] = lv.toUpperCase().padEnd(5);
        logger[lv] = async (...args) => {
            // const timestamp = dayjs().format(DURIAN_LOG_TIME);
            let prefix = ''; // `${timestamp} ${LV[lv]}`;
            // if (colors[lv]) prefix = chalk[colors[lv]](prefix);
            args.unshift(prefix);
            const msg = args
                .map((o) => {
                    if (typeof o === "string") return o;
                    return common.inspect(o);
                })
                .join(" ");
            // (lv === "error" ? process.stderr : process.stdout).write(`${msg}\n`);
            log[lv](msg)
        };
    });
    logger.document = (text) => {
        text = text.replace(/^(#+[^\n\r]+)/gm, chalk.green('$1'));
        // .replace(/(\-\s+)((GET|POST|PUT|PATCH|DELETE)\s+\/\S+)/ig, '$1' + chalk.blue('$2'));
        log.debug(text)
    }

}

function initServer() {
    const winston = require("winston")
    require("winston-daily-rotate-file")
    const myFormat = winston.format.printf(log => {
        return `${log.timestamp} ${log.level}: ${log.message}`
    })
    const loggerWinston = winston.createLogger({
        level: 'debug',//输出级别
        format: winston.format.combine(winston.format.timestamp(), myFormat),//输出格式format.json(),format.simple()等也可自定义,timestamp:日志输出时间
        transports: [
            new winston.transports.Console(),//是否在控制台输出
            // new winston.transports.File({filename: './log/error.log', level: 'error'}),//将日志输出到根目录下的log/error.log,输出级别是error，不指定的话采用全局配置的
            // new winston.transports.File({filename: './log/warn.log', level: 'warn'}),//将日志输出到根目录下的log/error.log,输出级别是error，不指定的话采用全局配置的
            // new winston.transports.File({filename: './log/combined.log'}),////将日志输出到根目录下的log/error.log,输出级别使用全局配置的info
            new winston.transports.DailyRotateFile({
                level: "debug",
                dirname: "logs",
                filename: "index-%DATE%.log",
                datePattern: "YYYY-MM-DD-HH",
                maxSize: "512m",
                maxFiles: "7d",
            }),
        ],
    });

    const levels = ["debug", "log", "info", "warn", "error"];
    const colors = {
        info: "green",
        warn: "yellow",
        error: "red",
    };
    const LV = {}
    levels.map((lv) => {
        LV[lv] = lv.toUpperCase().padEnd(5);
        logger[lv] = async (...args) => {
            const timestamp = dayjs().format(DURIAN_LOG_TIME);
            let prefix = `${timestamp} ${LV[lv]}`;
            if (colors[lv]) prefix = chalk[colors[lv]](prefix);
            args.unshift(prefix);
            const msg = args
                .map((o) => {
                    if (typeof o === "string") return o;
                    return common.inspect(o);
                })
                .join(" ");
            // (lv === "error" ? process.stderr : process.stdout).write(`${msg}\n`);
            // log[lv](msg)
            loggerWinston[lv](msg)
        };
    });
    logger.document = (text) => {
        // text = text.replace(/^(#+[^\n\r]+)/gm, chalk.green('$1'));
        // .replace(/(\-\s+)((GET|POST|PUT|PATCH|DELETE)\s+\/\S+)/ig, '$1' + chalk.blue('$2'));
        // log.debug(text)
        loggerWinston.debug(text)
    }
}

if (DURIAN_MODE == "BACKEND") {
    initServer();
} else {
    initClient()
}
module.exports = logger;