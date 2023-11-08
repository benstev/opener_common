package database

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

type (
	DbIF interface {
		UpdatePhone(phoneId string, attributes *PhoneUpdateAttributes) error
		DeletePhone(phoneId string) error
		ActivatePhone(phoneId string, active bool) error
		CreatePhone(phone *PhoneCreatettributes) error
		SetActivation(phoneId string, invite *ActivationRec) error
		AddGatePhone(gateId string, key int64) error
		RemoveGatePhone(gateId string, key int64) error
		SetGateSeed(gateId string, seed uint32) error
		SetRollingCode(gateId string, rollingCode uint32) error
		ListenToGate(listener DbListener)
		ListenToPhone(listener DbListener)
		// SetVerified(phoneId string) error
	}

	DbListener interface {
		AddGate(string, *Gate)
		ChangeGate(string, *Gate)
		RemoveGate(string)
		AddPhone(string, *Phone)
		ChangePhone(string, *Phone)
		RemovePhone(string)
	}

	DeviceError struct {
		ErrorEnum string
		Detail    string
	}

	Gate struct {
		Id                string `firestore:"id" json:"id"`
		DeviceId          string `firestore:"deviceId" json:"deviceId"`
		Name              string `firestore:"name" json:"name"`
		OpenerServiceUuid string `firestore:"openerServiceUuid" json:"openerServiceUuid"`
		CounterUuid       string `firestore:"counterUuid" json:"counterUuid"`
		OpenerUuid        string `firestore:"openerUuid" json:"openerUuid"`

		AdminServiceUuid string `firestore:"adminServiceUuid" json:"adminServiceUuid"`
		PhoneFlasherUuid string `firestore:"phoneFlasherUuid" json:"phoneFlasherUuid"`
		RcFlasherUuid    string `firestore:"rcFlasherUuid" json:"rcFlasherUuid"`
		GetTokenUuid     string `firestore:"getTokenUuid" json:"getTokenUuid"`
		GetFlashUuid     string `firestore:"getFlashUuid" json:"getFlashUuid"`

		RollingCode uint32  `firestore:"rollingCode" json:"rollingCode"`
		Phones      []int64 `firestore:"phones" json:"phones"`
	}

	ActivationRec struct {
		Code int       `firestore:"code"`
		When time.Time `firestore:"when,serverTimestamp"`
	}

	Phone struct {
		Id           string         `firestore:"id" json:"id"`
		Key          int64          `firestore:"key" json:"key"`
		Pin          int            `firestore:"pin,omitempty" json:"pin,omitempty"`
		FriendlyName string         `firestore:"friendlyName" json:"friendlyName"`
		Owner        string         `firestore:"owner" json:"owner"`
		Uid          string         `firestore:"uid,omitempty" json:"uid,omitempty"`
		Number       string         `firestore:"number" json:"number"`
		Active       bool           `firestore:"active" json:"active"`
		Activation   *ActivationRec `firestore:"activation,omitempty" json:"activation,omitempty"`
		Email        string         `firestore:"email" json:"email"`
		// Verified     bool           `firestore:"verified,omitempty" json:"verified,omitempty"`
	}

	PhoneCreatettributes struct {
		Key          int64  `firestore:"key" json:"key"`
		FriendlyName string `firestore:"friendlyName" json:"friendlyName"`
		Owner        string `firestore:"owner" json:"owner"`
		Number       string `firestore:"number" json:"number"`
		Email        string `firestore:"email" json:"email"`
	}

	PhoneUpdateAttributes struct {
		Pin          *int    `json:"pin,omitempty"`
		FriendlyName *string `json:"friendlyName"`
		Owner        *string `json:"owner"`
		Uid          *string `json:"uid,omitempty"`
		Number       *string `json:"number"`
		Email        *string `json:"email"`
	}

	DbService struct {
		client *firestore.Client
		ctx    context.Context
		phones *firestore.CollectionRef
		gates  *firestore.CollectionRef
		grants *firestore.CollectionRef
	}
)

func NewDb(ctx context.Context, fb_Credentials_File string, projectID string) (*DbService, error) {
	opt := option.WithCredentialsFile(fb_Credentials_File)
	config := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	db := &DbService{
		client: client,
		ctx:    ctx,
		phones: client.Collection("phones"),
		gates:  client.Collection("gates"),
		grants: client.Collection("grants"),
	}

	return db, nil
}

func (db *DbService) ListenToGate(listener DbListener) {
	go func() error {
		it := db.gates.Snapshots(db.ctx)
		for {
			snap, err := it.Next()

			if err != nil {
				return err
			}
			if snap != nil {
				for _, change := range snap.Changes {
					id := change.Doc.Ref.ID
					gate := Gate{}
					err = change.Doc.DataTo(&gate)
					if err != nil {
						log.Fatal().Str("gate", id).Err(err).Msg("can't parse")
					}
					gate.Id = id

					switch change.Kind {

					case firestore.DocumentAdded:
						listener.AddGate(id, &gate)

					case firestore.DocumentModified:
						listener.ChangeGate(id, &gate)

					case firestore.DocumentRemoved:
						listener.RemoveGate(id)
					}
				}
			}
		}
	}()
}

