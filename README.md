# Parallel

This is meant to be a replacement for GNU Parallel written in Go. This started
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

Here are some benchmarks using the `time` command:

Below is the timing for the GNU version:

```
$ time find ~/src/go -type f | parallel md5sum
... output elided ...
parallel md5sum  61.22s user 44.63s system 286% cpu 36.896 total
```

Below is the timing for this version:

```
$ time find ~/src/go -type f | ./parallel md5sum
... output elided ...
./parallel md5sum  8.74s user 3.53s system 669% cpu 1.832 total
```

This represents a total execution time that is roughly 20x faster.

A few notes on my test environment:

- Intel Core i7 8700K
- 64GB of RAM
- 512GB Samsung 960 Pro NVMe SSD
- Ubuntu 18.04 in WSL2 on Windows 10 Pro Insider Preview 19603.rs_prerelease.200403-1523

## TODO

On both of my machines (One 6-core Core i7-8700K and one 4-core Ryzen 7 2700U)
performance seems to peak at a concurrency level of 4. I would like to hunt down
the cause of this bottleneck.

GNU Parallel supports building pipelines in its templating language. I would
like to emulate this, but I feel like it would add a fair amount of complexity.
