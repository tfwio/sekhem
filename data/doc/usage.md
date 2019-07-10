[home]: ../../readme.md "github.com/tfwio/sekhm/readme.md"
[features]: features.md
[configuration]: configuration.md
[build]: build.md
[usage]: usage.md
<!-- []:  -->

- [home]
    - [features]
    - [configuration]
    - [usage]
    - [build]


Usage
================

The first time you run `./srv-${GOARCH}`, it will generate a example configuration file, `data/conf.json`, and exit.

Running the application once more will run the server on your localhost (http://localhost:5500/) and serve a demo html file that's sitting in the `./public` directory.

You can run the server on a differt port by calling

```bash
./srv-${GOARCH} --port <number>
```

Though there are the sub-commands `run` | `go` which do the same thing, you don't need to use them as its the default context — to run the server.

Running Help:
```bash
NAME:
   srv-amd64.exe

USAGE:
    [global options] command [command options] [arguments...]

VERSION:
   v0.0.0

AUTHOR:
   tfw; et alia

COMMANDS:
     run, go    Runs the server.
     make-conf  srv-amd64.exe make-conf <[file-path].json>
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --tls          Wether or not to use TLS.
   --port value   UseHost is identifies the host to use to over-ride JSON config. (default: 0)
   --conf value   Points to a custom configuration file. (default: "data/conf.json")
   --version, -v  print the version

COPYRIGHT:
   github.com/tfwio/sekhm

   This is free, open-source software.
   disclaimer: use at own risk.

```
