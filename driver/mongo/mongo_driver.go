package mongo_driver

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ghetzel/go-stockutil/maputil"
	go_cake "github.com/skazanyNaGlany/go-cake"
	"github.com/skazanyNaGlany/go-cake/utils"
	"github.com/thoas/go-funk"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ENCODABLE_OBJECT_ID = "64177cafe338354a050543f7"

type MongoDriver struct {
	ConnectionString string
	DatabaseName     string
	client           *mongo.Client
	modelJSONTagMap  map[string]ModelSpecs
}

func NewMongoDriver(connectionString string, databaseName string, ctx context.Context) (*MongoDriver, error) {
	var err error

	driver := MongoDriver{
		ConnectionString: connectionString,
		DatabaseName:     databaseName,
	}

	driver.modelJSONTagMap = make(map[string]ModelSpecs)

	driver.client, err = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))

	if err != nil {
		return nil, &go_cake.UnableToInitDatabaseDriverError{}
	}

	return &driver, nil
}

func (d *MongoDriver) Close() error {
	if d.client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := d.client.Disconnect(ctx); err != nil {
		panic(err)
	}

	return nil
}

func (d *MongoDriver) TestModel(
	idField string,
	etagField string,
	model go_cake.GoKateModel,
	dbPath string) error {
	modelType := fmt.Sprintf("%T", model)

	if _, alreadyTested := d.modelJSONTagMap[modelType]; alreadyTested {
		return nil
	}

	newModelInstance := model.CreateInstance()

	if newModelInstance == nil {
		return fmt.Errorf("%T: unable to create new model instance", model)
	}

	// test ID
	if err := d.testModelID(model, newModelInstance); err != nil {
		return err
	}

	// test ETag
	if etagField != "" {
		if err := d.testModelETag(model, newModelInstance); err != nil {
			return err
		}
	}

	// test errors
	if err := d.testModelError(model, newModelInstance); err != nil {
		return err
	}

	tagMap, err := utils.StructUtilsInstance.StructToTagMap(
		model,
		[]string{"json", "bson", "name"},
		"name")

	if err != nil {
		return err
	}

	if err = d.testTagMap(idField, etagField, model, tagMap); err != nil {
		return err
	}

	d.modelJSONTagMap[modelType] = ModelSpecs{
		model:     model,
		tagMap:    tagMap,
		idField:   idField,
		etagField: etagField,
		dbPath:    dbPath,
	}

	return nil
}

func (d *MongoDriver) testModelID(
	model go_cake.GoKateModel,
	newModelInstance go_cake.GoKateModel) error {
	err := newModelInstance.SetID(ENCODABLE_OBJECT_ID)

	if err != nil {
		return fmt.Errorf("%T: cannot encode ID %v", model, model)
	}

	encodedID := newModelInstance.GetID()

	if encodedID == nil {
		return fmt.Errorf("%T: cannot encode ID %v", model, model)
	}

	finalValue := utils.StructUtilsInstance.GetFinalValue(encodedID)
	finalValueStr := finalValue.(primitive.ObjectID).Hex()

	if finalValueStr != ENCODABLE_OBJECT_ID {
		return fmt.Errorf("%T: cannot encode ID %v", model, model)
	}

	if err := newModelInstance.SetID(finalValueStr); err != nil {
		return fmt.Errorf("%T: cannot encode ID %v", model, model)
	}

	encodedEtag2 := newModelInstance.GetID()

	finalValue2 := utils.StructUtilsInstance.GetFinalValue(encodedEtag2)
	finalValueStr2 := finalValue2.(primitive.ObjectID).Hex()

	if finalValueStr != finalValueStr2 {
		return fmt.Errorf("%T: cannot encode ID %v", model, model)
	}

	return nil
}

func (d *MongoDriver) testModelETag(
	model go_cake.GoKateModel,
	newModelInstance go_cake.GoKateModel) error {
	encodedEtag := newModelInstance.CreateETag()

	if encodedEtag == nil {
		return fmt.Errorf("%T: cannot encode ETag %v", model, model)
	}

	if newModelInstance.GetETag() != encodedEtag {
		return fmt.Errorf("%T: cannot encode ETag %v", model, model)
	}

	finalValue := utils.StructUtilsInstance.GetFinalValue(encodedEtag)
	finalValueStr := fmt.Sprint(finalValue)

	if err := newModelInstance.SetETag(finalValueStr); err != nil {
		return fmt.Errorf("%T: cannot encode ETag %v", model, model)
	}

	encodedEtag2 := newModelInstance.GetETag()

	finalValue2 := utils.StructUtilsInstance.GetFinalValue(encodedEtag2)
	finalValueStr2 := fmt.Sprint(finalValue2)

	if finalValueStr != finalValueStr2 {
		return fmt.Errorf("%T: cannot encode ETag %v", model, model)
	}

	return nil
}

