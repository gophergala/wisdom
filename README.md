# Wisdom

Wisdom is an open API for [startupquote.com](http://startupquote.com/).

The idea was to provide an access to famous quote from startup ecosystem 
that can resused by developer inside their apps.

## Mission

1. Simple
2. Fast (because it's written in [Go](http://golang.or))
3. Open for everyone (no authentication required)

## API

### Object
`quote` has the following form

```json
{
    "id": 13,
    "post_id": "104365553355",
    "author": {
        "id": 12,
        "avatar_url": "",
        "name": "Leslie Bradshaw",
        "company": "JESS3",
        "twitter_username": "lesliebradshaw"
    },
    "content": "In my 20s I was thrashing around in the water, trying to keep my head above it. In my 30s, I realized it was only three feet deep and I stood up.",
    "permalink": "http://startupquote.com/post/104365553355",
    "picture_url": "http://36.media.tumblr.com/bc0698c00b443d0e5c6b9b814d74bbd9/tumblr_nfywvtqr2C1qz6pqio1_r1_1280.png",
    "tags": [
        {
            "id": 26,
            "label": "entrepreneur"
        },
        {
            "id": 27,
            "label": "founder"
        }
    ]
}
```

`author` has the following form

```json
{
    "id": 1,
    "avatar_url": "",
    "name": "Reid Hoffman",
    "company": "Linkedin",
    "twitter_username": "reidhoffman"
}
```

### Random

| Endpoint  | Description |
| --------- | ------ |
| `/v1/random` | return a random `quote`|

example valid request are:

```
curl -i -H "Accept: application/json" -X GET https://wisdomapi.herokuapp.com/v1/random
```

### Post

| Endpoint  | Description |
| --------- | ------ |
| `/v1/post/:id` | return an `object` that have given `:id`|

`:id` is a integer

TODO: add example of post response


### Author

| Endpoint  | Description |
| --------- | ------ |
| `/v1/authors` | return an array of `author`|
| `/v1/author/:twitter_username` | return an array of `quote` that have given `:twitter_username`|
| `/v1/author/:twitter_username/random` | return a random `quote` that have given `:twitter_username`|

`:twitter_username` is a string

Example request:

```
# Authors
curl -i -H "Accept: application/json" -X GET https://wisdomapi.herokuapp.com/v1/authors
```


### Tag

| Endpoint  | Description |
| --------- | ------ |
| `/v1/tags` | return an array of `tags`|
| `/v1/tags/:tags` | return an array of `object` that have given `:tags`|
| `/v1/tags/:tags/random` | return a random `object` that have given `:tags`|

`:tags` is a string or comma separated string.

TODO: add example of tag response
