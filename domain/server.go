package domain

import "context"

type Server interface {
	Start(context.Context) error
}
