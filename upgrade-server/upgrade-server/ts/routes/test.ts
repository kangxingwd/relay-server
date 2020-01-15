import express = require('express');

export const router: express.Router = express.Router();

router.get('/', async function (req: express.Request, res: express.Response, next: express.NextFunction) {

    res.send("hello world")
})

