podder
=====

Podcast feed paywall/authentication proxy.

## Synopsis

Many sites now sell access to podcast feeds and put their feeds behind a
paywall. Unfortunately not all podcast clients support authentication. This web
app, written in Go, exposes a proxy server that performs the authentication step
allowing any client to download paywalled podcasts.

## Features

* Streams audio/video files straight through to the client. Nothing is stored on
the server.
* Content headers are preserved so progress, pause and resume all work in
clients.
* Uses Go's excellent HTTP cookiejar library for compatiblity with many sites.
* Supports URL rewriting (needed on some servers)
* Designed against a specific site but can be extended and may just work
on other sites. Site designed against has been deliberately anonymized.
* No external dependancies.

## Usage

To compile and run locally, ensure you have [Go](http://golang.org), clone
the repo and run this command (make sure to replace the env vars with your own)

``` bash
$ REPLACE_URL=www.example.com/podcast FILE_URL=http://media.example.com/files LOGIN_URL=http://www.example.com/signin.php FEED_URL=http://www.example.com/podcast.rss PASS=pass USER=bob PORT=5000 go run podder.go
```

Then navigate to `http://localhost:5000/feed.xml` in your browser or podcast
client.

To run your own hosted version it's really easy to deploy to Heroku by running

```bash
$ heroku create -b https://github.com/kr/heroku-buildpack-go.git
$ git push heroku master
```

Make sure to set all the environment variables before deploying.

## Config

Environment variables are used to configure the app. Here is what they do.

Variable | Purpose
--- | --- | ---
`REPLACE_URL` | Base URL to replace in the feed.
`FILE_URL` | Base URL of where files are on target site.
`LOGIN_URL` | URL of the login page.
`FEED_URL` | Source podcast feed URL.
`PASS` | Password for the site.
`USER` | Username for the site.

## To do

* Support more sites
* Handle missing env vars

[1]: http://www.bleedingcool.com
