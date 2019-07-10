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

Features
-------------

CURRENT


- **JSON server configuration**.
- **XHR/JSON file-system index**
    - with GET|POST requests with manual refresh capability
    - file-extention and directory filters
    - XHR/JSON Tag requests (audio/video metadata) for multi-media (using [github.com/dhowden/tag])
- smart CLI interface for overriding config settings like
    - `--port <number>`
    - `--tls`: supply this flag to use tls when the config
      file has it off by default.
- More XHR request/data integrations to come perhaps including Calibre EBOOK
  data, plex (meta-info) and Chrome bookmarks/favicons, although
  is yet to be determined exactly how and when at this point.

IN PROGRESS

- Logon sessions (only sqlite3 data backend for now) are nearly complete.
- Separate demo sandbox projects (soon) 

KNOWN BUGS for expected fixes

- #1 **I'd like to see file time-stamps (CRD)**  
  This may only be implemented in windows since thats the main dev
  workstation.  PRs and discussions (bug section) are welcome.
- #2 [bug] *if file-extensions are poorly configured ATM:*  
  Multiple/duplicate files are returned in XHR/JSON due to extension definitions sharing the same extension.  
  (will be fixed soon)
- #3 [feature] **remove long-empty path entries**  
  The idea is to add a post indexing filter that strips out all empty directories and to provide a JSON option to apply such a filter.

Project Status: alpha (development) phase v0.0.0 has not changed yet for 90 revisions ATM.