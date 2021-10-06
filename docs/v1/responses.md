# Responses

This documents all of the JSON responses that API routes respond with.

## Check

This is returned when a client sends a request to check if a connection has a valid login session.

```json
{
  "valid": bool
}
```

## Gallery

This is returned when a client sends an API request to get all portfolio gallery items.

If an item for some reason has no slides, then the item will just have an empty slides array.

If there are no gallery items, an empty array is returned.

```json
{
  "items": [
    {
      "id": string,
      "title_line_1": string,
      "title_line_2": string,
      "thumbnail_location": string,
      "thumbnail_caption": string,
      "slides": [
        {
          "gallery_id": string,
          "name": string,
          "title": string,
          "caption": string,
          "location": string
        },
        . . . more slides
      ]
    },
    . . . more items
  ]
}
```

## Photos

This is returned when a client sends an API request to get all photography gallery items.

If there are no photo items, an empty array is returned.

```json
{
  "photos": [
    {
      "filename": string
    },
    . . . more items
  ]
}
```

## Users

This is returned when a client sends an API request to get all users.

```json
{
  "users": [
    . . . users
  ]
}
```

Here is what the returned User struct looks like:

```json
{
  "id": number,
  "username": string,
  "protected": bool
}
```
