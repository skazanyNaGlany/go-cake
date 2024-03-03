package go_cake

type ResourceCallback struct {
	AuthCallback        AuthCallback
	PreRequestCallback  PrePostRequestCallback
	PostRequestCallback PrePostRequestCallback
	FetchedDocuments    DocumentsCallback
	UpdatingDocuments   DocumentsCallback
	UpdatedDocuments    DocumentsCallback
	InsertingDocuments  DocumentsCallback
	InsertedDocuments   DocumentsCallback
	DeletingDocuments   DocumentsCallback
	DeletedDocuments    DocumentsCallback
	CreateContext       CreateContextCallback
}
