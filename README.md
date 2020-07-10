[![GoDoc](https://godoc.org/github.com/giantswarm/helmclient?status.svg)](http://godoc.org/github.com/giantswarm/helmclient) [![CircleCI](https://circleci.com/gh/giantswarm/helmclient.svg?&style=shield)](https://circleci.com/gh/giantswarm/helmclient)

# helmclient

Package helmclient implements [Helm] related primitives to work against helm
releases. Currently supports Helm 3.

## Branches

- `master`
    - Latest version using Helm 3.
- `helm2`
    - Legacy support for Helm 2.

## Interface

See `helmclient.Interface` in [spec.go] for supported methods.

## Getting Project

Clone the git repository: https://github.com/giantswarm/helmclient.git

### How to build

Build it using the standard `go build` command.

```
go build github.com/giantswarm/helmclient
```

## Contact

- Mailing list: [giantswarm](https://groups.google.com/forum/!forum/giantswarm)
- IRC: #[giantswarm](irc://irc.freenode.org:6667/#giantswarm) on freenode.org
- Bugs: [issues](https://github.com/giantswarm/helmclient/issues)

## Contributing & Reporting Bugs

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches, the
contribution workflow as well as reporting bugs.

## License

helmclient is under the Apache 2.0 license. See the [LICENSE](LICENSE) file
for details.

[Helm]: https://github.com/helm/helm
[spec.go]: https://github.com/giantswarm/helmclient/blob/master/spec.go
