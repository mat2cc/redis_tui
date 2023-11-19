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

Enter your `redis://` address, db, etc. and connect to redis. Check `--help` for more options
```sh
redis_tui -address localhost:6379 -db 2
```

While in the redis_tui, press `?` for movement instructions
