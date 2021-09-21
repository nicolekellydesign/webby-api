# Responses

This documents all of the JSON responses that API routes respond with.

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
    {
      "username": string
    },
    . . . more users
  ]
}
```
