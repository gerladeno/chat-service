package websocketstream

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/gerladeno/chat-service/internal/middlewares"
	eventstream "github.com/gerladeno/chat-service/internal/services/event-stream"
	"github.com/gerladeno/chat-service/internal/types"
)

const (
	writeTimeout = time.Second
	pingPeriod   = 250 * time.Millisecond
)

type eventStream interface {
	Subscribe(ctx context.Context, userID types.UserID) (<-chan eventstream.Event, error)
}

//go:generate options-gen -out-filename=handler_options.gen.go -from-struct=Options
type Options struct {
	pingPeriod time.Duration `default:"3s" validate:"omitempty,min=100ms,max=30s"`

	logger       *zap.Logger     `option:"mandatory" validate:"required"`
	eventStream  eventStream     `option:"mandatory" validate:"required"`
	eventAdapter EventAdapter    `option:"mandatory" validate:"required"`
	eventWriter  EventWriter     `option:"mandatory" validate:"required"`
	upgrader     Upgrader        `option:"mandatory" validate:"required"`
	shutdownCh   <-chan struct{} `option:"mandatory" validate:"required"`
}

type HTTPHandler struct {
	Options
}

func NewHTTPHandler(opts Options) (*HTTPHandler, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validating ws http handler options: %v", err)
	}
	return &HTTPHandler{
		Options: opts,
	}, nil
}

func (h *HTTPHandler) Serve(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()
	userID := middlewares.MustUserID(eCtx)
	ws, err := h.upgrader.Upgrade(eCtx.Response().Writer, eCtx.Request(), eCtx.Response().Header())
	if err != nil {
		return fmt.Errorf("upgrade connection to ws: %v", err)
	}
	closer := newWsCloser(h.logger, ws)
	eventsCh, err := h.eventStream.Subscribe(context.Background(), userID)
	if err != nil {
		return fmt.Errorf("subscribe on event stream: %v", err)
	}
	go func() {
		if err = h.readLoop(ctx, ws); err != nil {
			h.logger.Warn("ws readLoop", zap.Error(err))
		}
	}()
	go func() {
		if err = h.writeLoop(ctx, ws, eventsCh); err != nil {
			h.logger.Warn("ws writeLoop", zap.Error(err))
		}
	}()
	go func() {
		<-h.shutdownCh
		closer.Close(websocket.CloseNormalClosure)
	}()
	return nil
}

// readLoop listen PONGs.
func (h *HTTPHandler) readLoop(_ context.Context, ws Websocket) error {
	var err error
	for {
		_, _, err = ws.NextReader()
		if err != nil {
			return fmt.Errorf("get next reader: %v", err)
		}
		_ = ws.SetReadDeadline(time.Now().Add(pingPeriod))
		ws.SetPongHandler(func(string) error {
			_ = ws.SetReadDeadline(time.Now().Add(pingPeriod))
			return nil
		})
	}
}

// writeLoop listen events and writes them into Websocket.
func (h *HTTPHandler) writeLoop(_ context.Context, ws Websocket, events <-chan eventstream.Event) error {
	var err error
	var adapted any
	var event eventstream.Event
	var w io.WriteCloser
	t := time.NewTicker(pingPeriod)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			_ = ws.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err = ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return fmt.Errorf("ping error: %v", err)
			}
		case event = <-events:
			adapted, err = h.eventAdapter.Adapt(event)
			if err != nil {
				return fmt.Errorf("adapt event: %v", err)
			}
			_ = ws.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err = func() error {
				if w, err = ws.NextWriter(websocket.TextMessage); err != nil {
					return fmt.Errorf("get next writer: %v", err)
				}
				defer func() {
					if err = w.Close(); err != nil {
						h.logger.Warn("ws close error", zap.Error(err))
					}
				}()
				err = JSONEventWriter{}.Write(adapted, w)
				if err != nil {
					return fmt.Errorf("write encoded message to the connection: %v", err)
				}
				return nil
			}(); err != nil {
				return err
			}
		}
	}
}
