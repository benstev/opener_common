package device

import (
	"context"

	"github.com/benstev/opener_common/services/database"
)

type (
	Attributes = map[string]interface{}

	Gate struct {
		device *database.Gate
	}
)

func NewGate(ctx context.Context, device *database.Gate) *Gate {
	return &Gate{device}
}

func (d *Gate) Id() string             { return d.device.Id }
func (d *Gate) Device() *database.Gate { return d.device }
