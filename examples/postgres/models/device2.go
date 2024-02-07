package models

import (
	go_cake "github.com/skazanyNaGlany/go-cake"
	"github.com/skazanyNaGlany/go-cake/utils"
	"github.com/uptrace/bun"
)

/*
CREATE TABLE device2 (
   id              	SERIAL PRIMARY KEY,
   etag          	TEXT    NOT NULL,
   email           	TEXT    NOT NULL,
   max_contacts		INT     NOT NULL
);
*/

type Device2 struct {
	bun.BaseModel           `json:"-" bun:"table:device2"`
	go_cake.BaseGoKateModel `json:"-" bun:"-"`

	ID          *int64  `json:"id,omitempty" bun:"id,pk,autoincrement"`
	ETag        *string `json:"etag,omitempty" bun:"etag"`
	Email       *string `json:"email,omitempty" bun:"email"`
	MaxContacts *int64  `json:"max_contacts,omitempty" bun:"max_contacts"`
}

func (d *Device2) CreateInstance() go_cake.GoKateModel {
	newObj := Device2{}

	return &newObj
}

func (d *Device2) GetID() any {
	return d.ID
}

// func (d *Device2) SetID(id string) error {
// 	_id, err := primitive.ObjectIDFromHex(id)

// 	if err != nil {
// 		return err
// 	}

// 	d.ID = &_id

// 	return nil
// }

func (d *Device2) CreateETag() any {
	etag := utils.StringUtilsInstance.NewUUID()

	d.ETag = &etag

	return d.ETag
}

func (d *Device2) GetETag() any {
	return d.ETag
}

func (d *Device2) SetETag(etag string) error {
	d.ETag = &etag

	return nil
}
