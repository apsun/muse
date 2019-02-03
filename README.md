# muse

A logical extension of Python's `http.server` to include metadata in directory
listings. The output is returned in JSON format.

## Examples

Example response:

```
[
    {
        "name": "my secret stash",
        "type": "directory",
    },
    {
        "name": "Foo - Bar.mp3",
        "type": "audio/mpeg",
        "id3": {
            "artist": "Foo",
            "title": "Bar",
        },
        "duration": 353,
    },
    {
        "name": "lolcat.jpg",
        "type": "image/jpeg",
    },
    {
        "name": "hello.txt",
        "type": "text/plain",
    },
    {
        "name": "random.bin",
        "type": "application/octet-stream",
    },
]
```
