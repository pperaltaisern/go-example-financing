package esrcwatermill

import (
	"context"
	"fmt"
	"strings"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Handler[T any] interface {
	Handle(context.Context, *T) error
}

type HandlerWrapper[T any] struct {
	handlerName string
	h           Handler[T]
}

var _ cqrs.CommandHandler = (*HandlerWrapper[any])(nil)
var _ cqrs.EventHandler = (*HandlerWrapper[any])(nil)

func NewHandler[T any](h Handler[T]) *HandlerWrapper[T] {
	return &HandlerWrapper[T]{
		handlerName: handlerName[T](),
		h:           h,
	}
}

func handlerName[T any]() string {
	var h T
	t := fmt.Sprintf("%T", h)
	t = strings.Split(t, ".")[1]
	return t
}

func (h *HandlerWrapper[T]) HandlerName() string {
	return h.handlerName
}

func (h *HandlerWrapper[T]) NewCommand() any {
	return new(T)
}

func (h *HandlerWrapper[T]) NewEvent() any {
	return new(T)
}

func (h *HandlerWrapper[T]) Handle(ctx context.Context, c interface{}) error {
	return h.h.Handle(ctx, c.(*T))
}
