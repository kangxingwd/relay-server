import express = require('express');
import appdeal = require("../lib/appdeal");



export const router: express.Router = express.Router();

router.get('/', async function (req: express.Request, res: express.Response, next: express.NextFunction) {
    try{
        let {action} = req.query
        let ret: any = {}
        if(action === "check_version")
        {
            ret = appdeal.getAppVersion()
        }else{
            throw Error("action not found.")
        }
        res.send({r: 1, d: ret})
    }catch (e) {
        console.log(e.message)
        res.send({r: 0, errmsg: e.toString()})
    }
})