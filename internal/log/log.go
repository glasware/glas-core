package log

import "go.uber.org/zap"

type Logger struct {
	logger *zap.Logger
}

func New() (*Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &Logger{
		logger: logger,
	}, nil
}
