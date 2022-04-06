package cmd

import "context"

type Command interface {
	Execute(context.Context) error
}
