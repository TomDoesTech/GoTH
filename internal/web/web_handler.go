package web

import "go.uber.org/zap"

type WebHandler struct {
	logger *zap.Logger
}

type WebHandlerParams struct {
	Logger *zap.Logger
}

func NewWebHandler(p WebHandlerParams) *WebHandler {
	return &WebHandler{
		logger: p.Logger,
	}
}
