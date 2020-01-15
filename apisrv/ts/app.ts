var createError = require('http-errors');
var express = require('express');
var path = require('path');
var cookieParser = require('cookie-parser');
var logger = require('morgan');
var session = require("express-session");
var RedisStore = require('connect-redis')(session);
var app = express();

//mysql
import mysqlPool = require('./lib/mysqlpool')
import config = require('./config/config')
mysqlPool.initPool(config.database)

// view engine setup
app.set('views', path.join(__dirname, 'views'));

app.use(logger('dev'));
app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

// 设置Express的Session存储中间件
app.use(session({
    secret: '生成session的签名',
    resave: false,
    saveUninitialized: true,
    store: new RedisStore(config.session),
}));

app.use(function (req: any, res: any, next: any) {
    if ('/manage/login/login' == req.path) {
        next()
    } else if (req.session.name) {
        next()
    } else {
        //401:请求要求用户的身份认证
        //403:服务器理解请求客户端的请求，但是拒绝执行此请求
        res.status(401).send('Unauthorized');
    }
})

app.use('/manage/login', require('./routes/login').router)
app.use('/manage/system', require('./routes/system').router)
app.use('/manage/device', require('./routes/device').router)

// catch 404 and forward to error handler
app.use(function (req: any, res: any, next: any) {
    // next(createError(404));
    res.locals.message = '';
    let err = new Error('Not Found')
    next(err)
});

// error handler
app.use(function (err: any, req: any, res: any, next: any) {
    // set locals, only providing error in development
    res.locals.message = err.message;
    console.log(err.message)
    res.locals.error = req.app.get('env') === 'development' ? err : {};

    // render the error page
    res.status(err.status || 500);
    // res.render('error');
    res.end()
});

module.exports = app;