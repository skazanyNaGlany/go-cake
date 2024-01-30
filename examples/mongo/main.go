package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	go_cake "github.com/skazanyNaGlany/go-cake"
	driver "github.com/skazanyNaGlany/go-cake/driver/mongo"
	models "github.com/skazanyNaGlany/go-cake/examples/mongo/models"
)

func checkAuth(
	resource *go_cake.Resource,
	request *go_cake.Request,
	response *go_cake.ResponseJSON) bool {
	return true
}

func preRequest(
	resource *go_cake.Resource,
	request *go_cake.Request,
	response *go_cake.ResponseJSON) go_cake.HTTPError {
	log.Println(">>>", request.UniqueID, request.Method, request.URL)

	return nil
}

func postRequest(
	resource *go_cake.Resource,
	request *go_cake.Request,
	response *go_cake.ResponseJSON) go_cake.HTTPError {
	log.Println("<<<", request.UniqueID, request.Method, request.URL, "["+strconv.FormatInt(int64(response.Meta.StatusCode), 10)+"]")

	return nil
}

func printDocuments(action string, documents []go_cake.GoKateModel) {
	for _, doc := range documents {
		log.Println(action, doc)
	}
}

func fetchedDocuments(
	resource *go_cake.Resource,
	request *go_cake.Request,
	documents []go_cake.GoKateModel,
	currentHttpErr go_cake.HTTPError) go_cake.HTTPError {
	if currentHttpErr != nil {
		return currentHttpErr
	}

	printDocuments("fetched documents", documents)

	return nil
}

func updatingDocuments(
	resource *go_cake.Resource,
	request *go_cake.Request,
	documents []go_cake.GoKateModel,
	currentHttpErr go_cake.HTTPError) go_cake.HTTPError {
	if currentHttpErr != nil {
		return currentHttpErr
	}

	// product := documents[0].(*models.Product)

	printDocuments("updating documents", documents)

	return nil
}

func updatedDocuments(
	resource *go_cake.Resource,
	request *go_cake.Request,
	documents []go_cake.GoKateModel,
	currentHttpErr go_cake.HTTPError) go_cake.HTTPError {
	if currentHttpErr != nil {
		return currentHttpErr
	}

	printDocuments("updated documents", documents)

	return nil
}

func insertingDocuments(
	resource *go_cake.Resource,
	request *go_cake.Request,
	documents []go_cake.GoKateModel,
	currentHttpErr go_cake.HTTPError) go_cake.HTTPError {
	if currentHttpErr != nil {
		return currentHttpErr
	}

	printDocuments("inserting documents", documents)

	return nil
}

func insertedDocuments(
	resource *go_cake.Resource,
	request *go_cake.Request,
	documents []go_cake.GoKateModel,
	currentHttpErr go_cake.HTTPError) go_cake.HTTPError {
	if currentHttpErr != nil {
		return currentHttpErr
	}

	printDocuments("inserted documents", documents)

	return nil
}

func deletingDocuments(
	resource *go_cake.Resource,
	request *go_cake.Request,
	documents []go_cake.GoKateModel,
	currentHttpErr go_cake.HTTPError) go_cake.HTTPError {
	if currentHttpErr != nil {
		return currentHttpErr
	}

	printDocuments("deleting documents", documents)

	return nil
}

func deletedDocuments(
	resource *go_cake.Resource,
	request *go_cake.Request,
	documents []go_cake.GoKateModel,
	currentHttpErr go_cake.HTTPError) go_cake.HTTPError {
	if currentHttpErr != nil {
		return currentHttpErr
	}

	printDocuments("deleted documents", documents)

	return nil
}

type NotFoundHandler struct{}

