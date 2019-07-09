

Configuration
================

> Documentation is a little incomplete at the moment.

You may setup as many paths as you like to be served however
they must not conflict with one another!

All of our configuration takes place in a configuration file `./data/conf.json`...

We're still missing some additional configuration elements which need
to be hard-coded until this is resolved.

```json
{
  "serv": {
    "host": "localhost",
    "port": ":5500",
    "tls": false,
    "key": "data\\key.pem",
    "crt": "data\\cert.pem",
    "path": "v"
  },
  "root": {
    "path": "/",
    "dir": ".\\public",
    "files": [
      "json.json",
      "bundle.js",
      "favicon.ico"
    ],
    "alias": [
      "home",
      "index.htm",
      "index.html",
      "index",
      "default",
      "default.htm"
    ],
    "default": "index.html"
  },
  "stat": [
    {
      "src": "public\\images",
      "tgt": "/images/",
      "nav": true
    },
    {
      "src": "public\\static",
      "tgt": "/static/",
      "nav": true
    }
  ],
  "indx": [
    {
      "src": "multi-media\\public",
      "tgt": "/v/",
      "nav": true,
      "serve": true,
      "ignorePaths": [],
      "spec": [
        "Media"
      ]
    }
  ],
  "spec": [
    {
      "name": "Media",
      "ext": [
        ".mp4",
        ".m4a",
        ".mp3"
      ]
    },
    {
      "name": "Markdown",
      "ext": [
        ".md",
        ".mmd"
      ]
    }
  ]
}
```


<!-- os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600 -->
<!-- https://github.com/d4l3k/messagediff -->
<!-- https://github.com/davecgh/go-spew -->
<!-- https://github.com/sergi/go-diff -->
<!-- https://github.com/STRML/react-grid-layout#demos -->
<!-- https://transform.now.sh/css-to-js/ -->
<!-- https://github.com/ritz078/transform-www -->
<!-- https://jsvault.com/async-parallel -->
<!-- https://github.com/Microsoft/monaco-editor-webpack-plugin -->
<!-- https://github.com/gin-contrib/secure -->






