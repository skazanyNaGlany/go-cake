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

func printDocuments(action string, documents []go_cake.GoCakeModel) {
	for _, doc := range documents {
		log.Println(action, doc)
	}
}

func fetchedDocuments(
	resource *go_cake.Resource,
	request *go_cake.Request,
	documents []go_cake.GoCakeModel,
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
	documents []go_cake.GoCakeModel,
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
	documents []go_cake.GoCakeModel,
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
	documents []go_cake.GoCakeModel,
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
	documents []go_cake.GoCakeModel,
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
	documents []go_cake.GoCakeModel,
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
	documents []go_cake.GoCakeModel,
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

	ordersResource, err := go_cake.NewResource(
		`^(?P<version>\/\w+)(?P<url>\/api\/orders\/?)$`,
		"orders",
		"orders",
		dbDriver,
		&models.Order{},
		"ID",
		"id",
		"ETag",
		"_etag",
		[]string{"v1"},
		checkAuth)

	if err != nil {
		panic(err)
	}

	defer ordersResource.Close()

	ordersResource.JSONSchemaConfig.Validator = ordersValidator
	ordersResource.ResourceCallback.PreRequestCallback = preRequest
	ordersResource.ResourceCallback.PostRequestCallback = postRequest
	ordersResource.ResourceCallback.FetchedDocuments = fetchedDocuments
	ordersResource.ResourceCallback.UpdatingDocuments = updatingDocuments
	ordersResource.ResourceCallback.UpdatedDocuments = updatedDocuments
	ordersResource.ResourceCallback.InsertingDocuments = insertingDocuments
	ordersResource.ResourceCallback.InsertedDocuments = insertedDocuments
	ordersResource.ResourceCallback.DeletingDocuments = deletingDocuments
	ordersResource.ResourceCallback.DeletedDocuments = deletedDocuments

	if err := restHandler.AddResource(ordersResource); err != nil {
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

	usersResource, err := go_cake.NewResource(
		`^(?P<version>\/\w+)(?P<url>\/api\/users\/?)$`,
		"users",
		"users",
		dbDriver,
		&models.User{},
		"ID",
		"id",
		"ETag",
		"_etag",
		[]string{"v1"},
		checkAuth)

	if err != nil {
		panic(err)
	}

	defer usersResource.Close()

	usersResource.JSONSchemaConfig.Validator = usersValidator
	usersResource.ResourceCallback.PreRequestCallback = preRequest
	usersResource.ResourceCallback.PostRequestCallback = postRequest
	usersResource.ResourceCallback.FetchedDocuments = fetchedDocuments
	usersResource.ResourceCallback.UpdatingDocuments = updatingDocuments
	usersResource.ResourceCallback.UpdatedDocuments = updatedDocuments
	usersResource.ResourceCallback.InsertingDocuments = insertingDocuments
	usersResource.ResourceCallback.InsertedDocuments = insertedDocuments
	usersResource.ResourceCallback.DeletingDocuments = deletingDocuments
	usersResource.ResourceCallback.DeletedDocuments = deletedDocuments

	if err := restHandler.AddResource(usersResource); err != nil {
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

	productsResource, err := go_cake.NewResource(
		`^(?P<version>\/\w+)(?P<url>\/api\/products\/?)$`,
		"products",
		"products",
		dbDriver,
		&models.Product{},
		"ID",
		"id",
		"ETag",
		"_etag",
		[]string{"v1"},
		checkAuth)

	if err != nil {
		panic(err)
	}

	defer productsResource.Close()

	productsResource.JSONSchemaConfig.Validator = productsValidator
	productsResource.ResourceCallback.PreRequestCallback = preRequest
	productsResource.ResourceCallback.PostRequestCallback = postRequest
	productsResource.ResourceCallback.FetchedDocuments = fetchedDocuments
	productsResource.ResourceCallback.UpdatingDocuments = updatingDocuments
	productsResource.ResourceCallback.UpdatedDocuments = updatedDocuments
	productsResource.ResourceCallback.InsertingDocuments = insertingDocuments
	productsResource.ResourceCallback.InsertedDocuments = insertedDocuments
	productsResource.ResourceCallback.DeletingDocuments = deletingDocuments
	productsResource.ResourceCallback.DeletedDocuments = deletedDocuments

	if err := restHandler.AddResource(productsResource); err != nil {
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

	devicesResource, err := go_cake.NewResource(
		`^(?P<version>\/\w+)(?P<url>\/api\/devices\/?)$`,
		"devices",
		"devices",
		dbDriver,
		&models.Device{},
		"ID",
		"",
		"ETag",
		"",
		[]string{"v1"},
		checkAuth)

	if err != nil {
		panic(err)
	}

	defer devicesResource.Close()

	devicesResource.JSONSchemaConfig.Validator = devicesValidator
	devicesResource.ResourceCallback.PreRequestCallback = preRequest
	devicesResource.ResourceCallback.PostRequestCallback = postRequest
	devicesResource.ResourceCallback.FetchedDocuments = fetchedDocuments
	devicesResource.ResourceCallback.UpdatingDocuments = updatingDocuments
	devicesResource.ResourceCallback.UpdatedDocuments = updatedDocuments
	devicesResource.ResourceCallback.InsertingDocuments = insertingDocuments
	devicesResource.ResourceCallback.InsertedDocuments = insertedDocuments
	devicesResource.ResourceCallback.DeletingDocuments = deletingDocuments
	devicesResource.ResourceCallback.DeletedDocuments = deletedDocuments

	if err := restHandler.AddResource(devicesResource); err != nil {
		panic(err)
	}

	log.Println("Ready")

	http.ListenAndServe(":8080", restHandler)
}
