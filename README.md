# twirp-openapi-gen

A Twirp RPC OpenAPI generator implemented as `protoc` plugin

*Currently supports only OpenAPI 2.0*

# Usage

Installing the generator for protoc/buf:

```
go install github.com/albenik/twirp-openapi-gen/cmd/protoc-gen-twirp-openapi@latest
```

## Run whith the protoc

```
protoc -twirp-openapi_out=api/foo -twirp-openapi_opt=hostname=api.example.com,path_prefix=/api/v1 rpc/service.proto
```

## Run with [buf.build](https://buf.build):

`buf.gen.yaml`:

```yaml
version: v1
plugins:
  - name: twirp-openapi
    out: api/foo
    opt:
     - hostname=api.example.com
     - path_prefix=/api/v1
```

# Motivation

The only working generator referenced by [Twirp](https://github.com/twitchtv/twirp) project and found by me
is [go-bridget/protoc-gen-twirp_swagger](https://github.com/go-bridget/twirp-swagger-gen). It lacks support of map (at
least) and also does not work with advanced buf configuration when `buf.work.yaml` and `buf.gen.yaml` placed on the
different levels.

I started from a PR to the original [go-bridget/twirp-swagger-gen](https://github.com/go-bridget/twirp-swagger-gen) but later
decided to write generator from scratch. I dropped stadalone proto parser package and just rely on already parsed
structure from `protogen`. Also I don't use OpenAPI spec go-package and replaced it by hand-written structures, which
are incomplete buf much more comfortable to use.

The generated output is suitable for Swagger-UI.
