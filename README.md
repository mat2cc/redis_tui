# redis_tui
<p>
    <a href="https://github.com/mat2cc/redis_tui/releases"><img src="https://img.shields.io/github/release/mat2cc/redis_tui.svg" alt="Latest Release"></a>
    <a href="https://github.com/mat2cc/redis_tui/actions"><img src="https://github.com/mat2cc/redis_tui/workflows/test/badge.svg" alt="Build Status"></a>
</p>

Redis terminal browser written in go, with the use of charmbraclet's [bubbletea](https://github.com/charmbracelet/bubbletea)

https://github.com/mat2cc/redis_tui/assets/19396939/8389236f-71c7-4062-b103-4c53a1ca94db

## Install

**homebrew**:
```sh
brew install mat2cc/tap/redis_tui
```

**go**:
```sh
go install github.com/mat2cc/redis_tui@latest
```

## Usage

```sh
redis_tui -address localhost:6379 -db 2
```

Enter your `redis://` address, db, etc. and connect to redis. 
*note:* this will also work wish ssh tunneling into a redis database.
Check `--help` for more options!

While in the redis_tui, press `?` for movement instructions

## Under the hood

redis_tui is powered by the [redis SCAN command](https://redis.io/commands/scan) where we are gathering a variable number (configured by the `--scan-size` arguments, with default of 1000) of keys per scan. A cursor is kept so that every time the `m` key is hit, more keys will be fetched.

To get the redis type information, we are using a [Redis Pipeline](https://redis.io/docs/manual/pipelining/) which will batch all of the type reqests for each scan into one transaction. This can be disabled by setting `--include-types=false`, where type information will only be retrieved when opening the details view for a specific key.

## Production Considerations
The redis SCAN command can safely be used in production since we can limit the number of keys returned per request. Some things to keep in mind when querying production redis or when performance in a concern:
- reduce the number of elements returned per request with the `--scan-size` argument (the `m` key can always fetch more keys)
- if possible query against a replica redis db
- set the `--include-types=false` to omit the redis pipeline that gathers type data for all keys
