# \TodoListApi

All URIs are relative to *https://todo.api/todos*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddItem**](TodoListApi.md#AddItem) | **Post** /todos | Add a new item to the list
[**DeleteItem**](TodoListApi.md#DeleteItem) | **Delete** /todos/{id} | Delete an item
[**DeleteItems**](TodoListApi.md#DeleteItems) | **Delete** /todos | Delete all items
[**GetItem**](TodoListApi.md#GetItem) | **Get** /todos/{id} | Get an item
[**ListItems**](TodoListApi.md#ListItems) | **Get** /todos | List items
[**UpdateItem**](TodoListApi.md#UpdateItem) | **Patch** /todos/{id} | Update an existing item



## AddItem

> TodoItem AddItem(ctx, addTodoItemRequest)

Add a new item to the list

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**addTodoItemRequest** | [**AddTodoItemRequest**](AddTodoItemRequest.md)|  | 

### Return type

[**TodoItem**](TodoItem.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteItem

> DeleteItem(ctx, id)

Delete an item

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| Item ID | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DeleteItems

> DeleteItems(ctx, )

Delete all items

### Required Parameters

This endpoint does not need any parameter.

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetItem

> TodoItem GetItem(ctx, id)

Get an item

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| Item ID | 

### Return type

[**TodoItem**](TodoItem.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ListItems

> []TodoItem ListItems(ctx, )

List items

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**[]TodoItem**](TodoItem.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UpdateItem

> TodoItem UpdateItem(ctx, id, updateTodoItemRequest)

Update an existing item

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string**| Item ID | 
**updateTodoItemRequest** | [**UpdateTodoItemRequest**](UpdateTodoItemRequest.md)|  | 

### Return type

[**TodoItem**](TodoItem.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

