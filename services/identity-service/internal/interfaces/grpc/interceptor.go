package grpc

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryInterceptors returns the server interceptor chain, outermost first.
//
// Request/response payloads are deliberately never logged: CreateUserRequest
// carries a plaintext password.
func UnaryInterceptors(logger *slog.Logger) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		recoveryInterceptor(logger),
		loggingInterceptor(logger),
	}
}

func recoveryInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorContext(ctx, "panic recovered",
					slog.String("method", info.FullMethod),
					slog.Any("panic", r),
					slog.String("stack", string(debug.Stack())),
				)

				resp, err = nil, status.Error(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

func loggingInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		code := status.Code(err)
		level := slog.LevelInfo

		attrs := []slog.Attr{
			slog.String("method", info.FullMethod),
			slog.String("code", code.String()),
			slog.Duration("duration", time.Since(start)),
		}

		if err != nil {
			attrs = append(attrs, slog.String("error", status.Convert(err).Message()))

			switch code {
			case codes.Internal, codes.Unknown, codes.Unavailable, codes.DataLoss:
				level = slog.LevelError
			default:
				level = slog.LevelWarn
			}
		}

		logger.LogAttrs(ctx, level, "grpc request", attrs...)

		return resp, err
	}
}
