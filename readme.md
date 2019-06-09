This is a little golang web-server app for servicing NPM sandboxes
for react or the like.

It generally creates a web-server and will likely use statik for
compiling once you might be happy with a simple react app distro.

The primary purpose of the app is to create a file-system index
wrapper that can serve JSON file-system indexes.

This was written to support go-mmd-fs or whatever its called at this point.

Configuration
----------------

You may setup as many paths as you like to be served however
they must not conflict with one another!

All of our configuration takes place in a configuration file `./data/conf.json`...

We're still missing some additional configuration elements which need
to be hard-coded until this is resolved.

```json
{
    "serv": {
        "host": "tfw.io",
        "port": ":5500",
        "tls": false,
        "key": "data/ia.key",
        "crt": "data/ia.crt",
        "path": "v"
    },
    "root": {
        "path": "/",
        "dir": ".\\public",
        "files": ["json.json", "bundle.js", "favicon.ico"],
        "alias": ["home", "index.htm", "index.html", "index", "default", "default.htm"],
        "default": "index.html"
    },
    "stat": [
      {
        "src": "public\\images",
        "tgt": "/images/",
        "nav": true
      }, {
        "src": "public\\static",
        "tgt": "/static/",
        "nav": true
      }, {
        "src": "[some directory with media files]",
        "tgt": "/v/",
        "nav": true
      }
    ]
}
```

