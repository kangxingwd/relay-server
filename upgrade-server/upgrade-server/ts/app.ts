// import path = require('path')
import logger = require('morgan');
import express = require('express');
import cookieParser = require('cookie-parser');
import bodyParser = require('body-parser');

import fs = require("fs");

const app = express();
app.use(logger('dev'));
app.use(cookieParser());
// app.use(bodyParser.json());
// app.use(bodyParser.urlencoded({ extended: false }));

//database
import mysql = require("./lib/mysqlpool");
import mysqlCfg = require("./config/mysql");
mysql.initPool(mysqlCfg.poolOpt)


app.use(bodyParser.json({ limit:"100000kb"}));  //据需求更改limit大小
app.use(bodyParser.urlencoded({ extended: false, limit:"100000kb"})); //根据需求更改limit大小
app.use(bodyParser.raw({ limit:"100000kb"}))
app.use(bodyParser.text({ limit:"100000kb"}))

import appdeal = require('./routes/appdeal');
import test = require('./routes/test');
import acs = require('./routes/update_version');
import feedback = require("./routes/feedback");

app.use('/api/test', test.router);
app.use('/appdeal/handle', appdeal.router)
app.use('/appdeal/update_now', acs.router)
app.use('/api/v1/feedback', feedback.router)


// catch 404 and forward to error handler
app.use(function (req: express.Request, res: express.Response, next: express.NextFunction) {
    res.locals.message = '';
    let err = new Error('Not Found');
    next(err)
});

// error handler
app.use(function (err: any, req: any, res: any, next: any) {
  // set locals, only providing error in development
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  // render the error page
  res.status(err.status || 500);
  res.end()
});

module.exports = app;
require('./daemon/daemon').run();


