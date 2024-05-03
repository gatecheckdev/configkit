# ConfigKit

A Go Library for managing CLI configurations.

## Why is this?

If you make a lot of CLI tools, you likely have some pattern you like to follow for managing a config file.
This project is an attempt to capture a common pattern and make it reusable across a variety of projects.

It is intended to be used with the fantastic [Cobra CLI Library](https://github.com/spf13/cobra).

## Features

The library is still under active development and isn't recommended in a production environment.

Below are the intended features planned for the library:

- [X] Cobra Flag binding strategy
- [X] Environment Variable binding strategy
- [X] Value retrival based on common sense order of precedence evaluation

## Examples

Code examples can be found in the `examples/` directory of this project.
