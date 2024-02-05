package postgres

import (
	"context"
	"fmt"
	"log"

	go_cake "github.com/skazanyNaGlany/go-cake"
	"github.com/skazanyNaGlany/go-cake/utils"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDriver struct {
	// ConnectionString string
	modelJSONTagMap map[string]ModelSpecs
	// db               *bun.DB
	db *sqlx.DB
}

// func NewPostgresDriver(connectionString string, ctx context.Context) (*PostgresDriver, error) {
func NewPostgresDriver(driverName string, dataSourceName string, ctx context.Context) (*PostgresDriver, error) {
	var err error

	driver := PostgresDriver{
		// ConnectionString: connectionString,
	}

	// //// working but buggy
	// sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connectionString)))
	// driver.db = bun.NewDB(sqldb, pgdialect.New())

	// config, err := pgx.ParseConfig(connectionString)
	// if err != nil {
	// 	panic(err)
	// }
	// // config.PreferSimpleProtocol = true

	// sqldb := stdlib.OpenDB(*config)
	// driver.db = bun.NewDB(sqldb, pgdialect.New())

	driver.db, err = sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		log.Fatalln(err)
	}

	driver.modelJSONTagMap = make(map[string]ModelSpecs)

	return &driver, nil
}

func (pd *PostgresDriver) Close() error {
	// if d.client == nil {
	// 	return nil
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// if err := d.client.Disconnect(ctx); err != nil {
	// 	panic(err)
	// }

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

func (pd *PostgresDriver) Find(
	model go_cake.GoKateModel,
	where, sort string,
	page, perPage int64) ([]go_cake.GoKateModel, go_cake.HTTPError) {
	resultDocuments := make([]go_cake.GoKateModel, 0)

	rows, err := pd.db.Queryx("SELECT * FROM public.device2")

	if err != nil {
		return nil, go_cake.NewLowLevelDriverHTTPError(err)
	}

	for rows.Next() {
		newInstance := model.CreateInstance()

		err := rows.StructScan(newInstance)
		if err != nil {
			newInstance.SetHTTPError(go_cake.NewLowLevelDriverHTTPError(err))
			continue
		}

		resultDocuments = append(resultDocuments, newInstance)
	}

	return resultDocuments, nil
}

func (pd *PostgresDriver) Total(
	model go_cake.GoKateModel,
	where string) (uint64, go_cake.HTTPError) {
	return 0, nil
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
