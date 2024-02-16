package go_cake

type BaseGoCakeModel struct {
	gk_GetHTTPError HTTPError
}

func (bgkm *BaseGoCakeModel) GetHTTPError() HTTPError {
	return bgkm.gk_GetHTTPError
}

func (bgkm *BaseGoCakeModel) SetHTTPError(httpError HTTPError) {
	bgkm.gk_GetHTTPError = httpError
}

func (bgkm *BaseGoCakeModel) CreateInstance() GoCakeModel {
	panic("not implemented")
}

func (bgkm *BaseGoCakeModel) GetID() any {
	panic("not implemented")
}

func (bgkm *BaseGoCakeModel) SetID(id string) error {
	panic("not implemented")
}

func (bgkm *BaseGoCakeModel) CreateETag() any {
	panic("not implemented")
}

func (bgkm *BaseGoCakeModel) GetETag() any {
	panic("not implemented")
}

func (bgkm *BaseGoCakeModel) SetETag(etag string) error {
	panic("not implemented")
}
