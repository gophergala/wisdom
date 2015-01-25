# Wisdom

Wisdom is an open API for [startupquote.com](http://startupquote.com/).

The idea was to provide an access to famous quote from startup ecosystem 
that can resused by developer inside their apps.

## Mission

1. Simple
2. Fast (because it's written in [Go](http://golang.or))
3. Open for everyone (no authentication required)

## API

### Quote object
every `object` quote has the following form

```json
{
    "number": 1,
    "post_id": 123456,
    "author": {
        "name":"name of author",
        "company":"company name",
        "twitter_url":"http://twitter.com/author_username",
        "profile_url":"http://some.url/on/the/internet"
    },
    "content": "content of quote",
    "permalink": "http://startupquote.com/post/post_id_of_quote",
    "tag":"comma,separated,value",
    "picture_url":"http://some.url/on/the/internet"
}
```

### Random

| Endpoint  | Description |
| --------- | ------ |
| `/v1/random` | return a random `object`|

example valid request are:

```
$ curl -i -H "Accept: application/json" -X GET https://wisdomapi.herokuapp.com/v1/random
HTTP/1.1 200 OK
Connection: keep-alive
Content-Type: application/json; charset=utf-8
Server: Wisdom powered by Gophergala
X-Wisdom-Media-Type: wisdom.V1
Date: Sun, 25 Jan 2015 12:55:32 GMT
Content-Length: 620
Via: 1.1 vegur

{
    "id":17,
    "post_id":"81784898975",
    "author":{
        "id":0,
        "avatar_url":"",
        "name":"Steve Jobs",
        "company":"Apple",
        "twitter_username":""
    },
    "content":"Your work is going to fill a large part of your life, and the only way to be truly satisfied is to do what you believe is great work. And the only way to do great work is to love what you do.",
    "permalink":"http://startupquote.com/post/81784898975",
    "picture_url":"http://41.media.tumblr.com/0915e5401eb9a4e0e01e27c525198562/tumblr_n3keo0aVDc1qz6pqio1_r1_1280.png",
    "tags":[
        {"
            id":23,
            "label":"classic"
        },
        {
            "id":29,
            "label":"work"
        },
        {
            "id":31,
            "label":"life"
        },
        {
            "id":32,
            "label":"enthusiasm"
        }
    ]
}
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
| `/v1/author/:twitter_username` | return an array of `object` that have given `:twitter_username`|
| `/v1/author/:twitter_username/random` | return an `object` that have given `:twitter_username`|

`:twitter_username` is a string

TODO: add example of author response

### Tag

| Endpoint  | Description |
| --------- | ------ |
| `/v1/tags` | return an array of `tags`|
| `/v1/tags/:tags` | return an array of `object` that have given `:tags`|
| `/v1/tags/:tags/random` | return a random `object` that have given `:tags`|

`:tags` is a string or comma separated string.

TODO: add example of tag response
