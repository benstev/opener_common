package fakers

import "github.com/benstev/opener_common/services/database"

type (
	MPhone struct {
		Dp *database.Phone
		Id string
	}

	Manager struct {
		gates        []*database.Gate
		phones       []*MPhone
		phoneUpdated bool
		gateUpdated  bool
	}
)

func NewManager() *Manager {
	return &Manager{make([]*database.Gate, 0), make([]*MPhone, 0), false, false}
}

func (m *Manager) NGates() int  { return len(m.gates) }
func (m *Manager) NPhones() int { return len(m.phones) }

func (m *Manager) FindPhoneByKey(key int64) *MPhone {
	for _, phone := range m.phones {
		if phone.Dp.Key == key {
			return phone
		}
	}
	return nil
}

func (m *Manager) CheckPhoneFiendlyName(key int64, friendlyName string) bool {
	if p := m.FindPhoneByKey(key); p != nil {
		return p.Dp.FriendlyName == friendlyName
	}
	return false
}

func (m *Manager) CheckPhoneHasActivationRecord(key int64, code int) bool {
	if p := m.FindPhoneByKey(key); p != nil {
		return p.Dp.Activation != nil && p.Dp.Activation.Code == code
	}
	return false
}

func (m *Manager) FindGateById(id string) *database.Gate {
	for _, gate := range m.gates {
		if gate.Id == id {
			return gate
		}
	}
	return nil
}

func (m *Manager) GateUpdated() bool {
	updated := m.gateUpdated
	m.gateUpdated = false
	return updated
}

func (m *Manager) AddGate(id string, g *database.Gate) {
	m.gates = append(m.gates, g)
}

func (m *Manager) ChangeGate(id string, g *database.Gate) {
	gates := make([]*database.Gate, 0)
	for _, gate := range m.gates {
		if gate.Id != id {
			gates = append(gates, gate)
		} else {
			gates = append(gates, g)
			m.gateUpdated = true
		}
	}
	m.gates = gates
}

func (m *Manager) RemoveGate(id string) {
	gates := make([]*database.Gate, 0)
	for _, gate := range m.gates {
		if gate.Id != id {
			gates = append(gates, gate)
		}
	}
	m.gates = gates
}

func (m *Manager) AddPhone(id string, p *database.Phone) {
	m.phones = append(m.phones, &MPhone{p, id})
}

func (m *Manager) ChangePhone(id string, p *database.Phone) {
	for _, phone := range m.phones {
		if phone.Id == id {
			phone.Dp = p
		}
	}
}

func (m *Manager) RemovePhone(id string) {
	phones := make([]*MPhone, 0)
	for _, phone := range m.phones {
		if phone.Id != id {
			phones = append(phones, phone)
		}
	}
	m.phones = phones
}
