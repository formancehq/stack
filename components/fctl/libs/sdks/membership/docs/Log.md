# Log

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Seq** | **string** |  | 
**OrganizationId** | **string** |  | 
**UserId** | **string** |  | 
**Action** | **string** |  | 
**Date** | **time.Time** |  | 
**Data** | **map[string]interface{}** |  | 

## Methods

### NewLog

`func NewLog(seq string, organizationId string, userId string, action string, date time.Time, data map[string]interface{}, ) *Log`

NewLog instantiates a new Log object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLogWithDefaults

`func NewLogWithDefaults() *Log`

NewLogWithDefaults instantiates a new Log object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSeq

`func (o *Log) GetSeq() string`

GetSeq returns the Seq field if non-nil, zero value otherwise.

### GetSeqOk

`func (o *Log) GetSeqOk() (*string, bool)`

GetSeqOk returns a tuple with the Seq field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSeq

`func (o *Log) SetSeq(v string)`

SetSeq sets Seq field to given value.


### GetOrganizationId

`func (o *Log) GetOrganizationId() string`

GetOrganizationId returns the OrganizationId field if non-nil, zero value otherwise.

### GetOrganizationIdOk

`func (o *Log) GetOrganizationIdOk() (*string, bool)`

GetOrganizationIdOk returns a tuple with the OrganizationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrganizationId

`func (o *Log) SetOrganizationId(v string)`

SetOrganizationId sets OrganizationId field to given value.


### GetUserId

`func (o *Log) GetUserId() string`

GetUserId returns the UserId field if non-nil, zero value otherwise.

### GetUserIdOk

`func (o *Log) GetUserIdOk() (*string, bool)`

GetUserIdOk returns a tuple with the UserId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserId

`func (o *Log) SetUserId(v string)`

SetUserId sets UserId field to given value.


### GetAction

`func (o *Log) GetAction() string`

GetAction returns the Action field if non-nil, zero value otherwise.

### GetActionOk

`func (o *Log) GetActionOk() (*string, bool)`

GetActionOk returns a tuple with the Action field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAction

`func (o *Log) SetAction(v string)`

SetAction sets Action field to given value.


### GetDate

`func (o *Log) GetDate() time.Time`

GetDate returns the Date field if non-nil, zero value otherwise.

### GetDateOk

`func (o *Log) GetDateOk() (*time.Time, bool)`

GetDateOk returns a tuple with the Date field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDate

`func (o *Log) SetDate(v time.Time)`

SetDate sets Date field to given value.


### GetData

`func (o *Log) GetData() map[string]interface{}`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *Log) GetDataOk() (*map[string]interface{}, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *Log) SetData(v map[string]interface{})`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


