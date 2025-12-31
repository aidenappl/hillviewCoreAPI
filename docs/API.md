# Hillview Core API Documentation

## Overview

The Hillview Core API is a RESTful service that powers the Hillview TV platform. It provides endpoints for managing videos, assets, playlists, links, users, checkouts, and spotlight content.

**Base URL:** `/core/v1.1`

**API Version:** 1.1

---

## Table of Contents

- [Authentication](#authentication)
- [Common Response Format](#common-response-format)
- [Error Handling](#error-handling)
- [Data Models](#data-models)
- [Endpoints](#endpoints)
  - [Health Check](#health-check)
  - [Videos](#videos)
  - [Assets](#assets)
  - [Playlists](#playlists)
  - [Links](#links)
  - [Checkouts](#checkouts)
  - [Users](#users)
  - [Mobile Users](#mobile-users)
  - [Spotlight](#spotlight)
  - [Upload](#upload)

---

## Authentication

Most endpoints require authentication via an access token. Include the token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

### Authentication Levels

| Level | Short Name | Description                                           |
| ----- | ---------- | ----------------------------------------------------- |
| 1     | student    | Limited permissions, cannot edit certain video fields |
| 2     | member     | Standard member access                                |
| 3     | admin      | Full administrative access                            |

**Note:** Students are restricted from editing `allow_downloads`, `download_url`, `url`, and `status` fields on videos.

---

## Common Response Format

All successful responses return JSON with the requested data directly in the response body.

### Success Response

```json
{
  "id": 1,
  "title": "Example",
  ...
}
```

### List Response

```json
[
  { "id": 1, ... },
  { "id": 2, ... }
]
```

---

## Error Handling

### Error Response Format

```json
{
  "error": "error message description"
}
```

### HTTP Status Codes

| Code | Description                                                 |
| ---- | ----------------------------------------------------------- |
| 200  | Success                                                     |
| 400  | Bad Request - Invalid parameters or missing required fields |
| 401  | Unauthorized - Missing or invalid access token              |
| 403  | Forbidden - Insufficient permissions                        |
| 404  | Not Found - Resource does not exist                         |
| 409  | Conflict - Database or query error                          |
| 500  | Internal Server Error                                       |

---

## Data Models

### Video

| Field             | Type       | Description                    |
| ----------------- | ---------- | ------------------------------ |
| `id`              | integer    | Unique identifier              |
| `uuid`            | string     | Unique UUID identifier         |
| `title`           | string     | Video title                    |
| `description`     | string     | Video description              |
| `thumbnail`       | string     | Thumbnail image URL            |
| `cloudflare_id`   | string     | Cloudflare video ID (nullable) |
| `url`             | string     | Video URL                      |
| `download_url`    | string     | Download URL (nullable)        |
| `allow_downloads` | boolean    | Whether downloads are allowed  |
| `views`           | integer    | View count                     |
| `downloads`       | integer    | Download count                 |
| `status`          | GeneralNSN | Video status object            |
| `inserted_at`     | string     | Creation timestamp             |

### Asset

| Field         | Type                 | Description              |
| ------------- | -------------------- | ------------------------ |
| `id`          | integer              | Unique identifier        |
| `name`        | string               | Asset name               |
| `image_url`   | string               | Asset image URL          |
| `identifier`  | string               | Unique identifier string |
| `description` | string               | Asset description        |
| `category`    | GeneralNSN           | Category object          |
| `status`      | GeneralNSN           | Status object            |
| `metadata`    | AssetMetadata        | Asset metadata           |
| `active_tab`  | AssetCheckoutOmitted | Active checkout info     |
| `inserted_at` | string               | Creation timestamp       |

### AssetMetadata

| Field           | Type   | Description       |
| --------------- | ------ | ----------------- |
| `serial_number` | string | Serial number     |
| `manufacturer`  | string | Manufacturer name |
| `model`         | string | Model name        |
| `notes`         | string | Additional notes  |

### Playlist

| Field          | Type       | Description          |
| -------------- | ---------- | -------------------- |
| `id`           | integer    | Unique identifier    |
| `name`         | string     | Playlist name        |
| `status`       | GeneralNSN | Status object        |
| `description`  | string     | Playlist description |
| `banner_image` | string     | Banner image URL     |
| `route`        | string     | URL route/slug       |
| `inserted_at`  | string     | Creation timestamp   |
| `videos`       | Video[]    | Array of videos      |

### Link

| Field         | Type    | Description            |
| ------------- | ------- | ---------------------- |
| `id`          | integer | Unique identifier      |
| `route`       | string  | Short URL route        |
| `destination` | string  | Destination URL        |
| `active`      | boolean | Whether link is active |
| `creator`     | User    | Creator user object    |
| `clicks`      | integer | Click count            |
| `inserted_at` | string  | Creation timestamp     |

### Checkout

| Field             | Type       | Description          |
| ----------------- | ---------- | -------------------- |
| `id`              | integer    | Unique identifier    |
| `user`            | MobileUser | User who checked out |
| `associated_user` | integer    | Associated user ID   |
| `asset`           | Asset      | Checked out asset    |
| `asset_id`        | integer    | Asset ID             |
| `offsite`         | integer    | Offsite flag         |
| `checkout_status` | GeneralNSN | Checkout status      |
| `checkout_notes`  | string     | Checkout notes       |
| `time_out`        | string     | Checkout time        |
| `time_in`         | string     | Return time          |
| `expected_in`     | string     | Expected return time |

### User

| Field               | Type               | Description                |
| ------------------- | ------------------ | -------------------------- |
| `id`                | integer            | Unique identifier          |
| `username`          | string             | Username                   |
| `name`              | string             | Full name                  |
| `email`             | string             | Email address              |
| `profile_image_url` | string             | Profile image URL          |
| `authentication`    | GeneralNSN         | Auth level object          |
| `inserted_at`       | string             | Creation timestamp         |
| `last_active`       | string             | Last activity timestamp    |
| `strategies`        | UserAuthStrategies | Auth strategies (optional) |

### MobileUser

| Field               | Type       | Description        |
| ------------------- | ---------- | ------------------ |
| `id`                | integer    | Unique identifier  |
| `name`              | string     | Full name          |
| `email`             | string     | Email address      |
| `identifier`        | string     | Unique identifier  |
| `status`            | GeneralNSN | Status object      |
| `profile_image_url` | string     | Profile image URL  |
| `inserted_at`       | string     | Creation timestamp |

### Spotlight

| Field         | Type    | Description             |
| ------------- | ------- | ----------------------- |
| `rank`        | integer | Spotlight position/rank |
| `video_id`    | integer | Associated video ID     |
| `inserted_at` | string  | Creation timestamp      |
| `updated_at`  | string  | Last update timestamp   |
| `video`       | Video   | Associated video object |

### GeneralNSN

| Field        | Type    | Description        |
| ------------ | ------- | ------------------ |
| `id`         | integer | Status/category ID |
| `name`       | string  | Full name          |
| `short_name` | string  | Short name/code    |

---

## Endpoints

### Health Check

#### GET `/healthcheck`

Check if the API is running.

**Authentication:** Not required

**Response:** `200 OK` (empty body)

---

## Videos

### List Videos

#### GET `/core/v1.1/admin/videos`

Retrieve a paginated list of videos.

**Authentication:** Required

**Query Parameters:**

| Parameter | Type    | Required | Description                 |
| --------- | ------- | -------- | --------------------------- |
| `limit`   | integer | Yes      | Number of results to return |
| `offset`  | integer | Yes      | Number of results to skip   |
| `search`  | string  | No       | Search query for filtering  |
| `sort`    | string  | No       | Sort order                  |

**Response:**

```json
[
  {
    "id": 1,
    "uuid": "abc123-def456",
    "title": "Video Title",
    "description": "Video description",
    "thumbnail": "https://example.com/thumb.jpg",
    "cloudflare_id": "cf123",
    "url": "https://example.com/video.mp4",
    "download_url": "https://example.com/download.mp4",
    "allow_downloads": true,
    "views": 100,
    "downloads": 10,
    "status": {
      "id": 1,
      "name": "Published",
      "short_name": "published"
    },
    "inserted_at": "2024-01-01T00:00:00Z"
  }
]
```

---

### Get Video

#### GET `/core/v1.1/admin/video/{query}`

Retrieve a single video by ID or identifier.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type           | Description                         |
| --------- | -------------- | ----------------------------------- |
| `query`   | string/integer | Video ID (integer) or UUID (string) |

**Response:**

```json
{
  "id": 1,
  "uuid": "abc123-def456",
  "title": "Video Title",
  "description": "Video description",
  "thumbnail": "https://example.com/thumb.jpg",
  "cloudflare_id": "cf123",
  "url": "https://example.com/video.mp4",
  "download_url": "https://example.com/download.mp4",
  "allow_downloads": true,
  "views": 100,
  "downloads": 10,
  "status": {
    "id": 1,
    "name": "Published",
    "short_name": "published"
  },
  "inserted_at": "2024-01-01T00:00:00Z"
}
```

---

### Create Video

#### POST `/core/v1.1/admin/video`

Create a new video.

**Authentication:** Required

**Request Body:**

| Field          | Type    | Required | Description                    |
| -------------- | ------- | -------- | ------------------------------ |
| `title`        | string  | Yes      | Video title                    |
| `description`  | string  | Yes      | Video description              |
| `thumbnail`    | string  | Yes      | Thumbnail image URL            |
| `url`          | string  | Yes      | Video URL                      |
| `download_url` | string  | No       | Download URL                   |
| `status`       | integer | No       | Status ID (default: 2 - draft) |

**Request Example:**

```json
{
  "title": "New Video",
  "description": "Video description",
  "thumbnail": "https://example.com/thumb.jpg",
  "url": "https://example.com/video.mp4",
  "download_url": "https://example.com/download.mp4"
}
```

**Response:** Returns the created Video object

**Note:** If the authenticated user is a student, admin users will receive an email notification about the new upload.

---

### Update Video

#### PUT `/core/v1.1/admin/video/{query}`

Update an existing video.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type           | Description                         |
| --------- | -------------- | ----------------------------------- |
| `query`   | string/integer | Video ID (integer) or UUID (string) |

**Request Body:**

| Field     | Type   | Required | Description                        |
| --------- | ------ | -------- | ---------------------------------- |
| `changes` | object | Yes      | Object containing fields to update |

**Changes Object:**

| Field             | Type    | Description          |
| ----------------- | ------- | -------------------- |
| `title`           | string  | Video title          |
| `description`     | string  | Video description    |
| `thumbnail`       | string  | Thumbnail URL        |
| `allow_downloads` | boolean | Allow downloads flag |
| `download_url`    | string  | Download URL         |
| `url`             | string  | Video URL            |
| `status`          | integer | Status ID            |

**Request Example:**

```json
{
  "changes": {
    "title": "Updated Title",
    "description": "Updated description",
    "status": 1
  }
}
```

**Response:** Returns the updated Video object

**Note:** Students cannot modify `allow_downloads`, `download_url`, `url`, or `status` fields.

---

## Assets

### List Assets

#### GET `/core/v1.1/admin/assets`

Retrieve a paginated list of assets.

**Authentication:** Required

**Query Parameters:**

| Parameter | Type    | Required | Description                 |
| --------- | ------- | -------- | --------------------------- |
| `limit`   | integer | Yes      | Number of results to return |
| `offset`  | integer | Yes      | Number of results to skip   |
| `search`  | string  | No       | Search query for filtering  |
| `sort`    | string  | No       | Sort order                  |

**Response:**

```json
[
  {
    "id": 1,
    "name": "Camera A",
    "image_url": "https://example.com/camera.jpg",
    "identifier": "CAM-001",
    "description": "Professional camera",
    "category": {
      "id": 1,
      "name": "Cameras",
      "short_name": "camera"
    },
    "status": {
      "id": 1,
      "name": "Available",
      "short_name": "available"
    },
    "metadata": {
      "serial_number": "SN123456",
      "manufacturer": "Canon",
      "model": "EOS R5",
      "notes": "Handle with care"
    },
    "active_tab": null,
    "inserted_at": "2024-01-01T00:00:00Z"
  }
]
```

---

### Get Asset

#### GET `/core/v1.1/admin/asset/{query}`

Retrieve a single asset by ID or identifier.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type           | Description                               |
| --------- | -------------- | ----------------------------------------- |
| `query`   | string/integer | Asset ID (integer) or identifier (string) |

**Response:** Returns an Asset object

---

### Create Asset

#### POST `/core/v1.1/admin/asset`

Create a new asset.

**Authentication:** Required

**Request Body:**

| Field         | Type    | Required | Description       |
| ------------- | ------- | -------- | ----------------- |
| `name`        | string  | Yes      | Asset name        |
| `identifier`  | string  | Yes      | Unique identifier |
| `category`    | integer | Yes      | Category ID       |
| `image_url`   | string  | Yes      | Image URL         |
| `description` | string  | Yes      | Description       |
| `metadata`    | object  | Yes      | Asset metadata    |

**Metadata Object (all required):**

| Field           | Type   | Description   |
| --------------- | ------ | ------------- |
| `serial_number` | string | Serial number |
| `manufacturer`  | string | Manufacturer  |
| `model`         | string | Model name    |

**Request Example:**

```json
{
  "name": "Camera A",
  "identifier": "CAM-001",
  "category": 1,
  "image_url": "https://example.com/camera.jpg",
  "description": "Professional camera for video production",
  "metadata": {
    "serial_number": "SN123456",
    "manufacturer": "Canon",
    "model": "EOS R5"
  }
}
```

**Response:** Returns the created Asset object

---

### Update Asset

#### PUT `/core/v1.1/admin/asset/{query}`

Update an existing asset.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type           | Description                               |
| --------- | -------------- | ----------------------------------------- |
| `query`   | string/integer | Asset ID (integer) or identifier (string) |

**Request Body:**

| Field     | Type   | Required | Description                        |
| --------- | ------ | -------- | ---------------------------------- |
| `changes` | object | Yes      | Object containing fields to update |

**Changes Object:**

| Field         | Type    | Description           |
| ------------- | ------- | --------------------- |
| `name`        | string  | Asset name            |
| `image_url`   | string  | Image URL             |
| `identifier`  | string  | Unique identifier     |
| `description` | string  | Description           |
| `category`    | integer | Category ID           |
| `status`      | integer | Status ID             |
| `metadata`    | object  | Asset metadata object |

**Request Example:**

```json
{
  "changes": {
    "name": "Updated Camera A",
    "status": 2,
    "metadata": {
      "notes": "Needs maintenance"
    }
  }
}
```

**Response:** Returns the updated Asset object

---

## Playlists

### List Playlists

#### GET `/core/v1.1/admin/playlists`

Retrieve a paginated list of playlists.

**Authentication:** Required

**Query Parameters:**

| Parameter | Type    | Required | Description                 |
| --------- | ------- | -------- | --------------------------- |
| `limit`   | integer | Yes      | Number of results to return |
| `offset`  | integer | Yes      | Number of results to skip   |
| `search`  | string  | No       | Search query for filtering  |
| `sort`    | string  | No       | Sort order                  |

**Response:**

```json
[
  {
    "id": 1,
    "name": "Featured Videos",
    "status": {
      "id": 1,
      "name": "Active",
      "short_name": "active"
    },
    "description": "Our top featured content",
    "banner_image": "https://example.com/banner.jpg",
    "route": "featured",
    "inserted_at": "2024-01-01T00:00:00Z",
    "videos": []
  }
]
```

---

### Get Playlist

#### GET `/core/v1.1/admin/playlist/{query}`

Retrieve a single playlist by ID or route.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type           | Description                             |
| --------- | -------------- | --------------------------------------- |
| `query`   | string/integer | Playlist ID (integer) or route (string) |

**Response:** Returns a Playlist object with associated videos

---

### Create Playlist

#### POST `/core/v1.1/admin/playlist`

Create a new playlist.

**Authentication:** Required

**Request Body:**

| Field          | Type      | Required | Description          |
| -------------- | --------- | -------- | -------------------- |
| `name`         | string    | Yes      | Playlist name        |
| `description`  | string    | Yes      | Playlist description |
| `banner_image` | string    | Yes      | Banner image URL     |
| `route`        | string    | Yes      | URL route/slug       |
| `videos`       | integer[] | Yes      | Array of video IDs   |

**Request Example:**

```json
{
  "name": "Featured Videos",
  "description": "Our top featured content",
  "banner_image": "https://example.com/banner.jpg",
  "route": "featured",
  "videos": [1, 2, 3]
}
```

**Response:** Returns the created Playlist object

---

### Update Playlist

#### PUT `/core/v1.1/admin/playlist/{query}`

Update an existing playlist.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type    | Description |
| --------- | ------- | ----------- |
| `query`   | integer | Playlist ID |

**Request Body:**

| Field     | Type   | Required | Description                        |
| --------- | ------ | -------- | ---------------------------------- |
| `changes` | object | Yes      | Object containing fields to update |

**Changes Object:**

| Field           | Type      | Description          |
| --------------- | --------- | -------------------- |
| `name`          | string    | Playlist name        |
| `status`        | integer   | Status ID            |
| `description`   | string    | Playlist description |
| `banner_image`  | string    | Banner image URL     |
| `route`         | string    | URL route/slug       |
| `remove_videos` | integer[] | Video IDs to remove  |
| `add_videos`    | integer[] | Video IDs to add     |

**Request Example:**

```json
{
  "changes": {
    "name": "Updated Playlist Name",
    "add_videos": [4, 5],
    "remove_videos": [1]
  }
}
```

**Response:** Returns the updated Playlist object

---

## Links

### List Links

#### GET `/core/v1.1/admin/links`

Retrieve a paginated list of short links.

**Authentication:** Required

**Query Parameters:**

| Parameter | Type    | Required | Description                 |
| --------- | ------- | -------- | --------------------------- |
| `limit`   | integer | Yes      | Number of results to return |
| `offset`  | integer | Yes      | Number of results to skip   |
| `search`  | string  | No       | Search query for filtering  |
| `sort`    | string  | No       | Sort order                  |

**Response:**

```json
[
  {
    "id": 1,
    "route": "promo",
    "destination": "https://example.com/promotion",
    "active": true,
    "creator": {
      "id": 1,
      "name": "John Doe"
    },
    "clicks": 150,
    "inserted_at": "2024-01-01T00:00:00Z"
  }
]
```

---

### Get Link

#### GET `/core/v1.1/admin/link/{query}`

Retrieve a single link by ID or route.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type           | Description                         |
| --------- | -------------- | ----------------------------------- |
| `query`   | string/integer | Link ID (integer) or route (string) |

**Response:** Returns a Link object

---

### Create Link

#### POST `/core/v1.1/admin/link`

Create a new short link.

**Authentication:** Required

**Request Body:**

| Field         | Type   | Required | Description     |
| ------------- | ------ | -------- | --------------- |
| `route`       | string | Yes      | Short URL route |
| `destination` | string | Yes      | Destination URL |

**Request Example:**

```json
{
  "route": "promo",
  "destination": "https://example.com/promotion"
}
```

**Response:** Returns the created Link object

**Note:** The `created_by` field is automatically set to the authenticated user.

---

### Update Link

#### PUT `/core/v1.1/admin/link/{query}`

Update an existing link.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type    | Description |
| --------- | ------- | ----------- |
| `query`   | integer | Link ID     |

**Request Body:**

| Field     | Type   | Required | Description                        |
| --------- | ------ | -------- | ---------------------------------- |
| `changes` | object | Yes      | Object containing fields to update |

**Changes Object:**

| Field         | Type    | Description            |
| ------------- | ------- | ---------------------- |
| `route`       | string  | Short URL route        |
| `destination` | string  | Destination URL        |
| `active`      | boolean | Whether link is active |

**Request Example:**

```json
{
  "changes": {
    "destination": "https://example.com/new-promo",
    "active": false
  }
}
```

**Response:** Returns the updated Link object

---

## Checkouts

### List Checkouts

#### GET `/core/v1.1/admin/checkouts`

Retrieve a paginated list of asset checkouts.

**Authentication:** Required

**Query Parameters:**

| Parameter | Type    | Required | Description                 |
| --------- | ------- | -------- | --------------------------- |
| `limit`   | integer | Yes      | Number of results to return |
| `offset`  | integer | Yes      | Number of results to skip   |
| `sort`    | string  | No       | Sort order                  |

**Response:**

```json
[
  {
    "id": 1,
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    },
    "asset": {
      "id": 1,
      "name": "Camera A"
    },
    "offsite": 0,
    "checkout_status": {
      "id": 1,
      "name": "Checked Out",
      "short_name": "out"
    },
    "checkout_notes": null,
    "time_out": "2024-01-01T10:00:00Z",
    "time_in": null,
    "expected_in": "2024-01-02T10:00:00Z"
  }
]
```

---

### Update Checkout

#### PUT `/core/v1.1/admin/checkout/{query}`

Update a checkout (primarily for checking in assets).

**Authentication:** Required

**Path Parameters:**

| Parameter | Type    | Description |
| --------- | ------- | ----------- |
| `query`   | integer | Checkout ID |

**Request Body:**

| Field     | Type   | Required | Description                        |
| --------- | ------ | -------- | ---------------------------------- |
| `changes` | object | Yes      | Object containing fields to update |

**Changes Object:**

| Field      | Type    | Description                         |
| ---------- | ------- | ----------------------------------- |
| `check_in` | boolean | Set to `true` to check in the asset |

**Request Example:**

```json
{
  "changes": {
    "check_in": true
  }
}
```

**Response:** Returns the updated Checkout object

**Note:** When `check_in` is `true`, the checkout status is set to 2 (returned), `time_in` is set to the current time, and notes are updated to "checked in by admin".

---

## Users

### List Users

#### GET `/core/v1.1/admin/users`

Retrieve a paginated list of admin/team users.

**Authentication:** Required

**Query Parameters:**

| Parameter | Type    | Required | Description                 |
| --------- | ------- | -------- | --------------------------- |
| `limit`   | integer | Yes      | Number of results to return |
| `offset`  | integer | Yes      | Number of results to skip   |
| `search`  | string  | No       | Search query for filtering  |
| `sort`    | string  | No       | Sort order                  |

**Response:**

```json
[
  {
    "id": 1,
    "username": "johndoe",
    "name": "John Doe",
    "email": "john@example.com",
    "profile_image_url": "https://example.com/avatar.jpg",
    "authentication": {
      "id": 3,
      "name": "Administrator",
      "short_name": "admin"
    },
    "inserted_at": "2024-01-01T00:00:00Z",
    "last_active": "2024-01-15T14:30:00Z"
  }
]
```

---

### Get User

#### GET `/core/v1.1/admin/user/{query}`

Retrieve a single user by ID.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type    | Description |
| --------- | ------- | ----------- |
| `query`   | integer | User ID     |

**Response:** Returns a User object

---

### Update User

#### PUT `/core/v1.1/admin/user/{query}`

Update an existing user.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type    | Description |
| --------- | ------- | ----------- |
| `query`   | integer | User ID     |

**Request Body:**

| Field     | Type   | Required | Description                        |
| --------- | ------ | -------- | ---------------------------------- |
| `changes` | object | Yes      | Object containing fields to update |

**Changes Object:**

| Field               | Type    | Description                          |
| ------------------- | ------- | ------------------------------------ |
| `username`          | string  | Username                             |
| `name`              | string  | Full name                            |
| `email`             | string  | Email address (must be valid format) |
| `profile_image_url` | string  | Profile image URL                    |
| `authentication`    | integer | Authentication level ID              |

**Request Example:**

```json
{
  "changes": {
    "name": "John Smith",
    "email": "john.smith@example.com"
  }
}
```

**Response:** Returns the updated User object

---

## Mobile Users

### List Mobile Users (Public)

#### GET `/core/v1.1/list/mobileUsers`

Retrieve a paginated list of mobile users. This endpoint is publicly accessible.

**Authentication:** Not required

**Query Parameters:**

| Parameter | Type    | Required | Description                 |
| --------- | ------- | -------- | --------------------------- |
| `limit`   | integer | Yes      | Number of results to return |
| `offset`  | integer | Yes      | Number of results to skip   |
| `search`  | string  | No       | Search query for filtering  |
| `sort`    | string  | No       | Sort order                  |

---

### List Mobile Users (Admin)

#### GET `/core/v1.1/admin/mobileUsers`

Retrieve a paginated list of mobile users (admin version).

**Authentication:** Required

**Query Parameters:** Same as public endpoint

**Response:**

```json
[
  {
    "id": 1,
    "name": "Jane Doe",
    "email": "jane@example.com",
    "identifier": "JANE001",
    "status": {
      "id": 1,
      "name": "Active",
      "short_name": "active"
    },
    "profile_image_url": "https://example.com/avatar.jpg",
    "inserted_at": "2024-01-01T00:00:00Z"
  }
]
```

---

### Get Mobile User

#### GET `/core/v1.1/admin/mobileUser/{query}`

Retrieve a single mobile user by ID or identifier.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type           | Description                              |
| --------- | -------------- | ---------------------------------------- |
| `query`   | string/integer | User ID (integer) or identifier (string) |

**Response:** Returns a MobileUser object

---

### Create Mobile User

#### POST `/core/v1.1/admin/mobileUser`

Create a new mobile user.

**Authentication:** Required

**Request Body:**

| Field               | Type   | Required | Description                          |
| ------------------- | ------ | -------- | ------------------------------------ |
| `name`              | string | Yes      | Full name                            |
| `email`             | string | Yes      | Email address (must be valid format) |
| `identifier`        | string | Yes      | Unique identifier                    |
| `profile_image_url` | string | No       | Profile image URL                    |

**Request Example:**

```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "identifier": "JANE001",
  "profile_image_url": "https://example.com/avatar.jpg"
}
```

**Response:** Returns the created MobileUser object

---

### Update Mobile User

#### PUT `/core/v1.1/admin/mobileUser/{query}`

Update an existing mobile user.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type    | Description    |
| --------- | ------- | -------------- |
| `query`   | integer | Mobile User ID |

**Request Body:**

| Field     | Type   | Required | Description                        |
| --------- | ------ | -------- | ---------------------------------- |
| `changes` | object | Yes      | Object containing fields to update |

**Changes Object:**

| Field               | Type    | Description       |
| ------------------- | ------- | ----------------- |
| `name`              | string  | Full name         |
| `email`             | string  | Email address     |
| `identifier`        | string  | Unique identifier |
| `status`            | integer | Status ID         |
| `profile_image_url` | string  | Profile image URL |

**Request Example:**

```json
{
  "changes": {
    "name": "Jane Smith",
    "status": 2
  }
}
```

**Response:** Returns the updated MobileUser object

---

## Spotlight

### List Spotlights

#### GET `/core/v1.1/admin/spotlight`

Retrieve a list of spotlight entries.

**Authentication:** Required

**Query Parameters:**

| Parameter | Type    | Required | Description                 |
| --------- | ------- | -------- | --------------------------- |
| `limit`   | integer | Yes      | Number of results to return |
| `offset`  | integer | Yes      | Number of results to skip   |

**Response:**

```json
[
  {
    "rank": 1,
    "video_id": 42,
    "inserted_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-15T00:00:00Z",
    "video": {
      "id": 42,
      "title": "Featured Video",
      ...
    }
  }
]
```

---

### Get Spotlight

#### GET `/core/v1.1/admin/spotlight/{rank}`

Retrieve a single spotlight entry by rank.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type    | Description             |
| --------- | ------- | ----------------------- |
| `rank`    | integer | Spotlight rank position |

**Response:** Returns a Spotlight object with associated video

---

### Update Spotlight

#### PUT `/core/v1.1/admin/spotlight/{rank}`

Update a spotlight entry.

**Authentication:** Required

**Path Parameters:**

| Parameter | Type    | Description             |
| --------- | ------- | ----------------------- |
| `rank`    | integer | Spotlight rank position |

**Request Body:**

| Field      | Type    | Required | Description         |
| ---------- | ------- | -------- | ------------------- |
| `video_id` | integer | Yes      | Video ID to feature |

**Request Example:**

```json
{
  "video_id": 42
}
```

**Response:** Returns the updated Spotlight object

---

## Upload

### Upload Image

#### POST `/core/v1.1/admin/upload`

Upload an image file to S3.

**Authentication:** Required

**Content-Type:** `multipart/form-data`

**Form Data:**

| Field   | Type    | Required | Description                                     |
| ------- | ------- | -------- | ----------------------------------------------- |
| `id`    | integer | Yes      | Associated resource ID                          |
| `route` | string  | Yes      | S3 route category (must be in permitted routes) |
| `image` | file    | Yes      | Image file (max 10MB)                           |

**Response:**

```json
{
  "url": "https://s3.example.com/uploads/route/id/image.jpg"
}
```

**Note:** The `route` parameter must be a valid, permitted S3 route as defined in the application configuration.

---

## Rate Limiting

Currently, the API does not implement rate limiting. However, it's recommended to implement reasonable request intervals in client applications.

---

## CORS

The API supports CORS with the following configuration:

- **Allowed Origins:** `*` (all origins)
- **Allowed Methods:** `GET`, `HEAD`, `POST`, `PUT`, `OPTIONS`
- **Allowed Headers:** `X-Requested-With`, `Content-Type`, `Origin`, `Authorization`, `Accept`, `Cookie`, `X-CSRF-Token`
- **Credentials:** Allowed

---

## Changelog

### v1.1 (Current)

- Initial documented version
- Full CRUD operations for videos, assets, playlists, links
- User and mobile user management
- Checkout system for asset management
- Spotlight feature for featured content
- Image upload functionality
