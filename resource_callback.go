package go_cake

type ResourceCallback struct {
	AuthCallback        AuthAppFunc
	PreRequestCallback  PrePostRequestAppFunc
	PostRequestCallback PrePostRequestAppFunc
	FetchedDocuments    DocumentsAppFunc
	UpdatingDocuments   DocumentsAppFunc
	UpdatedDocuments    DocumentsAppFunc
	InsertingDocuments  DocumentsAppFunc
	InsertedDocuments   DocumentsAppFunc
	DeletingDocuments   DocumentsAppFunc
	DeletedDocuments    DocumentsAppFunc
}
