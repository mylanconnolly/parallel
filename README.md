# Parallel

This is meant to be a replacement for GNU parallel written in Go. This started
as a learning exercise in dealing with parallelism in Go, but has since become
a tool that I regularly use.

The tool will start a worker for each CPU and work through the list of jobs that
you give it. The amount of workers is configurable.

This tool is striving to only use stdlib packages.

## Usage

Install using `go get github.com/mylanconnolly/parallel` or some other means.

### Simple usage

The most straightforward usage would be:

```shell
# Want to calculate the MD5 sum of every file in /etc?
$ find /etc -type f | parallel md5sum

# Want to only use two workers for the same thing?
$ find /etc -type f | parallel -j 2 md5sum
```

### Command templating

You can utilize Go templates when performing a command using the `-t` flag. When
using the `-t` flag, you do not need to specify the command (it will be ignored
if you do).

The following fields are available when using templates:

| Field        | Definition                                                    |
| :----------- | :------------------------------------------------------------ |
| `{{.Cmd}}`   | The path of the command specified, for example echo or md5sum |
| `{{.Input}}` | The current input that we received via stdin or input file    |
| `{{.Start}}` | The time that parallel was started                            |
| `{{.Time}}`  | The time that the current operation began                     |

In addition, the following functions are available in templates:

| Function       | Help                                   |
| :------------- | :------------------------------------- |
| `toUpper`      | Transform the string to uppercase      |
| `toLower`      | Transform the string to lowercase      |
| `absolutePath` | Get the absolute path of a filename    |
| `basename`     | Get the basename of a file path        |
| `dirname`      | Get the directory of a file path       |
| `ext`          | Get the extension of a file            |
| `noExt`        | Get the file path without an extension |

Some examples below:

```shell
# Copy some files up a level (utilizing template pipelines).
parallel -a ./files.txt -t 'cp {{.Input}} {{.Input | dirname | dirname}}'

# Create a directory named after the file (without extension).
parallel -a ./files.txt -t 'mkdir -p {{.Input}} {{noExt .Input}}'

# Echo the base name of the file without the extension (utilizing template
# pipelines).
parallel -a ./files.txt -t 'mkdir -p {{.Input}} {{.Input | basename | noExt}}'
```

For more general information about Go templates, check
[here](https://golang.org/pkg/text/template/#pkg-overview).

## Real world examples

Here are some benchmarks using the `time` command. The benchmark I put together
is to run `md5sum` for every file in the Go source repository as of commit
14bec27743.

Below is the timing for the GNU version:

```
$ time find ~/src/go -type f | parallel md5sum > /dev/null
noglob find ~/src/go -type f  0.01s user 0.07s system 0% cpu 22.580 total
parallel md5sum > /dev/null  22.65s user 42.48s system 246% cpu 26.432 total
```

Below is the timing for this version:

```
$ time find ~/src/go -type f | ./parallel md5sum > /dev/null
noglob find ~/src/go -type f  0.02s user 0.05s system 3% cpu 1.845 total
./parallel md5sum > /dev/null  7.46s user 2.72s system 396% cpu 2.569 total
```

In this example it took GNU parallel around 10 times longer to complete the same
amount of work.

A few notes on my test environment:

- Thinkpad A485
- AMD Ryzen Pro 2700U
- 16GB of RAM
- 256GB NVMe SSD (though I believe it might be a pretty low-quality one)
- Ubuntu 20.04 LTS (kernel version 5.4.0-21-generic)
