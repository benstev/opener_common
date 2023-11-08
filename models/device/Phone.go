package device

import (
	"github.com/benstev/opener_common/services/database"
	"github.com/rs/zerolog/log"
)

type (
	PhoneIF interface {
		PNumber() string
		PEmail() string
		POwner() string
	}

	Phone struct {
		Id string
		*database.Phone
	}
)

func NewPhone(phoneId string, p *database.Phone) *Phone {
	return &Phone{phoneId, p}
}

func (p *Phone) PNumber() string { return p.Number }

func (p *Phone) PEmail() string { return p.Email }

func (p *Phone) POwner() string { return p.Owner }

func (p *Phone) SetVerified() {
	// p.db.SetVerified(p.Id)
}

func (p *Phone) Notify() {
	log.Logger.Debug().Str("phone", p.Id).Msg(("send notification"))
}
