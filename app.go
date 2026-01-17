package main

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"sdim_pc/backend/api/convapi"
	"sdim_pc/backend/api/msgapi"
	"sdim_pc/backend/api/userapi"
	"sdim_pc/backend/appctx"
	"sdim_pc/backend/chat"
	"sdim_pc/backend/client"
	"sdim_pc/backend/config"
	"sdim_pc/backend/frmhandler"
	"sdim_pc/backend/mylog"
	preinld "sdim_pc/backend/preinld"
	"sdim_pc/backend/user"
	"sdim_pc/backend/utils"
	"sdim_pc/backend/utils/unet"
	"time"
)

// App struct
type App struct {
	ctx context.Context
	cli *client.Client
	fh  *frmhandler.FrameHandler
	cm  *chat.ConvManager
	ci  *convapi.ConvApi
	ui  *userapi.UserApi
	mi  *msgapi.MsgApi
	cfg config.Config
	lg  *zerolog.Logger
}

func NewApp(cfg config.Config) (*App, error) {
	cli, err := client.NewClient(cfg.AppCfg.EngineServerAddr, nil)
	if err != nil {
		return nil, err
	}

	cm := chat.NewConvManager()
	sender := unet.NewHttpSender(cfg.HttpReqCfg)

	ci := convapi.NewConvApi(cfg, sender)
	ui := userapi.NewUserApi(cfg, sender)
	mi := msgapi.NewMsgApi(cfg, sender)

	fh := frmhandler.NewFrameHandler(cli.GetFrameChan(), cm, ci)

	app := &App{
		cfg: cfg,
		cli: cli,
		fh:  fh,
		cm:  cm,
		ci:  ci,
		ui:  ui,
		mi:  mi,
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

	err := a.cli.Connect(uid, preinld.Pc)
	if err != nil {
		a.lg.Error().Stack().Err(err).Msgf("uid=%s connect to engine failed", uid)
		return err
	}

	/*// 拉取用户资料
	up, err := a.ui.UserProfile(uid)
	if err != nil {
		a.lg.Error().Stack().Err(err).Msgf("fetch user profile failed afater connected, uid=%s", uid)
		return nil
	}

	user.ModifyUnitInfo(
		up.Nickname,
		up.Avatar,
	)
	*/
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

func (a *App) SendMsg(msd *preinld.MsgSendData) error {
	if msd.ChatType == preinld.GroupChat {
		if msd.ConvId == "" {
			return errors.New("convId is required")
		}
	}

	if msd.ConvId == "" {
		clientId := utils.RandStr(32)
		_msd := *msd
		_msd.ClientId = clientId
		return a.cli.SendMsgFrame(_msd)
	}

	convItems, idx, clientMsgId, ok := a.cm.InsertMsgWhileSend(*msd)
	if !ok {
		return errors.New("can not found conv with id:" + msd.ConvId)
	}

	// 通知js, 会话更新(消息发送中)
	a.cm.EmitConvListUpdateEvent(convItems, idx)

	_msd := *msd
	_msd.ClientId = clientMsgId

	return a.cli.SendMsgFrame(_msd)
}
