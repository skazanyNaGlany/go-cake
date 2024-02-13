package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/auxten/postgresql-parser/pkg/sql/parser"
	"github.com/auxten/postgresql-parser/pkg/sql/sem/tree"
	"github.com/auxten/postgresql-parser/pkg/walk"
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

// New PostgresDriver using github.com/uptrace/bun/driver/pgdriver driver
// NOTE pgdriver does not support LastInsertId(), it will fill ID field
// value automatically
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

	// log.Println("tagMap", tagMap)

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

func (pd *PostgresDriver) modelSpecsJSONToBUNField(
	jsonField string,
	modelSpecs *ModelSpecs) string {
	for _, specs := range modelSpecs.tagMap {
		if specs["json"] == jsonField {
			return specs["bun"]
		}
	}

	return ""
}

func (pd *PostgresDriver) selectQueryJSONFieldsToBun(
	query *bun.SelectQuery,
	modelSpecs *ModelSpecs) (*bun.SelectQuery, go_cake.HTTPError) {
	statements, err := parser.Parse(query.String())

	if err != nil {
		return nil, go_cake.NewMalformedWhereHTTPError(err)
	}

	// FIRST - replace all fields names from JSON to BUN
	walker := &walk.AstWalker{
		Fn: func(ctx any, node any) (stop bool) {
			unresolvedName, isUnresolvedName := node.(*tree.UnresolvedName)

			if !isUnresolvedName {
				return false
			}

			if unresolvedName.NumParts < 1 {
				return false
			}

			bunName := pd.modelSpecsJSONToBUNField(unresolvedName.Parts[0], modelSpecs)

			if bunName == "" {
				return false
			}

			unresolvedName.Parts[0] = bunName

			return false
		},
	}

	_, err = walker.Walk(statements, nil)

	if err != nil {
		return nil, go_cake.NewMalformedWhereHTTPError(err)
	}

	// log.Println("statements", statements)

	newWhere := ""
	newOrderBy := ""
	countNewWhere := 0
	var offset *int64
	var limit *int64

	// SECOND - grab the new WHERE, ORDER BY and old OFFSET and LIMIT
	walker2 := &walk.AstWalker{
		Fn: func(ctx any, node any) (stop bool) {
			treeWhere, isTreeWhere := node.(*tree.Where)
			treeOrder, isTreeOrder := node.(*tree.Order)
			treeLimit, isTreeLimit := node.(*tree.Limit)

			if isTreeWhere {
				newWhere = treeWhere.Expr.String()

				countNewWhere++
			} else if isTreeOrder {
				formatter := tree.NewFmtCtx(tree.FmtSimple)
				treeOrder.Format(formatter)

				if newOrderBy != "" {
					newOrderBy += ", "
				}

				newOrderBy += formatter.String()
			} else if isTreeLimit {
				if treeLimit.Count != nil {
					parsedLimit, _ := strconv.ParseInt(treeLimit.Count.String(), 10, 64)

					limit = &parsedLimit
				}

				if treeLimit.Offset != nil {
					parsedOffset, _ := strconv.ParseInt(treeLimit.Offset.String(), 10, 64)

					offset = &parsedOffset
				}
			}

			return false
		},
	}

	_, err = walker2.Walk(statements, nil)

	if err != nil {
		return nil, go_cake.NewMalformedWhereHTTPError(err)
	}

	if countNewWhere > 1 {
		return nil, go_cake.NewMalformedWhereHTTPError(nil)
	}

	// log.Println("newWhere", newWhere)
	// log.Println("newOrderBy", newOrderBy)
	// log.Println("offset", offset)
	// log.Println("limit", limit)

	// THIRDLY - rebuild the whole select query since
	// all these ORMs and parsing SQL libs su..s
	translatedQuery := pd.db.NewSelect().Table(modelSpecs.dbPath)

	if newWhere != "" {
		translatedQuery.Where(newWhere)
	}

	if newOrderBy != "" {
		translatedQuery.OrderExpr(newOrderBy)
	}

	if offset != nil {
		translatedQuery.Offset(int(*offset))
	}

	if limit != nil {
		translatedQuery.Limit(int(*limit))
	}

	return translatedQuery, nil
}

func (pd *PostgresDriver) selectQueryGetJSONFields(
	query *bun.SelectQuery,
	modelSpecs *ModelSpecs,
	getWhereFields bool,
	getOrderByFields bool) ([]string, []string, go_cake.HTTPError) {
	statements, err := parser.Parse(query.String())

	if err != nil {
		return nil, nil, go_cake.NewMalformedWhereHTTPError(err)
	}

	whereFields := make([]string, 0)
	orderByFields := make([]string, 0)

	treeWhereNodes := make([]any, 0)
	treeOrderNodes := make([]any, 0)

	walker := &walk.AstWalker{
		Fn: func(ctx any, node any) (stop bool) {
			_, isTreeWhere := node.(*tree.Where)
			_, isTreeOrder := node.(*tree.Order)
			unresolvedName, isUnresolvedName := node.(*tree.UnresolvedName)

			if isTreeWhere {
				treeWhereNodes = make([]any, 0)
				treeWhereNodes = append(treeWhereNodes, node)
			}

			if isTreeOrder {
				treeOrderNodes = make([]any, 0)
				treeOrderNodes = append(treeOrderNodes, node)
			}

			if isUnresolvedName && unresolvedName.NumParts > 0 {
				if len(treeWhereNodes) > 0 && getWhereFields {
					whereFields = append(whereFields, unresolvedName.Parts[0])
				} else if len(treeOrderNodes) > 0 && getOrderByFields {
					orderByFields = append(orderByFields, unresolvedName.Parts[0])
				}
			}

			return false
		},
	}

	_, err = walker.Walk(statements, nil)

	if err != nil {
		return nil, nil, go_cake.NewMalformedWhereHTTPError(err)
	}

	return whereFields, orderByFields, nil
}

