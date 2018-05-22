# Axis
> CLI tool for Axis Communications cameras

## Examples

List all recordings on the camera:
```console
$ axis record ls
```

Download all recordings that occur between 7AM and 6PM with a duration under 20 seconds:

```console
$ axis record pull -d -m 20 ./
```

## Installation

```console
go get github.com/magahet/axis
```

## Config

Uses Viper. Default config path is `$HOME/.axis.yaml`

Config example:

```yaml
host: mycam.example.com
verbose: true
```