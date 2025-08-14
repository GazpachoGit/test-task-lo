package asynclog

import (
	"context"
	"log/slog"
	"sync"
)

type LogMessage struct {
	Level slog.Level
	Msg   string
	Args  []any
}

type AsyncLog struct {
	log     *slog.Logger
	msgChan chan LogMessage
	wg      *sync.WaitGroup
}

func NewAsyncLog(handler slog.Handler, size int) *AsyncLog {
	ch := make(chan LogMessage, size)
	log := slog.New(handler)

	return &AsyncLog{
		log:     log,
		msgChan: ch,
	}
}

func StartLogger(al *AsyncLog, wg *sync.WaitGroup) {
	al.wg = wg
	al.wg.Add(1)
	go func() {
		defer al.wg.Done()
		for msg := range al.msgChan {
			al.log.Log(context.Background(), msg.Level, msg.Msg, msg.Args...)
		}

	}()
}

func StopLogger(al *AsyncLog) {

	close(al.msgChan)
}

func (al *AsyncLog) Info(msg string, args ...any) {
	al.msgChan <- LogMessage{
		Level: slog.LevelInfo,
		Msg:   msg,
		Args:  args,
	}
}

func (al *AsyncLog) Warn(msg string, args ...any) {
	al.msgChan <- LogMessage{
		Level: slog.LevelWarn,
		Msg:   msg,
		Args:  args,
	}
}

func (al *AsyncLog) Error(msg string, args ...any) {
	al.msgChan <- LogMessage{
		Level: slog.LevelError,
		Msg:   msg,
		Args:  args,
	}
}

func (al *AsyncLog) Debug(msg string, args ...any) {
	al.msgChan <- LogMessage{
		Level: slog.LevelDebug,
		Msg:   msg,
		Args:  args,
	}
}
