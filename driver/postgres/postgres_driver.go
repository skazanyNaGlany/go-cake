package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	go_cake "github.com/skazanyNaGlany/go-cake"
	"github.com/skazanyNaGlany/go-cake/utils"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type PostgresDriver struct {
	modelJSONTagMap map[string]ModelSpecs
	db              *bun.DB
}

func NewPostgresDriver(connectionString string, ctx context.Context) (*PostgresDriver, error) {
	driver := PostgresDriver{}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionString)))
	driver.db = bun.NewDB(sqldb, pgdialect.New())

	driver.modelJSONTagMap = make(map[string]ModelSpecs)

	return &driver, nil
}

func (pd *PostgresDriver) Close() error {
	if pd.db == nil {
		return nil
	}

	pd.db.Close()
	pd.db = nil

	return nil
}

func (pd *PostgresDriver) TestModel(
	idField string,
	etagField string,
	model go_cake.GoKateModel,
	dbPath string) error {
	modelType := fmt.Sprintf("%T", model)

	if _, alreadyTested := pd.modelJSONTagMap[modelType]; alreadyTested {
		return nil
	}

	newModelInstance := model.CreateInstance()

	if newModelInstance == nil {
		return fmt.Errorf("%T: unable to create new model instance", model)
	}

	// test ID
	if err := pd.testModelID(model, newModelInstance); err != nil {
		return err
	}

	// test ETag
	if etagField != "" {
		if err := pd.testModelETag(model, newModelInstance); err != nil {
			return err
		}
	}

	// test errors
	if err := pd.testModelError(model, newModelInstance); err != nil {
		return err
	}

	tagMap, err := utils.StructUtilsInstance.StructToTagMap(
		model,
		[]string{"json", "bun", "name"},
		"name")

	if err != nil {
		return err
	}

	if err = pd.testTagMap(idField, etagField, model, tagMap); err != nil {
		return err
	}

	pd.modelJSONTagMap[modelType] = ModelSpecs{
		model:     model,
		tagMap:    tagMap,
		idField:   idField,
		etagField: etagField,
		dbPath:    dbPath,
	}

	return nil
}

func (pd *PostgresDriver) testModelID(
	model go_cake.GoKateModel,
	newModelInstance go_cake.GoKateModel) error {
	return nil
}

func (pd *PostgresDriver) testModelETag(
	model go_cake.GoKateModel,
	newModelInstance go_cake.GoKateModel) error {
	return nil
}

func (pd *PostgresDriver) testModelError(
	model go_cake.GoKateModel,
	newModelInstance go_cake.GoKateModel) error {
	okHttpErr := go_cake.NewOKHTTPError(nil)

	newModelInstance.SetHTTPError(okHttpErr)

	if newModelInstance.GetHTTPError() != okHttpErr {
		return fmt.Errorf("%T: cannot set HTTPError %T", model, okHttpErr)
	}

	return nil
}

func (pd *PostgresDriver) testTagMap(
	idField string,
	etagField string,
	model go_cake.GoKateModel,
	tagMap utils.TagMap) error {

	_, jsonTagExists := tagMap[idField]

	if !jsonTagExists {
		return fmt.Errorf("%T: unable to find JSON ID field tag (%v)", model, idField)
	}

	// idBsonFieldName, bsonTagExists := idJsonData["bson"]

	// if !bsonTagExists || idBsonFieldName == "" {
	// 	return fmt.Errorf("%T: unable to find BSON ID field tag (%v)", model, idField)
	// }

	if etagField == "" {
		return nil
	}

	// ETag field is defined in the resource
	// need to check it
	etagJsonData, jsonTagExists := tagMap[etagField]

	if !jsonTagExists {
		return fmt.Errorf("%T: unable to find JSON ETag field tag (%v)", model, etagField)
	}

	etagBsonFieldName, bsonTagExists := etagJsonData["bson"]

	if !bsonTagExists || etagBsonFieldName == "" {
		return fmt.Errorf("%T: unable to find JSON ETag field tag (%v)", model, etagField)
	}

	return nil
}

func (pd *PostgresDriver) prepareResultDocuments(model go_cake.GoKateModel, howMany int) []go_cake.GoKateModel {
	resultDocuments := make([]go_cake.GoKateModel, 0)

	for howMany > 0 {
		howMany--

		resultDocuments = append(resultDocuments, model.CreateInstance())
	}

	return resultDocuments
}

func (pd *PostgresDriver) Find(
	model go_cake.GoKateModel,
	where, sort string,
	page, perPage int64) ([]go_cake.GoKateModel, go_cake.HTTPError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	modelType := fmt.Sprintf("%T", model)
	modelSpec := pd.modelJSONTagMap[modelType]

	resultDocuments := pd.prepareResultDocuments(model, int(perPage))

	query := pd.db.NewSelect().Table(modelSpec.dbPath)

	if where != "" {
		query = query.Where(where)
	}

	if sort != "" {
		query = query.OrderExpr(sort)
	}

	query = query.Offset(int(perPage) * int(page)).Limit(int(perPage))

	// TODO translate query JSON fields to DB fields

	err := query.Scan(ctx, &resultDocuments)

	if err != nil {
		return nil, go_cake.NewLowLevelDriverHTTPError(err)
	}

	return resultDocuments, nil
}

func (pd *PostgresDriver) Total(
	model go_cake.GoKateModel,
	where string) (uint64, go_cake.HTTPError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	modelType := fmt.Sprintf("%T", model)
	modelSpec := pd.modelJSONTagMap[modelType]

	query := pd.db.NewSelect().Table(modelSpec.dbPath)

	if where != "" {
		query = query.Where(where)
	}

	// TODO translate query JSON fields to DB fields

	count, err := query.Count(ctx)

	if err != nil {
		return 0, go_cake.NewLowLevelDriverHTTPError(err)
	}

	return uint64(count), nil
}

func (pd *PostgresDriver) Insert(
	model go_cake.GoKateModel,
	documents []go_cake.GoKateModel,
) go_cake.HTTPError {
	return nil
}

func (pd *PostgresDriver) Delete(
	model go_cake.GoKateModel,
	documents []go_cake.GoKateModel,
) go_cake.HTTPError {
	return nil
}

func (pd *PostgresDriver) Update(
	model go_cake.GoKateModel,
	documents []go_cake.GoKateModel,
) go_cake.HTTPError {
	return nil
}

func (pd *PostgresDriver) GetWhereFields(where string) ([]string, go_cake.HTTPError) {
	return []string{}, nil
}

func (pd *PostgresDriver) GetSortFields(sort string) ([]string, go_cake.HTTPError) {
	return []string{}, nil
}

func (pd *PostgresDriver) GetProjectionFields(projection string) (map[string]bool, go_cake.HTTPError) {
	return map[string]bool{}, nil
}
