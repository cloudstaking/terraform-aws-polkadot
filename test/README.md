# Infrastructure tests

When changes happen, it is important to run infrastructure tests in order verify that some of the important things are still working as expected:

- Disk size 
- Tools installed (docker, docker-compose)
- etc

## How to run the tests

```sh
$ go test -v -timeout 30m
...
```
