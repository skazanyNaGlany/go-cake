package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	go_cake "github.com/skazanyNaGlany/go-cake"
	"github.com/skazanyNaGlany/go-cake/driver/postgres"
	"github.com/skazanyNaGlany/go-cake/examples/postgres/models"
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

	// dbDriver, err := postgres.NewPostgresDriver(
	// 	"host=db.devbox user=postgres password= dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai",
	// 	ctx)

	// if err != nil {
	// 	log.Panicln("Unable to init DB driver:", err)
	// }

	// defer dbDriver.Close()

	// BUN
	dbDriver, err := postgres.NewPostgresDriver(
		"postgres://postgres:@db.devbox:5432/postgres?sslmode=disable",
		ctx)

	if err != nil {
		log.Panicln("Unable to init DB driver:", err)
	}

	defer dbDriver.Close()

	restHandler := go_cake.NewHandler()

	devicesResource, err := go_cake.NewResource(
		`^(?P<version>\/\w+)(?P<url>\/api\/devices2\/?)$`,
		"public.device2",
		"devices2",
		dbDriver,
		&models.Device2{},
		"ID",
		"id",
		"",
		"",
		[]string{"v1"},
		checkAuth)

	if err != nil {
		panic(err)
	}

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
	/////////////////

	// usersResource, err := go_cake.NewResource(
	// 	`^(?P<version>\/\w+)(?P<url>\/api\/users2\/?)$`,
	// 	"public.user2",
	// 	"users2",
	// 	dbDriver,
	// 	&models.User2{},
	// 	"ID",
	// 	"id",
	// 	"",
	// 	"",
	// 	[]string{"v1"},
	// 	checkAuth)

	// if err != nil {
	// 	panic(err)
	// }

	// usersResource.ResourceCallback.PreRequestCallback = preRequest
	// usersResource.ResourceCallback.PostRequestCallback = postRequest
	// usersResource.ResourceCallback.FetchedDocuments = fetchedDocuments
	// usersResource.ResourceCallback.UpdatingDocuments = updatingDocuments
	// usersResource.ResourceCallback.UpdatedDocuments = updatedDocuments
	// usersResource.ResourceCallback.InsertingDocuments = insertingDocuments
	// usersResource.ResourceCallback.InsertedDocuments = insertedDocuments
	// usersResource.ResourceCallback.DeletingDocuments = deletingDocuments
	// usersResource.ResourceCallback.DeletedDocuments = deletedDocuments

	// if err := restHandler.AddResource(usersResource); err != nil {
	// 	panic(err)
	// }

	log.Println("Ready")

	http.ListenAndServe(":8080", restHandler)
}
