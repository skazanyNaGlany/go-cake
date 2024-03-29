package postgres

import (
	go_cake "github.com/skazanyNaGlany/go-cake"
	"github.com/skazanyNaGlany/go-cake/utils"
)

type ModelSpecs struct {
	model     go_cake.GoCakeModel
	tagMap    utils.TagMap
	idField   string
	etagField string
	dbPath    string
}
