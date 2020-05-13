# Parallel

This is meant to be a replacement for GNU parallel written in Go. This started
as a learning exercise in dealing with parallelism in Go, but has since become
a tool that I regularly use.

The tool will start a worker for each CPU and work through the list of jobs that
you give it. The amount of workers is configurable.

This tool is striving to only use stdlib packages.

## Usage

Install using `go get github.com/mylanconnolly/parallel` or some other means.

The most straightforward usage would be:

```shell
# Want to calculate the MD5 sum of every file in /etc?
$ find /etc -type f | parallel md5sum

# Want to only use two workers for the same thing?
$ find /etc -type f | parallel -j 2 md5sum
```

Let's say you want to do something more complex... maybe you want to use some
templating?

```shell
# This will take every file in /etc and copies it with an extension of today's
# date. (please don't run this command)
$ find /etc -type f | parallel -t 'cp {{.Input}} {{.Input}}.{{.Start.Format "20060101"}}'
```

Maybe you need to use an input file as your source:

```shell
# Maybe you want to calculate the MD5 sum of all the files in a text file.
parallel -a ./files.txt -t 'md5sum {{.Input}}'
```

The following fields are available when using templates:

| Field   | Definition                                                    |
| :------ | :------------------------------------------------------------ |
| `Cmd`   | The path of the command specified, for example echo or md5sum |
| `Input` | The current input that we received via stdin or input file    |
| `Start` | The time that parallel was started                            |
| `Time`  | The time that the current operation began                     |

For more general information about Go templates, check
[here](https://golang.org/pkg/text/template/#pkg-overview).

## Real world examples:

Here are some benchmarks using the `time` command. The benchmark I put together
is to run `md5sum` for every file in the Go source repository as of commit
14bec27743.

Below is the timing for the GNU version:

```
$ time find ~/src/go -type f | parallel md5sum > /dev/null
... output elided ...
parallel md5sum > /dev/null  22.70s user 43.62s system 239% cpu 27.667 total
```

Below is the timing for this version:

```
$ time find ~/src/go -type f | ./parallel md5sum > /dev/null
... output elided ...
./parallel md5sum > /dev/null  8.06s user 3.11s system 333% cpu 3.344 total
```

In this example it took GNU parallel around 8 times longer to complete the same
amount of work.

A few notes on my test environment:

- Thinkpad A485
- AMD Ryzen Pro 2700U
- 16GB of RAM
- 256GB NVMe SSD (though I believe it might be a pretty low-quality one)
- Ubuntu 20.04 LTS (kernel version 5.4.0-21-generic)

## TODO

On both of my machines (One 6-core Core i7-8700K and one 4-core Ryzen 7 2700U)
performance seems to peak at a concurrency level of 4. I would like to hunt down
the cause of this bottleneck.

GNU parallel supports building pipelines in its templating language. I would
like to emulate this, but I feel like it would add a fair amount of complexity.
