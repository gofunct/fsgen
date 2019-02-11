package pkg

import "go.uber.org/zap"

var l, _ = zap.NewDevelopment()
var (
	L *Logger
)

type Logger struct {
	sug *zap.SugaredLogger
}

func init() {
	zap.ReplaceGlobals(l)
	L.sug = l.Sugar()
}

func (l *Logger) FatalIfErr(err error, key string, msg string) {
	if err != nil {
		l.sug.Fatal(zap.Error(err), zap.String(key, msg))
	}
}

func (l *Logger) DebugIfErr(err error, key string, msg string) {
	if err != nil {
		l.sug.Debug(zap.Error(err), zap.String(key, msg))
	}
}

func (l *Logger) WarnIfErr(err error, key string, msg string) {
	if err != nil {
		l.sug.Warn(zap.Error(err), zap.String(key, msg))
	}
}

func (l *Logger) PanicIfErr(err error, key string, msg string) {
	if err != nil {
		l.sug.Panic(zap.Error(err), zap.String(key, msg))
	}
}