func (d *MongoDriver) testModelError(
	model go_cake.GoKateModel,
	newModelInstance go_cake.GoKateModel) error {
	okHttpErr := go_cake.NewOKHTTPError(nil)

	newModelInstance.SetHTTPError(okHttpErr)

	if newModelInstance.GetHTTPError() != okHttpErr {
		return fmt.Errorf("%T: cannot set HTTPError %T", model, okHttpErr)
	}

	return nil
}

func (d *MongoDriver) testTagMap(
	idField string,
	etagField string,
	model go_cake.GoKateModel,
	tagMap utils.TagMap) error {

	idJsonData, jsonTagExists := tagMap[idField]

	if !jsonTagExists {
		return fmt.Errorf("%T: unable to find JSON ID field tag (%v)", model, idField)
	}

	idBsonFieldName, bsonTagExists := idJsonData["bson"]

	if !bsonTagExists || idBsonFieldName == "" {
		return fmt.Errorf("%T: unable to find BSON ID field tag (%v)", model, idField)
	}

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

func (d *MongoDriver) jsonFieldsToBSONMap(jsonStr string, modelSpecs *ModelSpecs) (map[string]any, error) {
	jsonMap, err := utils.StructUtilsInstance.JSONStringToMap(jsonStr)

	if err != nil {
		return nil, err
	}

	modelNewInstance := modelSpecs.model.CreateInstance()

	err = maputil.Walk(jsonMap, func(value any, path []string, isLeaf bool) error {
		if len(path) == 0 {
			return nil
		}

		key := path[len(path)-1]

		for _, specs := range modelSpecs.tagMap {
			if specs["json"] == key {
				parentPath := path[0 : len(path)-1]
				parentValue := maputil.DeepGet(jsonMap, parentPath).(map[string]any)

				maputil.Delete(parentValue, specs["json"])

				if key == modelSpecs.tagMap[modelSpecs.idField]["json"] {
					valueStr := fmt.Sprintf("%v", value)

					err := modelNewInstance.SetID(valueStr)

					if err != nil {
						return err
					}

					parentValue[specs["bson"]] = utils.StructUtilsInstance.GetFinalValue(modelNewInstance.GetID())
				} else if modelSpecs.etagField != "" {
					if key == modelSpecs.tagMap[modelSpecs.etagField]["json"] {
						valueStr := fmt.Sprintf("%v", value)

						err := modelNewInstance.SetETag(valueStr)

						if err != nil {
							return err
						}

						parentValue[specs["bson"]] = utils.StructUtilsInstance.GetFinalValue(modelNewInstance.GetETag())
					}
				} else {
					parentValue[specs["bson"]] = value
				}

				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

func (d *MongoDriver) jsonWhereToFilter(jsonStr string, modelSpecs *ModelSpecs) (primitive.M, error) {
	jsonMap, err := d.jsonFieldsToBSONMap(jsonStr, modelSpecs)

	if err != nil {
		return nil, err
	}

	return bson.M(jsonMap), nil
}

func (d *MongoDriver) Find(
	model go_cake.GoKateModel,
	where, sort string,
	page, perPage int64) ([]go_cake.GoKateModel, go_cake.HTTPError) {
	var filter bson.M
	var err error

	resultDocuments := make([]go_cake.GoKateModel, 0)

	modelType := fmt.Sprintf("%T", model)
	modelSpec := d.modelJSONTagMap[modelType]

	if where != "" {
		filter, err = d.jsonWhereToFilter(where, &modelSpec)

		if err != nil {
			return nil, go_cake.NewMalformedWhereHTTPError(err)
		}
	}

	options, _, httpErr := d.getFindOptions(sort, page, perPage, &modelSpec)

	if httpErr != nil {
		return nil, httpErr
	}

	collection := d.client.Database(d.DatabaseName).Collection(modelSpec.dbPath)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter, &options)

	if err != nil {
		httpErr := go_cake.NewLowLevelDriverHTTPError(err)

		return nil, httpErr
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		modelNewInstance := model.CreateInstance()

		if err := cursor.Decode(modelNewInstance); err != nil {
			httpErr := go_cake.NewServerObjectMalformedHTTPError(modelNewInstance, err)
			modelNewInstance.SetHTTPError(httpErr)
		}

		resultDocuments = append(resultDocuments, modelNewInstance)
	}

	if err := cursor.Err(); err != nil {
		httpErr := go_cake.NewLowLevelDriverHTTPError(err)

		return nil, httpErr
	}

	return resultDocuments, nil
}

func (d *MongoDriver) Total(
	model go_cake.GoKateModel,
	where string) (uint64, go_cake.HTTPError) {
	var filter bson.M
	var err error
	var count int64

	modelType := fmt.Sprintf("%T", model)
	modelSpec := d.modelJSONTagMap[modelType]

	if where != "" {
		filter, err = d.jsonWhereToFilter(where, &modelSpec)

		if err != nil {
			return 0, go_cake.NewMalformedWhereHTTPError(err)
		}
	}

	collection := d.client.Database(d.DatabaseName).Collection(modelSpec.dbPath)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count, err = collection.CountDocuments(ctx, filter)

	if err != nil {
		httpErr := go_cake.NewLowLevelDriverHTTPError(err)

		return 0, httpErr
	}

	return uint64(count), nil
}

func (d *MongoDriver) Insert(
	model go_cake.GoKateModel,
	documents []go_cake.GoKateModel,
) go_cake.HTTPError {
	if len(documents) == 0 {
		return nil
	}

	modelType := fmt.Sprintf("%T", model)
	modelSpec := d.modelJSONTagMap[modelType]

	collection := d.client.Database(d.DatabaseName).Collection(modelSpec.dbPath)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, item := range documents {
		if item.GetHTTPError() != nil {
			continue
		}

		// update etag
		item.CreateETag()

		result, err := collection.InsertOne(ctx, item)

		if err != nil {
			item.SetHTTPError(go_cake.NewLowLevelDriverHTTPError(err))
			continue
		}

		hexId := result.InsertedID.(primitive.ObjectID).Hex()

		if err := item.SetID(hexId); err != nil {
			item.SetHTTPError(go_cake.NewClientObjectMalformedHTTPError(err))
		}
	}

	return nil
}

func (d *MongoDriver) Delete(
	model go_cake.GoKateModel,
	documents []go_cake.GoKateModel,
) go_cake.HTTPError {
	if len(documents) == 0 {
		return nil
	}

	modelType := fmt.Sprintf("%T", model)
	modelSpec := d.modelJSONTagMap[modelType]

	collection := d.client.Database(d.DatabaseName).Collection(modelSpec.dbPath)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, item := range documents {
		if item.GetHTTPError() != nil {
			continue
		}

		filter, httpErr := d.documentToFilter2(&modelSpec, item)

		if httpErr != nil {
			item.SetHTTPError(httpErr)
			continue
		}

		result, err := collection.DeleteOne(ctx, filter)

		if err != nil {
			item.SetHTTPError(go_cake.NewLowLevelDriverHTTPError(err))
			continue
		}

		if result.DeletedCount > 1 {
			item.SetHTTPError(go_cake.NewTooManyOBjectsHTTPError(nil))
			continue
		}

		if result.DeletedCount < 1 {
			item.SetHTTPError(go_cake.NewObjectNotFoundHTTPError(nil))
			continue
		}
	}

	return nil
}

func (d *MongoDriver) Update(
	model go_cake.GoKateModel,
	documents []go_cake.GoKateModel,
) go_cake.HTTPError {
	if len(documents) == 0 {
		return nil
	}

	modelType := fmt.Sprintf("%T", model)
	modelSpec := d.modelJSONTagMap[modelType]

	collection := d.client.Database(d.DatabaseName).Collection(modelSpec.dbPath)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, item := range documents {
		if item.GetHTTPError() != nil {
			continue
		}

		filter, httpErr := d.documentToFilter2(&modelSpec, item)

		if httpErr != nil {
			item.SetHTTPError(httpErr)
			continue
		}

		// update etag
		item.CreateETag()

		result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": item})

		if err != nil {
			item.SetHTTPError(go_cake.NewLowLevelDriverHTTPError(err))
			continue
		}

		if result.MatchedCount > 1 {
			item.SetHTTPError(go_cake.NewTooManyOBjectsHTTPError(nil))
			continue
		}

		if result.MatchedCount < 1 {
			item.SetHTTPError(go_cake.NewObjectNotFoundHTTPError(nil))
			continue
		}
	}

	return nil
}

func (d *MongoDriver) documentToFilter2(
	modelSpec *ModelSpecs,
	document go_cake.GoKateModel) (map[string]any, go_cake.HTTPError) {
	ifilter := make(map[string]any)
	idValue := document.GetID()

	if idValue == nil {
		return nil, go_cake.NewClientObjectMalformedHTTPError(nil)
	}

	idFinalValue := utils.StructUtilsInstance.GetFinalValue(idValue).(primitive.ObjectID)

	idFieldBSON := modelSpec.tagMap[modelSpec.idField]["bson"]
	etagFieldBSON := ""

	if modelSpec.etagField != "" {
		etagFieldBSON = modelSpec.tagMap[modelSpec.etagField]["bson"]
	} else {
		etagFieldBSON = ""
	}

	ifilter[idFieldBSON] = idFinalValue

	if etagFieldBSON != "" {
		etagFinalValue := utils.StructUtilsInstance.GetFinalValue(document.GetETag())

		if len(fmt.Sprintf("%v", etagFinalValue)) > 0 {
			ifilter[etagFieldBSON] = etagFinalValue
		}
	}

	return ifilter, nil
}

func (d *MongoDriver) GetWhereFields(where string) ([]string, go_cake.HTTPError) {
	var bsonWhere map[string]any

	if err := json.Unmarshal([]byte(where), &bsonWhere); err != nil {
		httpErr := go_cake.NewMalformedWhereHTTPError(err)

		return nil, httpErr
	}

	fields := utils.MapUtilsInstance.GetMapStringKeys(bsonWhere, true)

	reducedFields := funk.FilterString(fields, func(s string) bool {
		return s[0] != '$'
	})

	return reducedFields, nil
}

func (d *MongoDriver) GetSortFields(sort string) ([]string, go_cake.HTTPError) {
	jsonMap, err := utils.StructUtilsInstance.JSONStringToMap(sort)

	if err != nil {
		httpErr := go_cake.NewMalformedSortHTTPError(err)

		return nil, httpErr
	}

	return utils.MapUtilsInstance.GetMapStringKeys(jsonMap, true), nil
}

func (d *MongoDriver) GetProjectionFields(projection string) (map[string]bool, go_cake.HTTPError) {
	var bsonProjection map[string]any

	fields := make(map[string]bool)

	if projection == "" {
		return fields, nil
	}

	if err := json.Unmarshal([]byte(projection), &bsonProjection); err != nil {
		httpErr := go_cake.NewMalformedProjectionHTTPError(err)

		return nil, httpErr
	}

	for field, fieldData := range bsonProjection {
		fieldDataStr := fmt.Sprintf("%v", fieldData)

		fieldDataBool, err := strconv.ParseBool(fieldDataStr)

		if err != nil {
			httpErr := go_cake.NewMalformedProjectionHTTPError(err)

			return nil, httpErr
		}

		fields[field] = fieldDataBool
	}

	return fields, nil
}

func (d *MongoDriver) getFindOptions(
	sortStr string,
	page int64,
	perPage int64,
	modelSpecs *ModelSpecs) (options.FindOptions, map[string]any, go_cake.HTTPError) {
	options := options.FindOptions{}
	optionsMap := make(map[string]any)

	if perPage > 0 {
		options.SetLimit(int64(perPage))
		optionsMap["Limit"] = options.Limit
	}

	if page > 0 {
		options.SetSkip(int64(perPage) * int64(page))
		optionsMap["Skip"] = options.Limit
	}

	if sortStr != "" {
		sort, httpErr := d.getSort(sortStr, modelSpecs)

		if httpErr != nil {
			return options, optionsMap, httpErr
		}

		options.SetSort(sort)
		optionsMap["Sort"] = options.Sort
	}

	return options, optionsMap, nil
}

func (d *MongoDriver) getSort(sortStr string, modelSpecs *ModelSpecs) (bson.D, go_cake.HTTPError) {
	sort := bson.D{}

	if sortStr == "" {
		return sort, nil
	}

	jsonMap, err := utils.StructUtilsInstance.JSONStringToMap(sortStr)

	if err != nil {
		return nil, go_cake.NewMalformedSortHTTPError(err)
	}

	if len(jsonMap) == 0 {
		return sort, nil
	}

	jsonMapStr, err := utils.StructUtilsInstance.StructToJSONString(jsonMap)

	if err != nil {
		return nil, go_cake.NewMalformedSortHTTPError(err)
	}

	converted, err := d.jsonFieldsToBSONMap(jsonMapStr, modelSpecs)

	if err != nil {
		return nil, go_cake.NewMalformedSortHTTPError(err)
	}

	sliceOfMaps := make([]map[string]int, 0)

	for k, v := range converted {
		_m := make(map[string]int)

		direction, err := v.(json.Number).Int64()

		if err != nil {
			return nil, go_cake.NewMalformedSortHTTPError(err)
		}

		_m[k] = int(direction)

		sliceOfMaps = append(sliceOfMaps, _m)
	}

	return d.getSortByMap(sliceOfMaps), nil
}

func (d *MongoDriver) getSortByMap(sliceOfMaps []map[string]int) bson.D {
	sort := bson.D{}

	for _, mapOfString := range sliceOfMaps {
		for key, value := range mapOfString {
			key = strings.TrimSpace(key)

			sort = append(sort, primitive.E{Key: key, Value: value})
		}
	}

	return sort
}
