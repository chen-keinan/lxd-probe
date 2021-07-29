9[![Go Report Card](https://goreportcard.com/badge/github.com/chen-keinan/lxd-probe)](https://goreportcard.com/report/github.com/chen-keinan/lxd-probe)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/chen-keinan/lxd-probe/blob/main/LICENSE)
<br><img src="./pkg/img/lxd-gopher.png" width="300" alt="lxd-probe logo"><br>
# lxd-probe

###  Scan your Linux container (LXD / LXC) runtime !!
Lxd-Probe is an open source audit scanner who perform audit check on a linux container manager and output a security report.

The audit tests are the full implementation of [CIS Lxd Benchmark specification](https://www.cisecurity.org/benchmark/lxd/) <br>

audit result now can be leveraged as webhook via user plugin(using go plugin)
#### Audit checks are performed on linux containers, and output audit report include :
 1.  root cause of the security issue.
 2. proposed remediation for security issue

--------------------------------------------------------------------------------------------------------

* [Installation](#installation)
  
* [Quick Start](#quick-start)

## Installation

```
git clone https://github.com/chen-keinan/lxd-probe
cd lxd-probe
make build
./lxd-probe
```

Note : lxd-probe require privileged user to execute tests.

## Quick Start

```
Usage: lxd-probe [--version] [--help] <command> [<args>]

Available commands are:
  -r , --report :  run audit tests and generate failure report
  -i , --include:  execute only specific audit test,   example -i=1.2.3,1.4.5
  -e , --exclude:  ignore specific audit tests,  example -e=1.2.3,1.4.5
```
