package models

import (
	go_cake "github.com/skazanyNaGlany/go-cake"
	utils "github.com/skazanyNaGlany/go-cake/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	go_cake.BaseGoKateModel `json:"-" bson:"-"`

	ID              *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ETag            *string             `json:"_etag,omitempty" bson:"_etag,omitempty"`
	FilteredField   *string             `json:"filtered_field,omitempty" bson:"filtered_field,omitempty"`
	ProjectedField  *string             `json:"projected_field,omitempty" bson:"projected_field,omitempty"`
	SortedField     *string             `json:"sorted_field,omitempty" bson:"sorted_field,omitempty"`
	SortedField2    *string             `json:"sorted_field2,omitempty" bson:"sorted_field2,omitempty"`
	InsertedField   *string             `json:"inserted_field,omitempty" bson:"inserted_field,omitempty"`
	UpdatedField    *string             `json:"updated_field,omitempty" bson:"updated_field,omitempty"`
	HiddenProjField *string             `json:"hidden_proj_field,omitempty" bson:"hidden_proj_field,omitempty"`
	HiddenField     *string             `json:"hidden_field,omitempty" bson:"hidden_field,omitempty"`
	ErasedField     *string             `json:"erased_field,omitempty" bson:"erased_field,omitempty"`
	RequiredField   *string             `json:"required_field,omitempty" bson:"required_field,omitempty"`
	RequiredField2  *string             `json:"required_field2,omitempty" bson:"required_field2,omitempty"`
}

func (o *Order) CreateInstance() go_cake.GoKateModel {
	newObj := Order{}

	return &newObj
}

func (o *Order) GetID() any {
	return o.ID
}

func (o *Order) SetID(id string) error {
	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	o.ID = &_id

	return nil
}

func (o *Order) CreateETag() any {
	etag := utils.StringUtilsInstance.NewUUID()

	o.ETag = &etag

	return o.ETag
}

func (o *Order) GetETag() any {
	return o.ETag
}

func (o *Order) SetETag(etag string) error {
	o.ETag = &etag

	return nil
}
