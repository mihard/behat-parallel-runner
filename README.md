## Behat parallel runner [![Build Status](https://travis-ci.org/mihard/behat-parallel-runner.svg?branch=master)](https://travis-ci.org/mihard/behat-parallel-runner)

### What is it?

It is a simple launcher of behat tests, which allow:

 - to run each test in a separate process
 - to run several test in parallel

### How to use

- Put the executable file into the root of testable project (on the same level with `vendor` directory)
- Run following command

```
$ bpr 2 -s my_suite
```
where `bpr` is the name of executable file, 2 is a number of threads and everything else can be any regular behat arguments.

The executable binary you can find on the releases page.

### How to build 

- Clone the repository
- Run following command

```
$ cd <clonned project root>
$ GOPATH=<clonned project root> go build -i "-ldflags=-linkmode internal" -o ./build/bpr github.com/mihard/behat-parallel-runner
```

- Build on OSX machine for linux
```
$ env GOOS=linux GOARCH=amd64 GOPATH=<clonned project root> go build -i "-ldflags=-linkmode internal" -o ./build/bpr github.com/mihard/behat-parallel-runner
```