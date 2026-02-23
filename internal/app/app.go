package app

import (
	"context"
	"net"

	"shop-auth/internal/services/auth"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	v1 "github.com/legendaryneobatman/shop-proto-repo/gen/go/api/auth/v1"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	log    *logrus.Logger
	server *grpc.Server
	port   string
}

func NewApp(log *logrus.Logger, authHandler *auth.Handler, port string) *App {
	// Interceptors
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived,
			logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) error {
			log.WithField("panic", p).Error("Recovered from panic")
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	server := grpc.NewServer(
		grpc.MaxRecvMsgSize(1024*1024),
		grpc.MaxSendMsgSize(1024*1024),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpts...),
			logging.UnaryServerInterceptor(interceptorLogger(log), loggingOpts...),
		),
	)

	v1.RegisterAuthServiceServer(server, authHandler)

	return &App{
		log:    log,
		server: server,
		port:   port,
	}
}

func (a *App) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		return err
	}

	a.log.WithField("port", a.port).Info("Starting gRPC server")

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.server.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		a.log.Info("Shutting down gRPC server")
		a.server.GracefulStop()
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func interceptorLogger(l *logrus.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		entry := l.WithContext(ctx)

		if len(fields) > 0 {
			logrusFields := make(logrus.Fields, len(fields)/2)
			for i := 0; i < len(fields); i += 2 {
				if i+1 < len(fields) {
					if key, ok := fields[i].(string); ok {
						logrusFields[key] = fields[i+1]
					}
				}
			}
			entry = entry.WithFields(logrusFields)
		}

		entry.Log(convertLevel(lvl), msg)
	})
}

func convertLevel(lvl logging.Level) logrus.Level {
	switch lvl {
	case logging.LevelDebug:
		return logrus.DebugLevel
	case logging.LevelInfo:
		return logrus.InfoLevel
	case logging.LevelWarn:
		return logrus.WarnLevel
	case logging.LevelError:
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}
