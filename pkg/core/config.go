package core

import "context"

type ConfigStore interface {
	Save(context.Context, *TargetConfig) error
	Find(context.Context, string) (*TargetConfig, error)
}
