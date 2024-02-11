package models

import (
	"math"
	"strconv"

	go_cake "github.com/skazanyNaGlany/go-cake"
	"github.com/skazanyNaGlany/go-cake/utils"
	"github.com/uptrace/bun"
)

/*
CREATE TABLE user2 (
   id              	SERIAL PRIMARY KEY,
   etag          	INT    NOT NULL,
   email           	TEXT   NOT NULL,
   max_contacts		INT    NOT NULL
);
*/

type User2 struct {
	bun.BaseModel           `json:"-" bun:"table:user2"`
	go_cake.BaseGoKateModel `json:"-" bun:"-"`

	ID          *int64  `json:"id,omitempty" bun:"id,pk,autoincrement"`
	ETag        *int32  `json:"etag,omitempty" bun:"etag"`
	Email       *string `json:"email,omitempty" bun:"email"`
	MaxContacts *int64  `json:"max_contacts,omitempty" bun:"max_contacts"`
}

func (u *User2) CreateInstance() go_cake.GoKateModel {
	newObj := User2{}

	return &newObj
}

func (u *User2) GetID() any {
	return u.ID
}

func (u *User2) SetID(id string) error {
	i, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		return err
	}

	u.ID = &i

	return nil
}

func (u *User2) CreateETag() any {
	etag := utils.RandomUtilsInstance.RandomInt32(1, math.MaxInt32)

	u.ETag = &etag

	return u.ETag
}

func (u *User2) GetETag() any {
	return u.ETag
}

func (u *User2) SetETag(etag string) error {
	parsedEtag, err := strconv.ParseInt(etag, 10, 32)

	if err != nil {
		return err
	}

	parsedEtag32 := int32(parsedEtag)

	u.ETag = &parsedEtag32

	return nil
}
