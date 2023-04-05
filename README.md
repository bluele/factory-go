# factory-go

![Test](https://github.com/bluele/factory-go/workflows/Test/badge.svg)
[![GoDoc](https://godoc.org/github.com/bluele/factory-go?status.svg)](https://pkg.go.dev/github.com/bluele/factory-go?tab=doc)

factory-go is a fixtures replacement inspired by [factory_boy](https://github.com/FactoryBoy/factory_boy) and [factory_bot](https://github.com/thoughtbot/factory_bot).

It can be generated easily complex objects by using this, and maintain easily those objects generaters.

## Install

```
$ go get -u github.com/hyuti/factory-go/factory
```

## Example

* [Define a simple factory](https://github.com/hyuti/factory-go#define-a-simple-factory)
* [Define a factory which has input type different from output type](https://github.com/hyuti/factory-go#define-a-factory-which-has-input-type-different-from-output-type)
* [Use factory with random yet realistic values](https://github.com/hyuti/factory-go#use-factory-with-random-yet-realistic-values)
* [Define a factory includes sub-factory](https://github.com/hyuti/factory-go#define-a-factory-includes-sub-factory)
* [Define a factory includes a slice for sub-factory](https://github.com/hyuti/factory-go#define-a-factory-includes-a-slice-for-sub-factory)
* [Define a factory includes sub-factory that contains self-reference](https://github.com/hyuti/factory-go#define-a-factory-includes-sub-factory-that-contains-self-reference)
* [Define a sub-factory refers to parent factory](https://github.com/hyuti/factory-go#define-a-sub-factory-refers-to-parent-factory)
* [Define a factory has input type different from output type](https://github.com/hyuti/factory-go#define-a-sub-factory-refers-to-parent-factory)

## Features

## Roadmap
- 🚧️ Add Features section for README
- 🚧️ Bulk create feature
- 🚧️ Bulk update feature
- 🚧️ Bulk delete feature

## Persistent models

Here is a list of integration examples with some popular ORM libraries (Please pay atttention this is not an official integration):
- [gorm](https://github.com/hyuti/factory-go/blob/master/examples/gorm_integration.go) (what is [gorm](https://github.com/jinzhu/gorm))
- [ent](https://github.com/hyuti/factory-go/blob/master/examples/integration-with-ent/ent/factory.go) (what is [ent](https://github.com/ent/ent))

# Contributors
The [original repo](https://github.com/bluele/factory-go) of this project has not been actived for a long time, so this is a fork and maintained by me. Any contributions will be appreciated.

# Author

**Jun Kimura**
* <junkxdev@gmail.com>
