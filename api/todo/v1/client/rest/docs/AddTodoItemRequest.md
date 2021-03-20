# AddTodoItemRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | **string** |  | 
**Order** | **int32** |  | 

## Methods

### NewAddTodoItemRequest

`func NewAddTodoItemRequest(title string, order int32, ) *AddTodoItemRequest`

NewAddTodoItemRequest instantiates a new AddTodoItemRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAddTodoItemRequestWithDefaults

`func NewAddTodoItemRequestWithDefaults() *AddTodoItemRequest`

NewAddTodoItemRequestWithDefaults instantiates a new AddTodoItemRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTitle

`func (o *AddTodoItemRequest) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *AddTodoItemRequest) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *AddTodoItemRequest) SetTitle(v string)`

SetTitle sets Title field to given value.


### GetOrder

`func (o *AddTodoItemRequest) GetOrder() int32`

GetOrder returns the Order field if non-nil, zero value otherwise.

### GetOrderOk

`func (o *AddTodoItemRequest) GetOrderOk() (*int32, bool)`

GetOrderOk returns a tuple with the Order field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrder

`func (o *AddTodoItemRequest) SetOrder(v int32)`

SetOrder sets Order field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


