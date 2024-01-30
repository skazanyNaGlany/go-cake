package models

import (
	go_cake "github.com/skazanyNaGlany/go-cake"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Device struct {
	go_cake.BaseGoKateModel `json:"-" bson:"-"`

	ID          *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email       *string             `json:"email,omitempty" bson:"email,omitempty"`
	MaxContacts *uint64             `json:"max_contacts,omitempty" bson:"max_contacts,omitempty"`
}

func (d *Device) CreateInstance() go_cake.GoKateModel {
	newObj := Device{}

	return &newObj
}

func (d *Device) GetID() any {
	return d.ID
}

func (d *Device) SetID(id string) error {
	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	d.ID = &_id

	return nil
}

func (d *Device) CreateETag() any {
	return nil
}

func (d *Device) GetETag() any {
	return nil
}

func (d *Device) SetETag(etag string) error {
	return nil
}
