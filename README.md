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
```sh
curl -L https://raw.githubusercontent.com/g4s8/wts/master/scripts/download.sh | sh
```

On MacOS you can install it using `brew` tool:
```sh
brew tap g4s8/.tap https://github.com/g4s8/.tap
brew install g4s8/.tap/wts
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

Use this token for authentication: `wts -token yourtoken`
<!--
@todo #1:30min Read configuration file `~/.config/wts/config.yaml`
 where API token can be located. CLI `-token` parameter
 should override config token on conflicts.
-->

Get wallet id: `wts -token yourtoken id`

Get walle balance: `wts -token yourtoken balance`
