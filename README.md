# GoCronScheduler

GoCronScheduler is proof of concept application built using Golang (Go) to test the capabilities of creating and monitoring cron jobs in Kubernetes (k8s) in a dynamic manner.

## Installation

To install GoCronScheduler, you need to have Golang installed on your system. Then you can install it using `go get`:

```bash
go get github.com/copito/gocronscheduler
```

> The project does require a kubernetes instance like Minikube, K3s, K8s or others running and therefore a config file `.kube/config` located in the home path - with connection to said instance.

## Usage

GoCronScheduler is meant to be used as base code to build more interesting projects that require a cron scheduler to be leveraged, but it needs to be external to the current application (like k8s).

As mentioned before it is a POC (proof of concept) and therefore should not be used as-is in any projects as it does not contain any specific tests and is not ready for any production use case.

## Development

The code was developed using vscode and therefore contain a [launch file](./.vscode/launch.json) that will point directly to the [main.go file](./scheduler/main.go) for easy debugging visibility.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License - see the LICENSE [file](./LICENSE) for details.
Feel free to customize it further based on your project's specifics!
