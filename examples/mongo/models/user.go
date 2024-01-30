package models

import (
	"math"
	"strconv"

	go_cake "github.com/skazanyNaGlany/go-cake"
	"github.com/skazanyNaGlany/go-cake/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	go_cake.BaseGoKateModel `json:"-" bson:"-"`

	ID          *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ETag        *int32              `json:"_etag,omitempty" bson:"_etag,omitempty"`
	Email       *string             `json:"email,omitempty" bson:"email,omitempty"`
	MaxContacts *uint64             `json:"max_contacts,omitempty" bson:"max_contacts,omitempty"`
}

func (u *User) CreateInstance() go_cake.GoKateModel {
	newObj := User{}

	return &newObj
}

func (u *User) GetID() any {
	return u.ID
}

func (u *User) SetID(id string) error {
	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	u.ID = &_id

	return nil
}

func (u *User) CreateETag() any {
	etag := utils.RandomUtilsInstance.RandomInt32(1, math.MaxInt32)

	u.ETag = &etag

	return u.ETag
}

func (u *User) GetETag() any {
	return u.ETag
}

func (u *User) SetETag(etag string) error {
	parsedEtag, err := strconv.ParseInt(etag, 10, 32)

	if err != nil {
		return err
	}

	parsedEtag32 := int32(parsedEtag)

	u.ETag = &parsedEtag32

	return nil
}
