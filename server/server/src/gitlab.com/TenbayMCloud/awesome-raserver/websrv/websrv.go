package websrv

import (
	"context"
	"fmt"

	log "github.com/cihub/seelog"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"gitlab.com/TenbayMCloud/awesome-raserver/common"
	mysql "gitlab.com/TenbayMCloud/awesome-raserver/common/dbMgr"
	cfg "gitlab.com/TenbayMCloud/awesome-raserver/config"
	"gitlab.com/TenbayMCloud/awesome-raserver/model"
	redis "gitlab.com/TenbayMCloud/awesome-raserver/websrv/redisMgr"
)

type Context = context.Context

type WebSrv struct {
	ctx    Context
	addr   string
	router *router.Router
}

func setupRouters(r *router.Router) error {
	r.GET("/ping", RepPing)
	r.GET("/apiv1/login", login)
	r.POST("/apiv1/user/relay/:option", relay)
	r.POST("/apiv1/user/mcloud/:option", mcloud)
	r.GET("/apiv1/user/ad/:option", ad)
	r.POST("/apiv1/user/ad/:option", ad)

	return nil
}

func RepPing(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "{\"reponse\": \"pong\"}")
}

func NewWebSrv(ctx Context) *WebSrv {
	return &WebSrv{ctx: ctx}
}

func (ws *WebSrv) Initialize() error {

	log.Infof("Initialize")

	if err := redis.RedisConnectInit(); err != nil {
		log.Infof("RedisConnectInit error, err = %v", err)
		return err
	}

	log.Infof("Mysql Init")
	if err := mysql.MysqlConnectInit(); err != nil {
		log.Infof("MysqlConnectInit error, err = %v", err)
		return err
	}

	log.Infof("router Init")
	r := router.New()
	if err := setupRouters(r); err != nil {
		log.Infof("setupRouters")
		return err
	}
	ws.router = r

	log.Infof("model Init")
	if err := model.AutoMigrate(); err != nil {
		log.Infof("model AutoMigrate error, err = %v", err)
		return err
	}

	log.Infof("dns Init")
	if err := common.InitDnsConfig(); err != nil {
		log.Infof("InitDnsConfig error, err = %v", err)
		return err
	}

	log.Infof("Initialize end")
	return nil
}

func (dm *WebSrv) Name() string {
	return "websrv"
}

func (ws *WebSrv) Run() {
	addr := cfg.GetWebSvr().Port //cfg load --TODO
	server := &fasthttp.Server{
		Handler: ws.router.Handler,
	}

	ctx := ws.ctx
	go func() {
		<-ctx.Done()

		log.Infof("shutdown http server")
		err := server.Shutdown()
		if err != nil {
			log.Errorf("shutdown http server fail %", err.Error())
		}
	}()

	log.Infof("websrv listen on %s! \n", addr)
	if err := server.ListenAndServe(addr); err != nil {
		fmt.Println("websrv exit ! error: ", err.Error())
	}
	log.Infof("stop daemon: %s", ws.Name())

	if err := mysql.MysqlClose(); err != nil {
		log.Errorf("MysqlClose error, err = %v", err)
	}
}
