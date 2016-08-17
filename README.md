# Health Monitor

A standalone app to monitor and control OWTF - written in Go.


## Install

1. Clone or download and extract the repository to a folder, `health_monitor`, in `$GOPATH/src/github.com/owtf`. (You can also use `go get github.com/owtf/health_monitor`.)
2. Build a static binary by `go build -i` which creates an executable file, `health_monitor` in the project root.
3. Run the executable by `./health_monitor <flag>`. For help, run `./health_monitor --help`.


## Developement

1. First install **Go** and setup a development environment following the [this guide](https://golang.org/doc/install).
2. Clone the repository to `$GOPATH/src/`. (do not use `go get`!)
3. Make changes, and run `go run healthmon.go` to see your changes without explicitly building a new binary. When done, build a new binary using `go build` and test your changes.
4. Send a PR!

### Install script

Run ` bash -c "$(wget -c https://raw.githubusercontent.com/owtf/health_monitor/master/setup.sh -O -)"` to setup health monitor for development purpose. To install golang and setup its environment variables use `-g` option in above command.

> The logs for the monitor are stored in `$HOME/.owtf_monitor`. Until the web interface or CLI is completed, the working can be monitored by the logs.


## LICENSE

See [LICENSE](LICENSE)
