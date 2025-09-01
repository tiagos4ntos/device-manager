# Device Manager API

API for managing an inventory of mobile devices.

**Version:** 0.0.1-beta

## Endpoints

### `GET /devices`

*List devices*

Get all devices

#### Responses

| Status Code | Description | Schema |
|-------------|-------------|--------|
| 200 | OK | - |
| 500 | Internal Server Error | - |

### `POST /devices`

*Create a new device*

Registers a new device on the database with the provided information

#### Parameters

| Name | In | Required | Description | Type |
|------|----|----------|-------------|------|
| `device` | body | Yes | Device payload | - |

#### Responses

| Status Code | Description | Schema |
|-------------|-------------|--------|
| 201 | Created | - |
| 400 | Bad Request | - |
| 500 | Internal Server Error | - |

### `GET /devices/{id}`

*Get device by ID*

Returns a single device by its ID

#### Parameters

| Name | In | Required | Description | Type |
|------|----|----------|-------------|------|
| `id` | path | Yes | Device ID | - |

#### Responses

| Status Code | Description | Schema |
|-------------|-------------|--------|
| 200 | OK | - |
| 400 | Bad Request | - |
| 404 | Not Found | - |
| 500 | Internal Server Error | - |

### `PUT /devices/{id}`

*Updates device data by ID*

Update an existing device name and brand but only when the state is not "in-use", the fiel state can be updated anytime

#### Parameters

| Name | In | Required | Description | Type |
|------|----|----------|-------------|------|
| `id` | path | Yes | Device ID | - |
| `device` | body | Yes | Updated device payload | - |

#### Responses

| Status Code | Description | Schema |
|-------------|-------------|--------|
| 200 | OK | - |
| 400 | Bad Request | - |
| 500 | Internal Server Error | - |

### `DELETE /devices/{id}`

*Delete a device*

Removes a device by ID, only devices that are not "in-use" can be deleted

#### Parameters

| Name | In | Required | Description | Type |
|------|----|----------|-------------|------|
| `id` | path | Yes | Device ID | - |

#### Responses

| Status Code | Description | Schema |
|-------------|-------------|--------|
| 204 | No Content | - |
| 404 | Not Found | - |
| 500 | Internal Server Error | - |