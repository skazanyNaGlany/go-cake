package models

import (
	"math"
	"strconv"

	go_cake "github.com/skazanyNaGlany/go-cake"
	"github.com/skazanyNaGlany/go-cake/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	go_cake.BaseGoKateModel `json:"-" bson:"-"`

	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ETag         int32              `json:"_etag,omitempty" bson:"_etag,omitempty"`
	Email        string             `json:"email,omitempty" bson:"email,omitempty"`
	MaxContacts  uint64             `json:"max_contacts,omitempty" bson:"max_contacts,omitempty"`
	SomeBoolean2 bool               `json:"some_boolean2,omitempty" bson:"some_boolean2,omitempty"`
}

func (p *Product) CreateInstance() go_cake.GoKateModel {
	newObj := Product{}

	return &newObj
}

func (p *Product) GetID() any {
	return p.ID
}

func (u *Product) SetID(id string) error {
	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	u.ID = _id

	return nil
}

func (p *Product) CreateETag() any {
	etag := utils.RandomUtilsInstance.RandomInt32(1, math.MaxInt32)

	p.ETag = etag

	return p.ETag
}

func (p *Product) GetETag() any {
	return p.ETag
}

func (p *Product) SetETag(etag string) error {
	parsedEtag, err := strconv.ParseInt(etag, 10, 32)

	if err != nil {
		return err
	}

	parsedEtag32 := int32(parsedEtag)

	p.ETag = parsedEtag32

	return nil
}
