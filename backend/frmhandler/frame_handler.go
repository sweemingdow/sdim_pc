package frmhandler

import (
	"context"
	"github.com/panjf2000/ants/v2"
	"sdim_pc/backend/api/convapi"
	"sdim_pc/backend/client/frm"
	"sdim_pc/backend/conv"
	"sdim_pc/backend/mylog"
	"sdim_pc/backend/preinld"
	"sdim_pc/backend/user"
	"sdim_pc/backend/utils/parser/json"
	"sync/atomic"
	"time"
)

type FrameHandler struct {
	frmCh  <-chan *frm.Frame
	stopCh chan struct{}
	cm     *conv.ConvManager
	pool   *ants.Pool
	closed atomic.Bool
	ca     *convapi.ConvApi
}

func NewFrameHandler(frmCh <-chan *frm.Frame, cm *conv.ConvManager, ca *convapi.ConvApi) *FrameHandler {
	pool, _ := ants.NewPool(
		8,
		ants.WithMaxBlockingTasks(128),
		ants.WithPreAlloc(true),
		ants.WithNonblocking(true),
	)

	fh := &FrameHandler{
		frmCh:  frmCh,
		stopCh: make(chan struct{}),
		cm:     cm,
		ca:     ca,
		pool:   pool,
	}

	go fh.receiveFrame()

	return fh
}

func (fh *FrameHandler) StopReceive(_ context.Context) {
	if !fh.closed.CompareAndSwap(false, true) {
		return
	}
	close(fh.stopCh)
}

func (fh *FrameHandler) Cleanup(ctx context.Context, frames []*frm.Frame) {
	for _, frame := range frames {
		fh.handleFrame(frame)
	}
}

func (fh *FrameHandler) receiveFrame() {
	for {
		select {
		case <-fh.stopCh:
			return
		case frame, ok := <-fh.frmCh:
			if !ok {
				return
			}

			mylog.GetLogger().Trace().Msgf("reveive frame data, frameDesc:%s, frameData:%+v", frm.FrameType2desc(frame.Header.Ftype), frame)

			fh.handleFrame(frame)
		}
	}
}

func (fh *FrameHandler) handleFrame(frame *frm.Frame) {
	lg := mylog.GetLogger()
	// 发送的回帧
	if frame.Header.Ftype == frm.SendAck {
		var saf preinld.SendFrameAck
		err := json.Parse(frame.Payload.Body, &saf)
		if err != nil {
			lg.Error().Stack().Err(err).Msg("parse send frame ack body failed")
			return
		}

		lg.Trace().Msgf("parse send frame ack body success, data=%+v", saf.Data)

		if !preinld.IsOk(saf.ErrCode) {
			lg.Warn().Msgf("send frame ack errCode is not ok, errCode=%d, errDesc:%s", saf.ErrCode, saf.ErrDesc)
			return
		}

		// 本地会话不存在, 更新一次会话列表
		if !fh.cm.Exists(saf.Data.ConvId) {
			// 更新会话列表
			fh.submitTask(func() {
				time.Sleep(100 * time.Millisecond)

				convItems, ie := fh.ca.RecentlyConvList(user.GetUid())
				if ie != nil {
					lg.Error().Stack().Err(ie).Msg("call active conv list failed")
					return
				}

				contents, _ := json.Fmt(convItems)
				lg.Trace().Msgf("update conv list success, items:%s", string(contents))

				fh.cm.ReplaceConvList(convItems)
			})
		}

	}
}

func (fh *FrameHandler) submitTask(run func()) {
	if err := fh.pool.Submit(run); err == ants.ErrPoolOverload {
		mylog.GetLogger().Error().Err(err).Msg("too many tasks, check pls")
	}
}
