# Endpoints

All endpoints are inside the `/api` route. So for example, the endpoint to get all gallery items would be at `/api/gallery`.

## Public Routes

These routes can be used by anyone, though the logout and refresh endpoints need a valid session. They aren't considered admin endpoints, they simply need a valid session to make any sense.

### Session Endpoints

These endpoints are used for session management.

#### `/login`: POST

Endpoint to log a user in. If the username and password match what's in the database, a cookie will be set with a session token to be used for all further admin requests.

Each session is only valid for 5 minutes, unless the `extended` flag is set in the JSON request. Sessions may be refreshed as long as the session is still valid. Refreshing a session will extend the session by another 5 minutes.

If `extended` is set to `true` in the JSON request, the session will be valid for 30 days. This key is optional and may be omitted from the request.

```json
{
  "username": string,
  "password": string,
  "extended": bool
}
```

#### `/logout`: POST

Endpoint to log a user out. This requires a valid session to work, which should be pretty self-explanitory.

#### `/refresh`: POST

Endpoint to refresh an existing session. If the request has a valid session cookie, the session will be extended 5 minutes from the refresh request.

### Other endpoints

#### `/gallery`: GET

Endpoint to get all stored gallery items. See the responses doc for the structure of the returned JSON object.

#### `/photos`: GET

Endpoint to get all stored photography gallery items. See the responses doc for the structure of the returned JSON object.

## Admin Routes

All admin routes are in the `/api/admin` space and require a valid session to interact with.

### Gallery

These routes are for managing items and slides in the main portfolio gallery.

#### `/gallery`: POST

Adds a new gallery item from a JSON body.

```json
{
  "item": {
    "id": string,
    "title_line_1": string,
    "title_line_2": string,
    "thumbnail_caption": string,
    "thumbnail_location": string
  }
}
```

#### `/gallery/:id`: DELETE

Removes a gallery item with the given ID. If no item exists with the ID, HTTP status `404` will be returned.

#### `/gallery/:id/slides`: POST

Adds a new slide to the gallery item with the given ID.

```json
{
  "slide": {
    "name": string,
    "title": string,
    "caption": string,
    "location": string
  }
}
```

#### `/gallery/:id/slides/:name`: DELETE

Removes a slide with the given name from the gallery item that has the given ID.

### Photos

These routes are for managing pictures in the photography gallery.

#### `/photos`: POST

Add a new photo to the photography gallery from a JSON body.

```json
{
  "filename": string
}
```

#### `/photos/:fileName`: DELETE

Removes a photo that has the given file name.

### Users

These routes are for viewing and managing administrators.

#### `/users`: GET

Returns a list of usernames. This is a privileged endpoint for an extra layer of security.

#### `/users`: POST

Adds a new administrator. Expects the following JSON body:

```json
{
  "username": string,
  "password": string
}
```
