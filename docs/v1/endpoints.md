# Endpoints

All endpoints for the V1 API are inside the `/api/v1` route. So for example, the endpoint to get all gallery items would be at `/api/v1/gallery`.

## Public Routes

These routes can be used by anyone, though the logout endpoint needs a valid session. It's not considered an admin endpoint, it simply needs a valid session to make any sense.

### Session Endpoints

These endpoints are used for session management.

#### `/check`: GET

Checks if a connection has a valid login session. See the responses documentation for the returned object.

#### `/login`: POST

Logs a user in. If the username and password match what's in the database, a cookie will be set with a session token to be used for all further admin requests.

Each session is only valid for the current browser session, unless the `extended` flag is set in the JSON request.

If `extended` is set to `true` in the JSON request, the session will be valid for 30 days. This key is optional and may be omitted from the request.

```json
{
  "username": string,
  "password": string,
  "extended": bool
}
```

#### `/logout`: POST

Logs a user out. This requires a valid session to work, which should be pretty self-explanatory.

### Other endpoints

#### `/about`: GET

Gets the designer statement text from the server.

#### `/gallery`: GET

Gets all stored gallery items.

#### `/gallery/:name`: GET

Gets the details for a project with the given name.

#### `/photos`: GET

Endpoint to get all stored photography gallery items.

## Admin Routes

All admin routes are in the `/api/v1/admin` space and require a valid session to interact with.

### Gallery

These routes are for managing items and slides in the main portfolio gallery.

#### `/gallery`: POST

Adds a new gallery item to the database with a thumbnail. It expects a multipart-form body with these keys:

```
name: string
thumbnail: File
embed_url: string | undefined
title: string
caption: string
project_info: string
```

#### `/gallery/:id`: PUT

Updates a project. The values to update are taken from a JSON body with the format:

```json
{
  "title": string,
  "caption": string,
  "projectInfo": string,
  "embedURL": string | undefined
}
```

#### `/gallery/:id`: DELETE

Removes a gallery item with the given ID. If no item exists with the ID, HTTP status `404` will be returned.

#### `/gallery/:id/thumbnail`: PATCH

Updates the thumbnail for a project. The body should be a multipart-form with the image set to the `thumbnail` key.

#### `/gallery/:id/images`: POST

Adds a the given image names to the database, associating them with the project ID. The body should be a JSON array of the file names.

This doesn't handle the uploading of the images; see the `upload` endpoint.

#### `/gallery/:id/images`: DELETE

Removes images associated with a project from the database and filesystem. The body should be a JSON array of the file names to remove.

### Photos

These routes are for managing pictures in the photography gallery.

#### `/photos`: POST

Add new photos to the photography gallery. The body should be a JSON array of the file names.

This doesn't handle the uploading of the images; see the `upload` endpoint.

#### `/photos`: DELETE

Removes a list of photos from the database and filesystem. The body should be a JSON array of the file names to remove.

### Users

These routes are for viewing and managing administrators.

#### `/users`: GET

Returns a list of usernames. This is a privileged endpoint for an extra layer of security.

#### `/users`: POST

Adds a new administrator. The endpoint expects the following JSON body:

```json
{
  "username": string,
  "password": string
}
```

#### `/users/:id`: DELETE

Removes an administrator. An admin cannot delete themselves.

### Upload

#### `/upload`: POST

Upload a file to the server. If the file's MIME type is that of an image, the file will be saved to the `images` subdirectory of the project root. All other files will be saved to the `resources` subdirectory.

The request body should be a multipart-form with the file set to the key `file`.
