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
where `bpr` is the name of executable file, 2 is a number of threads and everything else can be any regular behat arguments