func (h *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("NotFoundHandler.ServeHTTP called")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dbDriver, err := driver.NewMongoDriver("mongodb://db.devbox:27017/", "local", ctx)

	if err != nil {
		log.Panicln("Unable to init DB driver:", err)
	}

	defer dbDriver.Close()

	corsConfig, _ := go_cake.NewDefaultCORSConfig()
	restHandler := go_cake.NewHandler()

	ordersValidator, err := go_cake.NewDefaultJSONValidator("orders.json", `{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"type": "object",
		"properties": {
		  "required_field": {
			"type": "string"
		  },
		  "projected_field": {
			"type": "string"
		  },
		  "sorted_field": {
			"type": "string"
		  },
		  "updated_field": {
			"type": "string"
		  },
		  "hidden_proj_field": {
			"type": "string"
		  },
		  "hidden_field": {
			"type": "string"
		  },
		  "erased_field": {
			"type": "string"
		  },
		  "required_field2": {
			"type": "string"
		  },
		  "filtered_field": {
			"type": "string"
		  },
		  "inserted_field": {
			"type": "string"
		  }
		},
		"required": [
		  "required_field",
		  "projected_field",
		  "sorted_field",
		  "updated_field",
		  "hidden_proj_field",
		  "hidden_field",
		  "erased_field",
		  "required_field2",
		  "filtered_field",
		  "inserted_field"
		]
	  }`)

	if err != nil {
		log.Panicln("Unable to create JSON validator:", err)
	}

	ordersResource := go_cake.Resource{
		Pattern:             `^(?P<version>\/\w+)(?P<url>\/api\/orders\/?)$`,
		DbPath:              "orders",
		ResourceName:        "orders",
		DatabaseDriver:      dbDriver,
		DbModel:             &models.Order{},
		DbModelIDField:      "ID",
		DbModelETagField:    "ETag",
		SupportedVersion:    []string{"v1"},
		GetAllowed:          true,
		DeleteAllowed:       true,
		InsertAllowed:       true,
		UpdateAllowed:       true,
		GetMaxOutputItems:   100,
		InsertMaxInputItems: 100,
		DeleteMaxInputItems: 100,
		UpdateMaxInputItems: 100,
		CORSConfig:          corsConfig,
		JSONSchemaConfig: &go_cake.JSONSchemaConfig{
			IDField:                "id",
			ETagField:              "_etag",
			FilterableFields:       []string{"filtered_field"},
			ProjectableFields:      []string{"hidden_proj_field", "projected_field"},
			SortableFields:         []string{"sorted_field", "sorted_field2"},
			InsertableFields:       []string{},
			UpdatableFields:        []string{},
			HiddenFields:           []string{"hidden_field"},
			ErasedFields:           []string{},
			RequiredOnInsertFields: []string{go_cake.FIELD_ANY},
			RequiredOnUpdateFields: []string{go_cake.FIELD_ANY},
			OptimizeOnInsertFields: []string{go_cake.FIELD_ANY},
			OptimizeOnUpdateFields: []string{go_cake.FIELD_ANY},
			Validator:              ordersValidator,
		},
		ResourceCallback: &go_cake.ResourceCallback{
			AuthCallback:        checkAuth,
			PreRequestCallback:  preRequest,
			PostRequestCallback: postRequest,
			FetchedDocuments:    fetchedDocuments,
			UpdatingDocuments:   updatingDocuments,
			UpdatedDocuments:    updatedDocuments,
			InsertingDocuments:  insertingDocuments,
			InsertedDocuments:   insertedDocuments,
			DeletingDocuments:   deletingDocuments,
			DeletedDocuments:    deletedDocuments,
		},
	}

	if err := ordersResource.Init(); err != nil {
		panic(err)
	}

	if err := restHandler.AddResource(&ordersResource); err != nil {
		panic(err)
	}

	usersValidator, err := go_cake.NewDefaultJSONValidator("users.json", `{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"type": "object",
		"properties": {
		  "email": {
			"type": "string"
		  },
		  "max_contacts": {
			"type": "number",
			"minimum": 1,
			"maximum": 9223372036854775807
		  }
		},
		"required": [
		  "email",
		  "max_contacts"
		]
	  }`)

	if err != nil {
		log.Panicln("Unable to create JSON validator:", err)
	}

	usersResource := go_cake.Resource{
		Pattern:             `^(?P<version>\/\w+)(?P<url>\/api\/users\/?)$`,
		DbPath:              "users",
		ResourceName:        "users",
		DatabaseDriver:      dbDriver,
		DbModel:             &models.User{},
		DbModelIDField:      "ID",
		DbModelETagField:    "ETag",
		SupportedVersion:    []string{"v1"},
		GetAllowed:          true,
		DeleteAllowed:       true,
		InsertAllowed:       true,
		UpdateAllowed:       true,
		GetMaxOutputItems:   100,
		InsertMaxInputItems: 100,
		DeleteMaxInputItems: 100,
		UpdateMaxInputItems: 100,
		CORSConfig:          corsConfig,
		JSONSchemaConfig: &go_cake.JSONSchemaConfig{
			IDField:                "id",
			ETagField:              "_etag",
			FilterableFields:       []string{go_cake.FIELD_ANY},
			ProjectableFields:      []string{go_cake.FIELD_ANY},
			SortableFields:         []string{go_cake.FIELD_ANY},
			InsertableFields:       []string{go_cake.FIELD_ANY},
			UpdatableFields:        []string{go_cake.FIELD_ANY},
			HiddenFields:           []string{},
			ErasedFields:           []string{},
			RequiredOnInsertFields: []string{go_cake.FIELD_ANY},
			RequiredOnUpdateFields: []string{go_cake.FIELD_ANY},
			OptimizeOnInsertFields: []string{go_cake.FIELD_ANY},
			OptimizeOnUpdateFields: []string{go_cake.FIELD_ANY},
			Validator:              usersValidator,
		},
		ResourceCallback: &go_cake.ResourceCallback{
			AuthCallback:        checkAuth,
			PreRequestCallback:  preRequest,
			PostRequestCallback: postRequest,
			FetchedDocuments:    fetchedDocuments,
			UpdatingDocuments:   updatingDocuments,
			UpdatedDocuments:    updatedDocuments,
			InsertingDocuments:  insertingDocuments,
			InsertedDocuments:   insertedDocuments,
			DeletingDocuments:   deletingDocuments,
			DeletedDocuments:    deletedDocuments,
		},
	}

	if err := usersResource.Init(); err != nil {
		panic(err)
	}

	if err := restHandler.AddResource(&usersResource); err != nil {
		panic(err)
	}

	productsValidator, err := go_cake.NewDefaultJSONValidator("products.json", `{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"type": "object",
		"properties": {
		  "email": {
			"type": "string"
		  },
		  "max_contacts": {
			"type": "number",
			"minimum": 1,
			"maximum": 9223372036854775807
		  }
		}
	  }
	`)

	if err != nil {
		log.Panicln("Unable to create JSON validator:", err)
	}

	productsResource := go_cake.Resource{
		Pattern:             `^(?P<version>\/\w+)(?P<url>\/api\/products\/?)$`,
		DbPath:              "products",
		ResourceName:        "products",
		DatabaseDriver:      dbDriver,
		DbModel:             &models.Product{},
		DbModelIDField:      "ID",
		DbModelETagField:    "ETag",
		SupportedVersion:    []string{"v1"},
		GetAllowed:          true,
		DeleteAllowed:       true,
		InsertAllowed:       true,
		UpdateAllowed:       true,
		GetMaxOutputItems:   100,
		InsertMaxInputItems: 100,
		DeleteMaxInputItems: 100,
		UpdateMaxInputItems: 100,
		CORSConfig:          corsConfig,
		JSONSchemaConfig: &go_cake.JSONSchemaConfig{
			IDField:           "id",
			ETagField:         "_etag",
			FilterableFields:  []string{go_cake.FIELD_ANY},
			ProjectableFields: []string{go_cake.FIELD_ANY},
			SortableFields:    []string{go_cake.FIELD_ANY},
			// InsertableFields:  []string{},
			InsertableFields: []string{go_cake.FIELD_ANY},
			UpdatableFields:  []string{"email"},
			// UpdatableFields:  []string{handler.FIELD_ANY},
			HiddenFields: []string{},
			ErasedFields: []string{},
			// RequiredOnInsertFields: []string{},
			RequiredOnInsertFields: []string{go_cake.FIELD_ANY},
			RequiredOnUpdateFields: []string{"email"},
			// RequiredOnUpdateFields: []string{handler.FIELD_ANY},
			OptimizeOnInsertFields: []string{go_cake.FIELD_ANY},
			OptimizeOnUpdateFields: []string{go_cake.FIELD_ANY},
			Validator:              productsValidator,
		},
		ResourceCallback: &go_cake.ResourceCallback{
			AuthCallback:        checkAuth,
			PreRequestCallback:  preRequest,
			PostRequestCallback: postRequest,
			FetchedDocuments:    fetchedDocuments,
			UpdatingDocuments:   updatingDocuments,
			UpdatedDocuments:    updatedDocuments,
			InsertingDocuments:  insertingDocuments,
			InsertedDocuments:   insertedDocuments,
			DeletingDocuments:   deletingDocuments,
			DeletedDocuments:    deletedDocuments,
		},
	}

	if err := productsResource.Init(); err != nil {
		panic(err)
	}

	if err := restHandler.AddResource(&productsResource); err != nil {
		panic(err)
	}

	devicesValidator, err := go_cake.NewDefaultJSONValidator("devices.json", `{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"type": "object",
		"properties": {
		  "email": {
			"type": "string"
		  },
		  "max_contacts": {
			"type": "number",
			"minimum": 1,
			"maximum": 9223372036854775807
		  }
		}
	  }
	`)

	if err != nil {
		log.Panicln("Unable to create JSON validator:", err)
	}

	devicesResource := go_cake.Resource{
		Pattern:             `^(?P<version>\/\w+)(?P<url>\/api\/devices\/?)$`,
		DbPath:              "devices",
		ResourceName:        "devices",
		DatabaseDriver:      dbDriver,
		DbModel:             &models.Device{},
		DbModelIDField:      "ID",
		SupportedVersion:    []string{"v1"},
		GetAllowed:          true,
		DeleteAllowed:       true,
		InsertAllowed:       true,
		UpdateAllowed:       true,
		GetMaxOutputItems:   100,
		InsertMaxInputItems: 100,
		DeleteMaxInputItems: 100,
		UpdateMaxInputItems: 100,
		CORSConfig:          corsConfig,
		JSONSchemaConfig: &go_cake.JSONSchemaConfig{
			IDField:           "id",
			FilterableFields:  []string{go_cake.FIELD_ANY},
			ProjectableFields: []string{go_cake.FIELD_ANY},
			SortableFields:    []string{go_cake.FIELD_ANY},
			// InsertableFields:  []string{},
			InsertableFields: []string{go_cake.FIELD_ANY},
			UpdatableFields:  []string{"email"},
			// UpdatableFields:  []string{handler.FIELD_ANY},
			HiddenFields: []string{},
			ErasedFields: []string{},
			// RequiredOnInsertFields: []string{},
			RequiredOnInsertFields: []string{go_cake.FIELD_ANY},
			RequiredOnUpdateFields: []string{"email"},
			// RequiredOnUpdateFields: []string{handler.FIELD_ANY},
			OptimizeOnInsertFields: []string{go_cake.FIELD_ANY},
			OptimizeOnUpdateFields: []string{go_cake.FIELD_ANY},
			Validator:              devicesValidator,
		},
		ResourceCallback: &go_cake.ResourceCallback{
			AuthCallback:        checkAuth,
			PreRequestCallback:  preRequest,
			PostRequestCallback: postRequest,
			FetchedDocuments:    fetchedDocuments,
			UpdatingDocuments:   updatingDocuments,
			UpdatedDocuments:    updatedDocuments,
			InsertingDocuments:  insertingDocuments,
			InsertedDocuments:   insertedDocuments,
			DeletingDocuments:   deletingDocuments,
			DeletedDocuments:    deletedDocuments,
		},
	}

	if err := devicesResource.Init(); err != nil {
		log.Printf("%T\n", err)
		panic(err)
	}

	if err := restHandler.AddResource(&devicesResource); err != nil {
		panic(err)
	}

	log.Println("Ready")

	http.ListenAndServe(":8080", restHandler)
}
