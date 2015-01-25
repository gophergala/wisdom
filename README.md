# Wisdom

Wisdom is a Startup Quote API.

The idea is to provide an access to famous quotes by people from startup ecosystem that can be reused by developer inside their apps.

Potential apps is like [Product Hunt](http://producthunt.com), [Beta List](http://betalist.com) 
and other related startup community.

## Mission

1. Simple
2. Fast (because it's written in [Go](http://golang.org))
3. Open for everyone (no authentication required)

## API

### Object
#### Quote
example

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

#### Author
example

```json
{
    "id": 1,
    "avatar_url": "",
    "name": "Reid Hoffman",
    "company": "Linkedin",
    "twitter_username": "reidhoffman"
}
```

#### Tag
example

```json
{
    "id": 1,
    "label": "launch"
}
```

### Random

| Endpoint  | Description |
| --------- | ------ |
| `/v1/random` | return a random `quote`|

#### Example request

```
GET https://wisdomapi.herokuapp.com/v1/random
```

### Author

| Endpoint  | Description |
| --------- | ------ |
| `/v1/authors` | return an array of `author`|
| `/v1/author/:twitter_username` | return an array of `quote` by author that have given `:twitter_username`. If author doesn't have twitter account response will be 404|
| `/v1/author/:twitter_username/random` | return a random `quote` by author that have given `:twitter_username`. If author only have 1 quote the response will be exactly the same, not random.|

`:twitter_username` is a string

#### Example request

```
# List of author
GET https://wisdomapi.herokuapp.com/v1/authors

# List of quotes by Paul Graham (@paulg)
GET https://wisdomapi.herokuapp.com/v1/author/paulg

# Random quote by Paul Graham (@paulg)
GET https://wisdomapi.herokuapp.com/v1/author/paulg/random
```


### Tags

| Endpoint  | Description |
| --------- | ------ |
| `/v1/tags` | return an array of `tag`|

#### Example request

```
# List of tags
GET https://wisdomapi.herokuapp.com/v1/tags
```