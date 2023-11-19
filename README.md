# redis_tui

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
