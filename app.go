package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"sdim_pc/backend/api/convapi"
	"sdim_pc/backend/appctx"
	"sdim_pc/backend/client"
	"sdim_pc/backend/config"
	"sdim_pc/backend/conv"
	"sdim_pc/backend/frmhandler"
	"sdim_pc/backend/mylog"
	preinld2 "sdim_pc/backend/preinld"
	"sdim_pc/backend/user"
	"sdim_pc/backend/utils/unet"
	"time"
)

// App struct
type App struct {
	ctx context.Context
	cli *client.Client
	fh  *frmhandler.FrameHandler
	cm  *conv.ConvManager
	ci  *convapi.ConvApi
	cfg config.Config
	lg  *zerolog.Logger
}

func NewApp(cfg config.Config) (*App, error) {
	cli, err := client.NewClient(cfg.AppCfg.EngineServerAddr, nil)
	if err != nil {
		return nil, err
	}

	cm := conv.NewConvManager()
	ci := convapi.NewConvApi(cfg, unet.NewHttpSender(cfg.HttpReqCfg))

	fh := frmhandler.NewFrameHandler(cli.GetFrameChan(), cm, ci)

	app := &App{
		cfg: cfg,
		cli: cli,
		fh:  fh,
		cm:  cm,
		ci:  ci,
		lg:  mylog.GetLogger(),
	}
	appctx.RegisterCtcProvider(app)
	return app, nil
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.lg.Info().Msg("application startup hook")
}

func (a *App) domReady(ctx context.Context) {
	a.lg.Info().Msg("application domReady hook")
}

func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	a.lg.Info().Msg("application beforeClose hook")

	_ctx, cancel := context.WithTimeout(ctx, a.cfg.AppCfg.StopTimeout)
	defer cancel()

	a.fh.StopReceive(ctx)

	frames, err := a.cli.Stop(_ctx)
	if err != nil {
		a.lg.Error().Stack().Err(err).Msg("stop client failed")
		return
	}

	if len(frames) > 0 {
		select {
		case <-ctx.Done():
			return
		default:
			a.fh.Cleanup(ctx, frames)
		}
	}

	return
}

func (a *App) shutdown(ctx context.Context) {
	a.lg.Info().Msg("application shutdown hook")
	// 注销所有事件
	runtime.EventsOffAll(ctx)

	time.Sleep(64 * time.Millisecond)
}

func (a *App) WailsCtx() context.Context {
	return a.ctx
}

func (a *App) Conn2Engine(uid string) error {
	a.lg.Debug().Msgf("uid=%s connecting to engine", uid)

	err := a.cli.Connect(uid, preinld2.Pc)
	if err != nil {
		a.lg.Error().Stack().Err(err).Msgf("uid=%s connect to engine failed", uid)
		return err
	}

	return nil
}

func (a *App) Disconnect() error {
	uid := user.GetUid()
	a.lg.Debug().Msgf("uid=%s disconnect", uid)

	err := a.cli.Disconnect()
	if err != nil {
		a.lg.Error().Stack().Err(err).Msgf("uid=%s disconnect failed", uid)
		return err
	}

	user.Reset()

	return nil
}

func (a *App) SendMsg(msd *preinld2.MsgSendData) error {
	return a.cli.SendMsgFrame(msd)
}
