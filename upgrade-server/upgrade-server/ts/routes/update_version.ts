import express = require('express');
import {update} from "../daemon/check_firmware";


export const router: express.Router = express.Router();

router.get('/', async function (req: express.Request, res: express.Response, next: express.NextFunction) {
    await update()
    res.send("ok")
})