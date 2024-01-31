// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"errors"
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
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

	taskStripe := TaskStripe{}
	if err := utils.UnmarshalJSON(data, &taskStripe, "", true, true); err == nil {
		u.TaskStripe = &taskStripe
		u.Type = TaskResponseDataTypeTaskStripe
		return nil
	}

	taskWise := TaskWise{}
	if err := utils.UnmarshalJSON(data, &taskWise, "", true, true); err == nil {
		u.TaskWise = &taskWise
		u.Type = TaskResponseDataTypeTaskWise
		return nil
	}

	taskCurrencyCloud := TaskCurrencyCloud{}
	if err := utils.UnmarshalJSON(data, &taskCurrencyCloud, "", true, true); err == nil {
		u.TaskCurrencyCloud = &taskCurrencyCloud
		u.Type = TaskResponseDataTypeTaskCurrencyCloud
		return nil
	}

	taskDummyPay := TaskDummyPay{}
	if err := utils.UnmarshalJSON(data, &taskDummyPay, "", true, true); err == nil {
		u.TaskDummyPay = &taskDummyPay
		u.Type = TaskResponseDataTypeTaskDummyPay
		return nil
	}

	taskModulr := TaskModulr{}
	if err := utils.UnmarshalJSON(data, &taskModulr, "", true, true); err == nil {
		u.TaskModulr = &taskModulr
		u.Type = TaskResponseDataTypeTaskModulr
		return nil
	}

	taskBankingCircle := TaskBankingCircle{}
	if err := utils.UnmarshalJSON(data, &taskBankingCircle, "", true, true); err == nil {
		u.TaskBankingCircle = &taskBankingCircle
		u.Type = TaskResponseDataTypeTaskBankingCircle
		return nil
	}

	taskMangoPay := TaskMangoPay{}
	if err := utils.UnmarshalJSON(data, &taskMangoPay, "", true, true); err == nil {
		u.TaskMangoPay = &taskMangoPay
		u.Type = TaskResponseDataTypeTaskMangoPay
		return nil
	}

	taskMoneycorp := TaskMoneycorp{}
	if err := utils.UnmarshalJSON(data, &taskMoneycorp, "", true, true); err == nil {
		u.TaskMoneycorp = &taskMoneycorp
		u.Type = TaskResponseDataTypeTaskMoneycorp
		return nil
	}

	return errors.New("could not unmarshal into supported union types")
}

func (u TaskResponseData) MarshalJSON() ([]byte, error) {
	if u.TaskStripe != nil {
		return utils.MarshalJSON(u.TaskStripe, "", true)
	}

	if u.TaskWise != nil {
		return utils.MarshalJSON(u.TaskWise, "", true)
	}

	if u.TaskCurrencyCloud != nil {
		return utils.MarshalJSON(u.TaskCurrencyCloud, "", true)
	}

	if u.TaskDummyPay != nil {
		return utils.MarshalJSON(u.TaskDummyPay, "", true)
	}

	if u.TaskModulr != nil {
		return utils.MarshalJSON(u.TaskModulr, "", true)
	}

	if u.TaskBankingCircle != nil {
		return utils.MarshalJSON(u.TaskBankingCircle, "", true)
	}

	if u.TaskMangoPay != nil {
		return utils.MarshalJSON(u.TaskMangoPay, "", true)
	}

	if u.TaskMoneycorp != nil {
		return utils.MarshalJSON(u.TaskMoneycorp, "", true)
	}

	return nil, errors.New("could not marshal union type: all fields are null")
}

type TaskResponse struct {
	Data TaskResponseData `json:"data"`
}

func (o *TaskResponse) GetData() TaskResponseData {
	if o == nil {
		return TaskResponseData{}
	}
	return o.Data
}