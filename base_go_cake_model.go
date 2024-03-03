package go_cake

import "encoding/json"

type BaseGoCakeModel struct {
	httpError HTTPError
	subModel  GoCakeModel
}

func (bgkm *BaseGoCakeModel) SetSubModel(model GoCakeModel) {
	bgkm.subModel = model
}

func (bgkm *BaseGoCakeModel) GetHTTPError() HTTPError {
	return bgkm.httpError
}

func (bgkm *BaseGoCakeModel) SetHTTPError(httpError HTTPError) {
	bgkm.httpError = httpError
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

func (bgkm *BaseGoCakeModel) ToMap() (map[string]any, error) {
	itemBytes, _ := json.Marshal(bgkm.subModel)
	objectMap := make(map[string]any)

	if err := json.Unmarshal(itemBytes, &objectMap); err != nil {
		return nil, err
	}

	if httpErr := bgkm.GetHTTPError(); httpErr != nil {
		_meta := make(map[string]any)

		_meta["status_code"] = httpErr.GetStatusCode()
		_meta["status_message"] = httpErr.GetStatusMessage()

		objectMap["_meta"] = _meta
	}

	return objectMap, nil
}