func (db *DbService) ListenToPhone(listener DbListener) {
	go func() error {
		it := db.phones.Snapshots(db.ctx)
		for {
			snap, err := it.Next()

			if err != nil {
				return err
			}
			if snap != nil {
				for _, change := range snap.Changes {
					id := change.Doc.Ref.ID
					phone := Phone{}
					err = change.Doc.DataTo(&phone)
					if err != nil {
						log.Fatal().Str("phone", id).Err(err).Msg("can't parse")
					}
					phone.Id = id
					switch change.Kind {

					case firestore.DocumentAdded:
						listener.AddPhone(id, &phone)

					case firestore.DocumentModified:
						listener.ChangePhone(id, &phone)

					case firestore.DocumentRemoved:
						listener.RemovePhone(id)
					}
				}
			}
		}
	}()
}

func (db *DbService) CreatePhone(phone *PhoneCreatettributes) error {

	return db.client.RunTransaction(db.ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		r := db.phones.Where("key", "==", phone.Key).Documents(db.ctx)
		ds, _ := r.GetAll()
		if len(ds) > 0 {
			return fmt.Errorf("key %d already exists", phone.Key)
		}

		_, _, err := db.phones.Add(db.ctx, phone)
		if err != nil {
			log.Error().Err(err).Msg("UpdatePhone")
		}
		return err
	})
}

func setUpdateMap(attributes *PhoneUpdateAttributes) []firestore.Update {
	updates := []firestore.Update{}

	if attributes.Pin != nil {
		updates = append(updates, firestore.Update{Path: "pin", Value: attributes.Pin})
	}
	if attributes.FriendlyName != nil {
		updates = append(updates, firestore.Update{Path: "friendlyName", Value: attributes.FriendlyName})
	}
	if attributes.Owner != nil {
		updates = append(updates, firestore.Update{Path: "owner", Value: attributes.Owner})
	}
	if attributes.Uid != nil {
		updates = append(updates, firestore.Update{Path: "uid", Value: attributes.Uid})
	}
	if attributes.Number != nil {
		updates = append(updates, firestore.Update{Path: "number", Value: attributes.Number})
	}
	if attributes.Email != nil {
		updates = append(updates, firestore.Update{Path: "email", Value: attributes.Email})
	}
	return updates
}

func (db *DbService) UpdatePhone(phoneId string, attributes *PhoneUpdateAttributes) error {
	_, err := db.phones.Doc(phoneId).Update(db.ctx, setUpdateMap(attributes))
	if err != nil {
		log.Error().Err(err).Msg("UpdatePhone")
	}
	return err
}

func (db *DbService) DeletePhone(phoneId string) error {
	_, err := db.phones.Doc(phoneId).Delete(db.ctx)
	if err != nil {
		log.Error().Err(err).Msg("DeletePhone")
	}
	return err
}

func (db *DbService) ActivatePhone(phoneId string, active bool) error {
	updates := []firestore.Update{{Path: "active", Value: active}}
	_, err := db.phones.Doc(phoneId).Update(db.ctx, updates)
	if err != nil {
		log.Error().Err(err).Msg("ActivatePhone")
	}
	return err
}

// func (db *DbService) SetVerified(phoneId string) error {
// 	updates := []firestore.Update{{Path: "verified", Value: true}}
// 	_, err := db.phones.Doc(phoneId).Update(db.ctx, updates)
// 	if err != nil {
// 		log.Error().Err(err)
// 	}
// 	return err
// }

func (db *DbService) CreateGate(Id string, gate *Gate) error {
	_, err := db.gates.Doc(Id).Set(db.ctx, gate)
	if err != nil {
		log.Error().Err(err).Msg("CreateGate")
	}
	return err
}

func (db *DbService) DeleteGate(gateId string) error {
	_, err := db.gates.Doc(gateId).Delete(db.ctx)
	if err != nil {
		log.Error().Err(err).Msg("DeleteGate")
	}
	return err
}

func (db *DbService) SetActivation(phoneId string, invite *ActivationRec) error {
	updates := []firestore.Update{{Path: "activation", Value: invite}}
	_, err := db.phones.Doc(phoneId).Update(db.ctx, updates)
	if err != nil {
		log.Error().Err(err).Msg("SetActivation")
	}
	return err
}

func (db *DbService) AddGatePhone(gateId string, key int64) error {
	updates := []firestore.Update{{Path: "phones", Value: firestore.ArrayUnion(key)}}
	_, err := db.gates.Doc(gateId).Update(db.ctx, updates)
	if err != nil {
		log.Error().Err(err).Msg("AddGatePhone")
	}
	return err
}

func (db *DbService) RemoveGatePhone(gateId string, key int64) error {
	updates := []firestore.Update{{Path: "phones", Value: firestore.ArrayRemove(key)}}
	_, err := db.gates.Doc(gateId).Update(db.ctx, updates)
	if err != nil {
		log.Error().Err(err).Msg("RemoveGatePhone")
	}
	return err
}

func (db *DbService) SetGateSeed(gateId string, seed uint32) error {
	updates := []firestore.Update{{Path: "rollingCode", Value: seed}}
	_, err := db.gates.Doc(gateId).Update(db.ctx, updates)
	if err != nil {
		log.Error().Err(err).Msg("SetGateSeed")
	}
	return err
}

func (db *DbService) SetRollingCode(gateId string, rollingCode uint32) error {
	updates := []firestore.Update{{Path: "rollingCode", Value: rollingCode}}
	_, err := db.gates.Doc(gateId).Update(db.ctx, updates)
	if err != nil {
		log.Error().Err(err).Msg("SetRollingCode	")
	}
	return err
}
