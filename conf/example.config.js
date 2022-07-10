// 应用配置
const config = {
    dependency: {
        mysql: {
            username: 'vivy',
            password: 'vivy',
            host: 'localhost',
            port: 3306,
            database: 'vivy',
            option: '?charset=utf8mb4',
            max_idle: 10,
            max_open: 100,
            show_exec_time: false,
            show_sql: false,
        },
        redis: {
            host: 'localhost',
            port: 6379,
            auth: null,
        },
        rsyslog: {
            host: 'localhost',
            port: 514,
            protocol: 'tcp', // 这个配置暂时无用，默认只使用tcp
        },
        mongodb: {
            // 应用里暂时没有用到mongodb。
        },
    },
    http: { // 对外服务的IP和端口
        host: '0.0.0.0',
        port: 8999,
    },
    session: {
        cookie_name: 'PHPSESSID', // session使用的cookie名称
        login_cookie_timeout: 300, // 登录校验阶段超时时间：默认5分钟
        secret_hex: '<random_value>', // cookie加解密、签名密钥
        /* 如果secret_hex不能解析，则默认使用随机密钥。可能会导致一些未知问题。 */
        session_timeout: 7200, // session时间，默认2小时
    },
    setting: {
        debug: false,
        logger_formatter: '[%(levelname)s]%(asctime)s[%(filename)s:%(lineno)s][%(name)s]: %(message)s',
        logger_directory: './logs',  // 如果rsyslog可用，则默认不会使用自带的文件日志服务（会降低很多性能）。
    },
}
