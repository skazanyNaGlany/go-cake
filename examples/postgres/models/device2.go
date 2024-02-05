package models

import (
	go_cake "github.com/skazanyNaGlany/go-cake"
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
	go_cake.BaseGoKateModel `json:"-"`

	ID          *int64  `json:"id,omitempty" db:"id"`
	ETag        *string `json:"etag,omitempty" db:"etag"`
	Email       *string `json:"email,omitempty" db:"email"`
	MaxContacts *int64  `json:"max_contacts,omitempty" db:"max_contacts"`
}

func (d *Device2) CreateInstance() go_cake.GoKateModel {
	newObj := Device2{}

	return &newObj
}

/*
// old
type Device2 struct {
	bun.BaseModel           `bun:"table:device2,alias:d2"`
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
*/

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

// func (d *Device2) CreateETag() any {
// 	return nil
// }

// func (d *Device2) GetETag() any {
// 	return nil
// }

// func (d *Device2) SetETag(etag string) error {
// 	return nil
// }
