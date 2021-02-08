# gitinfo

> Get git repo metagata (lib) and generate gitinfo.json via go generate (cmd)

[![GoDoc][gd1]][gd2]
 [![codecov][cc1]][cc2]
 [![Build Status][bs1]][bs2]
 [![GoCard][gc1]][gc2]
 [![GitHub Release][gr1]][gr2]
 [![GitHub license][gl1]][gl2]

[bs1]: https://cloud.drone.io/api/badges/pgmig/gitinfo/status.svg
[bs2]: https://cloud.drone.io/pgmig/gitinfo
[cc1]: https://codecov.io/gh/pgmig/gitinfo/branch/master/graph/badge.svg
[cc2]: https://codecov.io/gh/pgmig/gitinfo
[gd1]: https://godoc.org/github.com/pgmig/gitinfo?status.svg
[gd2]: https://godoc.org/github.com/pgmig/gitinfo
[gc1]: https://goreportcard.com/badge/github.com/pgmig/gitinfo
[gc2]: https://goreportcard.com/report/github.com/pgmig/gitinfo
[gr1]: https://img.shields.io/github/release/pgmig/gitinfo.svg
[gr2]: https://github.com/pgmig/gitinfo/releases
[gl1]: https://img.shields.io/github/license/pgmig/gitinfo.svg
[gl2]: https://github.com/pgmig/gitinfo/blob/master/LICENSE

This package uses external `git` binary for creating a file with git metagata like

```json
{
  "version": "v0.12.0-1-g99a5776",
  "repository": "git@github.com:pgmig/gitinfo.git",
  "modified": "2021-02-08T23:00:47+03:00"
}
```

This file (named `gitinfo.json by default) used later for

* embedding with filesystems
* showing project metagata

## Install

```sh
go get github.com/pgmig/gitinfo/...
```

## Usage

### Create gitinfo.json

Run go:generate just before embedding:

```go
// Generate gitinfo.json
//go:generate gitinfo ../../html

// Generate resource.go by [parcello](github.com/phogolabs/parcello)
//go:generate parcello -q -r -d ../../html
```

### Read gitinfo.json

Read metadata from .gitinfo.json if it exists, fetch from git otherwise

```go
var gi gitinfo.GitInfo
err = gitinfo.New(log, cfg).Make("cmd/", &gi)
```

### Generate gitinfo.json for single dir

```go
//go:generate gitinfo dir
```

### Generate gitinfo.json files for dir/*/ dirs

Used when dir contains git submodules

```go
//go:generate gitinfo dir/*
```

## License

The MIT License (MIT), see [LICENSE](LICENSE).

Copyright (c) 2019-2021 Aleksey Kovrizhkin <lekovr+pgmig@gmail.com>
