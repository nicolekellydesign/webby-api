# Responses

This documents all of the JSON responses that API routes respond with.

## About

This is returned when a client requests the about page designer statement.

```json
{
  "statement": string,
}
```

## Check

This is returned when a client sends a request to check if a connection has a valid login session.

```json
{
  "valid": bool
}
```

## Gallery

This is returned when a client sends an API request to get all portfolio gallery items.

If an item has no images, then the item will just have an empty images array.

If there are no gallery items, an empty array is returned.

```json
{
  "items": [
    {
      "id": string,
      "title": string,
      "caption": string,
      "projectInfo": string,
      "thumbnail": string,
      "embedURL": string,
      "images": [
        . . . string,
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
      "id": number,
      "username": string,
      "protected": bool
    },
    . . . users
  ]
}
```
