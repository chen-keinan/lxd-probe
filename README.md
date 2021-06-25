[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/chen-keinan/lxd-probe/blob/main/LICENSE)
<br><img src="./pkg/img/lxd-gopher.png" width="300" alt="lxd-probe logo"><br>
# lxd-probe

###  Scan your Linux container runtime !!
Lxd-Probe is an open source audit scanner who perform audit check on a linux container manager and output a security report.

The audit tests are the full implementation of [CIS Lxd Benchmark specification](https://www.cisecurity.org/benchmark/lxd/) <br>

audit result now can be leveraged as webhook via user plugin(using go plugin)
#### Audit checks are performed on linux containers, and output audit report include :
* root cause of the security issue
* proposed remediation for security issue