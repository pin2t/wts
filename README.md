# WTS client

[![DevOps By Rultor.com](http://www.rultor.com/b/g4s8/wts)](http://www.rultor.com/p/g4s8/wts)

[![GitHub release](https://img.shields.io/github/release/g4s8/wts.svg?label=version)](https://github.com/g4s8/wts/releases/latest)
[![Build Status](https://img.shields.io/travis/g4s8/wts.svg?style=flat-square)](https://travis-ci.org/g4s8/wts)
[![CircleCI](https://circleci.com/gh/g4s8/wts.svg?style=svg)](https://circleci.com/gh/g4s8/wts)
[![Hits-of-Code](https://hitsofcode.com/github/g4s8/wts)](https://hitsofcode.com/view/github/g4s8/wts)

[![PDD status](http://www.0pdd.com/svg?name=g4s8/wts)](http://www.0pdd.com/p?name=g4s8/wts)
[![License](https://img.shields.io/github/license/g4s8/wts.svg?style=flat-square)](https://github.com/g4s8/wts/blob/master/LICENSE)

## Install

Download binary for your platform from Github releases:
https://github.com/g4s8/wts/releases/latest

Use shell script to get latest release binary (only Linux and MacOSx):
```sh
curl -L https://raw.githubusercontent.com/g4s8/wts/master/scripts/download.sh | sh
```

On MacOS you can install it using `brew` tool:
```sh
brew tap g4s8/.tap https://github.com/g4s8/.tap
brew install wts
```

Build from sources:
```sh
git clone https://github.com/g4s8/wts.git
cd wts
go build ./cmd/wts/
# target binary will be placed at $PWD/wts
```

## Usage

Login to [WTS](https://wts.zold.io/) and get [API token](https://wts.zold.io/api).
You may put this token to configuration file or add as command parameter explicitly.

Configuration file should be localted at `~/.config/wts/config.yml`:
```yaml
---
version: V1
wts:
  token: "...API token..."
  debug: true # show debug output, default false
  pull: true # pull wallet before each operation, default false
```

Usage: `wts [options] (argument)`

where options are:
 - `-token <token>` (required if not specified in the `config.yml`) - API token
 - `-debug` - show debug output (default `false`)
 - `-progress` - use `-progress=false` to hide progress spinner (default `true`)
 - `-pull` - pull wallet before operation
 - `-config` - use custom config file location (default `$HOME/.config/wts/config.yml`)

actions are:
 - `id` - print wallet id
 - `balance` - show balance
 - `txns` - show transaction list, additional options are:
   - `-filter` - regex filter
   - `-limit` - transaction limit
 - `stats` - show statistic for period
   - `-period` - days for statistic

