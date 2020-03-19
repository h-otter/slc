# Super Lightweight Container

- This provide super lightweight container, only filesystem isolation.
  - No isolation of network, user, IPC, UTS and process namespaces with host
  - Executable Ping (with setuid / capability) on rootless container
  - No daemon

## Usage

### As a library

```go
	c, err := container.NewClient("./state")
	if err != nil {
		return err
	}

	if err := c.Run("alpine", []string{"/bin/bash"}); err != nil {
		return err
	}
```

### As a command

```sh
sudo slc pull busybox
slc run busybox ping 8.8.8.8
sudo slc rm busybox
sudo slc clear
```

## Build

```sh
git clone git@github.com:h-otter/slc.git
make build
```

## Benchmark (TBD)

- Container overhead
  - Execution time
  - Memory
  - CPU

### Execution time

- Execution time is benchmarked with https://github.com/gsauthof/utility/blob/master/benchmark.sh.
  - https://unix.stackexchange.com/a/313714

```sh
% time docker run debian echo test
test
docker run debian echo test  0.04s user 0.03s system 8% cpu 0.777 total

% time docker start -a 9e184a0bb079
test
docker start -a 9e184a0bb079  0.04s user 0.05s system 11% cpu 0.854 total
```


```sh
% time ./main echo test
test
./main echo test  0.00s user 0.01s system 110% cpu 0.010 total
```
