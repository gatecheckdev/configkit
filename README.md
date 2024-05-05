# ConfigKit

A Go Library for managing CLI configurations.

## Why is this?

If you make a lot of CLI tools, you likely have some pattern you like to follow for managing a config file.
This project is an attempt to capture a common pattern and make it reusable across a variety of projects.

It is intended to be used with the fantastic [Cobra CLI Library](https://github.com/spf13/cobra).

Scenario:

You're building a CLI based application and you want to support arguments from the following sources:

- cobra pFlags
- Environment Variables
- Config Files in json, yaml, or toml

Your options are:

1. Build the logic and order-of-precedence evaluation from scratch
2. Use something like Viper

While number 1 isn't difficult, it's easy to have bugs sneak in, even if you've done it a few times.
Viper is a great tool to use as well but some will find that a lot happens behind the scenes to accomodate a wide range of use cases.
It's super flexible but with flexibility, you get a lot of leadway to shoot yourself in the foot.

This package is intended to give you a way to define a "meta config", a mirrored config file with additional metadata
that can be used at runtime to bind flags, defaults, usage, env keys, and other options which, ideally, would elimate
a lot of teadious but necessary boiler plate.

## Features

The library is still under active development and isn't recommended in a production environment.

Below are the intended features planned for the library:

- [X] Cobra Flag binding strategy
- [X] Environment Variable binding strategy
- [X] Value retrival based on common sense order of precedence evaluation
- [ ] Config File strategy

## Examples

Code examples can be found in the `examples/` directory of this project.
