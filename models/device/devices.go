package device

import (
	"context"

	"github.com/benstev/opener_common/services/database"
	"github.com/rs/zerolog/log"
)

type (
	DevicesRepo struct {
		ctx context.Context
		db  database.DbIF

		Gates  []*database.Gate
		Phones map[string]*Phone
	}
)

func NewDevicesRepo(ctx context.Context, db database.DbIF) *DevicesRepo {
	r := &DevicesRepo{ctx, db, make([]*database.Gate, 0, 10), map[string]*Phone{}}
	db.ListenToGate(r)
	db.ListenToPhone(r)
	return r
}

func (m *DevicesRepo) AddGate(gateId string, gate *database.Gate) {
	log.Debug().Str("gate", gateId).Msg("Added")
	m.Gates = append(m.Gates, gate)
}

func (m *DevicesRepo) ChangeGate(gateId string, gate *database.Gate) {
	log.Debug().Str("gate", gateId).Msg("Changed")

	gates := make([]*database.Gate, 0, 10)
	for _, g := range m.Gates {
		if g.Id == gateId {
			gates = append(gates, gate)
		} else {
			gates = append(gates, g)
		}
	}
	m.Gates = gates
}

func (m *DevicesRepo) RemoveGate(gateId string) {
	log.Debug().Str("gate", gateId).Msg("Remove")

	gates := make([]*database.Gate, 0, 10)
	for _, g := range m.Gates {
		if g.Id != gateId {
			gates = append(gates, g)
		}
	}
	m.Gates = gates
}

// func (m *ManagerService) GetGate(gateId string) *database.Gate {
// 	for _, g := range m.Gates {
// 		if g.Id == gateId {
// 			return g
// 		}
// 	}
// 	return nil
// }

func (m *DevicesRepo) AddPhone(phoneId string, phone *database.Phone) {
	log.Debug().Str("phone", phoneId).Msg("Added")

	p := NewPhone(phoneId, phone)
	m.Phones[phoneId] = p
}

func (m *DevicesRepo) ChangePhone(phoneId string, phone *database.Phone) {
	log.Debug().Str("phone", phoneId).Msg("Changed")

	if p, ok := m.Phones[phoneId]; ok {
		m.Phones[p.Id].Phone = phone
	}
}

func (m *DevicesRepo) RemovePhone(phoneId string) {
	log.Debug().Str("phone", phoneId).Msg("Remove")
	delete(m.Phones, phoneId)
}
