package go_cake

type BaseGoKateModel struct {
	gk_GetHTTPError HTTPError
}

func (bgkm *BaseGoKateModel) GetHTTPError() HTTPError {
	return bgkm.gk_GetHTTPError
}

func (bgkm *BaseGoKateModel) SetHTTPError(httpError HTTPError) {
	bgkm.gk_GetHTTPError = httpError
}

func (bgkm *BaseGoKateModel) CreateInstance() GoKateModel {
	panic("not implemented")
}

func (bgkm *BaseGoKateModel) GetID() any {
	panic("not implemented")
}

func (bgkm *BaseGoKateModel) SetID(id string) error {
	panic("not implemented")
}

func (bgkm *BaseGoKateModel) CreateETag() any {
	panic("not implemented")
}

func (bgkm *BaseGoKateModel) GetETag() any {
	panic("not implemented")
}

func (bgkm *BaseGoKateModel) SetETag(etag string) error {
	panic("not implemented")
}