func (pd *PostgresDriver) buildSelectQuery(
	modelSpec *ModelSpecs,
	where string,
	sort string,
	page *int64,
	perPage *int64) *bun.SelectQuery {
	query := pd.db.NewSelect().Table(modelSpec.dbPath)

	if where != "" {
		query = query.Where(where)
	}

	if sort != "" {
		query = query.OrderExpr(sort)
	}

	if page != nil && perPage != nil {
		query = query.Offset(int(*perPage) * int(*page)).Limit(int(*perPage))
	}

	return query
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

	query := pd.buildSelectQuery(&modelSpec, where, sort, &page, &perPage)

	translatedQuery, httpErr := pd.selectQueryJSONFieldsToBun(query, &modelSpec)

	if httpErr != nil {
		return nil, httpErr
	}

	err := translatedQuery.Scan(ctx, &resultDocuments)

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

	query := pd.buildSelectQuery(&modelSpec, where, "", nil, nil)

	translatedQuery, httpErr := pd.selectQueryJSONFieldsToBun(query, &modelSpec)

	if httpErr != nil {
		return 0, httpErr
	}

	count, err := translatedQuery.Count(ctx)

	if err != nil {
		return 0, go_cake.NewLowLevelDriverHTTPError(err)
	}

	return uint64(count), nil
}

func (pd *PostgresDriver) Insert(
	model go_cake.GoKateModel,
	documents []go_cake.GoKateModel,
) go_cake.HTTPError {
	if len(documents) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, item := range documents {
		if item.GetHTTPError() != nil {
			continue
		}

		// update etag
		item.CreateETag()

		result, err := pd.db.NewInsert().Model(item).Exec(ctx)

		if err != nil {
			item.SetHTTPError(go_cake.NewLowLevelDriverHTTPError(err))
			continue
		}

		affectedRows, err := result.RowsAffected()

		if err != nil {
			item.SetHTTPError(go_cake.NewLowLevelDriverHTTPError(err))
			continue
		}

		if affectedRows <= 0 {
			item.SetHTTPError(go_cake.NewObjectNotAffectedHTTPError(nil))
			continue
		}

		if affectedRows > 1 {
			item.SetHTTPError(go_cake.NewTooManyAffectedObjectsHTTPError(nil))
			continue
		}
	}

	return nil
}

func (pd *PostgresDriver) Delete(
	model go_cake.GoKateModel,
	documents []go_cake.GoKateModel,
) go_cake.HTTPError {
	if len(documents) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, item := range documents {
		if item.GetHTTPError() != nil {
			continue
		}

		// update etag
		item.CreateETag()

		// TODO update by etag

		result, err := pd.db.NewDelete().Model(item).WherePK().Exec(ctx)

		if err != nil {
			item.SetHTTPError(go_cake.NewLowLevelDriverHTTPError(err))
			continue
		}

		affectedRows, _ := result.RowsAffected()

		if affectedRows <= 0 {
			item.SetHTTPError(go_cake.NewObjectNotFoundHTTPError(nil))
			continue
		}
	}

	return nil
}

func (pd *PostgresDriver) Update(
	model go_cake.GoKateModel,
	documents []go_cake.GoKateModel,
) go_cake.HTTPError {
	if len(documents) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, item := range documents {
		if item.GetHTTPError() != nil {
			continue
		}

		// update etag
		item.CreateETag()

		// TODO update by etag

		result, err := pd.db.NewUpdate().Model(item).WherePK().Exec(ctx)

		if err != nil {
			item.SetHTTPError(go_cake.NewLowLevelDriverHTTPError(err))
			continue
		}

		affectedRows, _ := result.RowsAffected()

		if affectedRows <= 0 {
			item.SetHTTPError(go_cake.NewObjectNotFoundHTTPError(nil))
			continue
		}
	}

	return nil
}

func (pd *PostgresDriver) GetWhereFields(
	model go_cake.GoKateModel,
	where string) ([]string, go_cake.HTTPError) {
	modelType := fmt.Sprintf("%T", model)
	modelSpec := pd.modelJSONTagMap[modelType]

	query := pd.buildSelectQuery(&modelSpec, where, "", nil, nil)

	whereFields, _, httpErr := pd.selectQueryGetJSONFields(query, &modelSpec, true, false)

	if httpErr != nil {
		return nil, httpErr
	}

	return whereFields, nil
}

func (pd *PostgresDriver) GetSortFields(
	model go_cake.GoKateModel,
	sort string) ([]string, go_cake.HTTPError) {
	modelType := fmt.Sprintf("%T", model)
	modelSpec := pd.modelJSONTagMap[modelType]

	query := pd.buildSelectQuery(&modelSpec, "", sort, nil, nil)

	_, orderByFields, httpErr := pd.selectQueryGetJSONFields(query, &modelSpec, false, true)

	if httpErr != nil {
		return nil, httpErr
	}

	return orderByFields, nil
}
