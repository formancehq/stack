// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

import (
	"errors"
	"fmt"
	"github.com/formancehq/formance-sdk-go/v2/pkg/utils"
)

type TasksCursorDataType string

const (
	TasksCursorDataTypeTaskStripe        TasksCursorDataType = "TaskStripe"
	TasksCursorDataTypeTaskWise          TasksCursorDataType = "TaskWise"
	TasksCursorDataTypeTaskCurrencyCloud TasksCursorDataType = "TaskCurrencyCloud"
	TasksCursorDataTypeTaskDummyPay      TasksCursorDataType = "TaskDummyPay"
	TasksCursorDataTypeTaskModulr        TasksCursorDataType = "TaskModulr"
	TasksCursorDataTypeTaskBankingCircle TasksCursorDataType = "TaskBankingCircle"
	TasksCursorDataTypeTaskMangoPay      TasksCursorDataType = "TaskMangoPay"
	TasksCursorDataTypeTaskMoneycorp     TasksCursorDataType = "TaskMoneycorp"
)

type TasksCursorData struct {
	TaskStripe        *TaskStripe
	TaskWise          *TaskWise
	TaskCurrencyCloud *TaskCurrencyCloud
	TaskDummyPay      *TaskDummyPay
	TaskModulr        *TaskModulr
	TaskBankingCircle *TaskBankingCircle
	TaskMangoPay      *TaskMangoPay
	TaskMoneycorp     *TaskMoneycorp

	Type TasksCursorDataType
}

func CreateTasksCursorDataTaskStripe(taskStripe TaskStripe) TasksCursorData {
	typ := TasksCursorDataTypeTaskStripe

	return TasksCursorData{
		TaskStripe: &taskStripe,
		Type:       typ,
	}
}

func CreateTasksCursorDataTaskWise(taskWise TaskWise) TasksCursorData {
	typ := TasksCursorDataTypeTaskWise

	return TasksCursorData{
		TaskWise: &taskWise,
		Type:     typ,
	}
}

func CreateTasksCursorDataTaskCurrencyCloud(taskCurrencyCloud TaskCurrencyCloud) TasksCursorData {
	typ := TasksCursorDataTypeTaskCurrencyCloud

	return TasksCursorData{
		TaskCurrencyCloud: &taskCurrencyCloud,
		Type:              typ,
	}
}

func CreateTasksCursorDataTaskDummyPay(taskDummyPay TaskDummyPay) TasksCursorData {
	typ := TasksCursorDataTypeTaskDummyPay

	return TasksCursorData{
		TaskDummyPay: &taskDummyPay,
		Type:         typ,
	}
}

func CreateTasksCursorDataTaskModulr(taskModulr TaskModulr) TasksCursorData {
	typ := TasksCursorDataTypeTaskModulr

	return TasksCursorData{
		TaskModulr: &taskModulr,
		Type:       typ,
	}
}

func CreateTasksCursorDataTaskBankingCircle(taskBankingCircle TaskBankingCircle) TasksCursorData {
	typ := TasksCursorDataTypeTaskBankingCircle

	return TasksCursorData{
		TaskBankingCircle: &taskBankingCircle,
		Type:              typ,
	}
}

func CreateTasksCursorDataTaskMangoPay(taskMangoPay TaskMangoPay) TasksCursorData {
	typ := TasksCursorDataTypeTaskMangoPay

	return TasksCursorData{
		TaskMangoPay: &taskMangoPay,
		Type:         typ,
	}
}

func CreateTasksCursorDataTaskMoneycorp(taskMoneycorp TaskMoneycorp) TasksCursorData {
	typ := TasksCursorDataTypeTaskMoneycorp

	return TasksCursorData{
		TaskMoneycorp: &taskMoneycorp,
		Type:          typ,
	}
}

func (u *TasksCursorData) UnmarshalJSON(data []byte) error {

	var taskStripe TaskStripe = TaskStripe{}
	if err := utils.UnmarshalJSON(data, &taskStripe, "", true, true); err == nil {
		u.TaskStripe = &taskStripe
		u.Type = TasksCursorDataTypeTaskStripe
		return nil
	}

	var taskWise TaskWise = TaskWise{}
	if err := utils.UnmarshalJSON(data, &taskWise, "", true, true); err == nil {
		u.TaskWise = &taskWise
		u.Type = TasksCursorDataTypeTaskWise
		return nil
	}

	var taskCurrencyCloud TaskCurrencyCloud = TaskCurrencyCloud{}
	if err := utils.UnmarshalJSON(data, &taskCurrencyCloud, "", true, true); err == nil {
		u.TaskCurrencyCloud = &taskCurrencyCloud
		u.Type = TasksCursorDataTypeTaskCurrencyCloud
		return nil
	}

	var taskDummyPay TaskDummyPay = TaskDummyPay{}
	if err := utils.UnmarshalJSON(data, &taskDummyPay, "", true, true); err == nil {
		u.TaskDummyPay = &taskDummyPay
		u.Type = TasksCursorDataTypeTaskDummyPay
		return nil
	}

	var taskModulr TaskModulr = TaskModulr{}
	if err := utils.UnmarshalJSON(data, &taskModulr, "", true, true); err == nil {
		u.TaskModulr = &taskModulr
		u.Type = TasksCursorDataTypeTaskModulr
		return nil
	}

	var taskBankingCircle TaskBankingCircle = TaskBankingCircle{}
	if err := utils.UnmarshalJSON(data, &taskBankingCircle, "", true, true); err == nil {
		u.TaskBankingCircle = &taskBankingCircle
		u.Type = TasksCursorDataTypeTaskBankingCircle
		return nil
	}

	var taskMangoPay TaskMangoPay = TaskMangoPay{}
	if err := utils.UnmarshalJSON(data, &taskMangoPay, "", true, true); err == nil {
		u.TaskMangoPay = &taskMangoPay
		u.Type = TasksCursorDataTypeTaskMangoPay
		return nil
	}

	var taskMoneycorp TaskMoneycorp = TaskMoneycorp{}
	if err := utils.UnmarshalJSON(data, &taskMoneycorp, "", true, true); err == nil {
		u.TaskMoneycorp = &taskMoneycorp
		u.Type = TasksCursorDataTypeTaskMoneycorp
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for TasksCursorData", string(data))
}

func (u TasksCursorData) MarshalJSON() ([]byte, error) {
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

	return nil, errors.New("could not marshal union type TasksCursorData: all fields are null")
}

type TasksCursorCursor struct {
	Data     []TasksCursorData `json:"data"`
	HasMore  bool              `json:"hasMore"`
	Next     *string           `json:"next,omitempty"`
	PageSize int64             `json:"pageSize"`
	Previous *string           `json:"previous,omitempty"`
}

func (o *TasksCursorCursor) GetData() []TasksCursorData {
	if o == nil {
		return []TasksCursorData{}
	}
	return o.Data
}

func (o *TasksCursorCursor) GetHasMore() bool {
	if o == nil {
		return false
	}
	return o.HasMore
}

func (o *TasksCursorCursor) GetNext() *string {
	if o == nil {
		return nil
	}
	return o.Next
}

func (o *TasksCursorCursor) GetPageSize() int64 {
	if o == nil {
		return 0
	}
	return o.PageSize
}

func (o *TasksCursorCursor) GetPrevious() *string {
	if o == nil {
		return nil
	}
	return o.Previous
}

type TasksCursor struct {
	Cursor TasksCursorCursor `json:"cursor"`
}

func (o *TasksCursor) GetCursor() TasksCursorCursor {
	if o == nil {
		return TasksCursorCursor{}
	}
	return o.Cursor
}
