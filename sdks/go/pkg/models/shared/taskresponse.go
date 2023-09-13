// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"bytes"
	"encoding/json"
	"errors"
)

type TaskResponseDataType string

const (
	TaskResponseDataTypeTaskStripe        TaskResponseDataType = "TaskStripe"
	TaskResponseDataTypeTaskWise          TaskResponseDataType = "TaskWise"
	TaskResponseDataTypeTaskCurrencyCloud TaskResponseDataType = "TaskCurrencyCloud"
	TaskResponseDataTypeTaskDummyPay      TaskResponseDataType = "TaskDummyPay"
	TaskResponseDataTypeTaskModulr        TaskResponseDataType = "TaskModulr"
	TaskResponseDataTypeTaskBankingCircle TaskResponseDataType = "TaskBankingCircle"
	TaskResponseDataTypeTaskMangoPay      TaskResponseDataType = "TaskMangoPay"
	TaskResponseDataTypeTaskMoneycorp     TaskResponseDataType = "TaskMoneycorp"
)

type TaskResponseData struct {
	TaskStripe        *TaskStripe
	TaskWise          *TaskWise
	TaskCurrencyCloud *TaskCurrencyCloud
	TaskDummyPay      *TaskDummyPay
	TaskModulr        *TaskModulr
	TaskBankingCircle *TaskBankingCircle
	TaskMangoPay      *TaskMangoPay
	TaskMoneycorp     *TaskMoneycorp

	Type TaskResponseDataType
}

func CreateTaskResponseDataTaskStripe(taskStripe TaskStripe) TaskResponseData {
	typ := TaskResponseDataTypeTaskStripe

	return TaskResponseData{
		TaskStripe: &taskStripe,
		Type:       typ,
	}
}

func CreateTaskResponseDataTaskWise(taskWise TaskWise) TaskResponseData {
	typ := TaskResponseDataTypeTaskWise

	return TaskResponseData{
		TaskWise: &taskWise,
		Type:     typ,
	}
}

func CreateTaskResponseDataTaskCurrencyCloud(taskCurrencyCloud TaskCurrencyCloud) TaskResponseData {
	typ := TaskResponseDataTypeTaskCurrencyCloud

	return TaskResponseData{
		TaskCurrencyCloud: &taskCurrencyCloud,
		Type:              typ,
	}
}

func CreateTaskResponseDataTaskDummyPay(taskDummyPay TaskDummyPay) TaskResponseData {
	typ := TaskResponseDataTypeTaskDummyPay

	return TaskResponseData{
		TaskDummyPay: &taskDummyPay,
		Type:         typ,
	}
}

func CreateTaskResponseDataTaskModulr(taskModulr TaskModulr) TaskResponseData {
	typ := TaskResponseDataTypeTaskModulr

	return TaskResponseData{
		TaskModulr: &taskModulr,
		Type:       typ,
	}
}

func CreateTaskResponseDataTaskBankingCircle(taskBankingCircle TaskBankingCircle) TaskResponseData {
	typ := TaskResponseDataTypeTaskBankingCircle

	return TaskResponseData{
		TaskBankingCircle: &taskBankingCircle,
		Type:              typ,
	}
}

func CreateTaskResponseDataTaskMangoPay(taskMangoPay TaskMangoPay) TaskResponseData {
	typ := TaskResponseDataTypeTaskMangoPay

	return TaskResponseData{
		TaskMangoPay: &taskMangoPay,
		Type:         typ,
	}
}

func CreateTaskResponseDataTaskMoneycorp(taskMoneycorp TaskMoneycorp) TaskResponseData {
	typ := TaskResponseDataTypeTaskMoneycorp

	return TaskResponseData{
		TaskMoneycorp: &taskMoneycorp,
		Type:          typ,
	}
}

func (u *TaskResponseData) UnmarshalJSON(data []byte) error {
	var d *json.Decoder

	taskStripe := new(TaskStripe)
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&taskStripe); err == nil {
		u.TaskStripe = taskStripe
		u.Type = TaskResponseDataTypeTaskStripe
		return nil
	}

	taskWise := new(TaskWise)
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&taskWise); err == nil {
		u.TaskWise = taskWise
		u.Type = TaskResponseDataTypeTaskWise
		return nil
	}

	taskCurrencyCloud := new(TaskCurrencyCloud)
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&taskCurrencyCloud); err == nil {
		u.TaskCurrencyCloud = taskCurrencyCloud
		u.Type = TaskResponseDataTypeTaskCurrencyCloud
		return nil
	}

	taskDummyPay := new(TaskDummyPay)
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&taskDummyPay); err == nil {
		u.TaskDummyPay = taskDummyPay
		u.Type = TaskResponseDataTypeTaskDummyPay
		return nil
	}

	taskModulr := new(TaskModulr)
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&taskModulr); err == nil {
		u.TaskModulr = taskModulr
		u.Type = TaskResponseDataTypeTaskModulr
		return nil
	}

	taskBankingCircle := new(TaskBankingCircle)
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&taskBankingCircle); err == nil {
		u.TaskBankingCircle = taskBankingCircle
		u.Type = TaskResponseDataTypeTaskBankingCircle
		return nil
	}

	taskMangoPay := new(TaskMangoPay)
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&taskMangoPay); err == nil {
		u.TaskMangoPay = taskMangoPay
		u.Type = TaskResponseDataTypeTaskMangoPay
		return nil
	}

	taskMoneycorp := new(TaskMoneycorp)
	d = json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(&taskMoneycorp); err == nil {
		u.TaskMoneycorp = taskMoneycorp
		u.Type = TaskResponseDataTypeTaskMoneycorp
		return nil
	}

	return errors.New("could not unmarshal into supported union types")
}

func (u TaskResponseData) MarshalJSON() ([]byte, error) {
	if u.TaskStripe != nil {
		return json.Marshal(u.TaskStripe)
	}

	if u.TaskWise != nil {
		return json.Marshal(u.TaskWise)
	}

	if u.TaskCurrencyCloud != nil {
		return json.Marshal(u.TaskCurrencyCloud)
	}

	if u.TaskDummyPay != nil {
		return json.Marshal(u.TaskDummyPay)
	}

	if u.TaskModulr != nil {
		return json.Marshal(u.TaskModulr)
	}

	if u.TaskBankingCircle != nil {
		return json.Marshal(u.TaskBankingCircle)
	}

	if u.TaskMangoPay != nil {
		return json.Marshal(u.TaskMangoPay)
	}

	if u.TaskMoneycorp != nil {
		return json.Marshal(u.TaskMoneycorp)
	}

	return nil, nil
}

// TaskResponse - OK
type TaskResponse struct {
	Data TaskResponseData `json:"data"`
}
