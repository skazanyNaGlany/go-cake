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
