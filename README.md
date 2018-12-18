# Parallel

This aims to be a replacement for the GNU parallel utility, written in Go.

## Installation

Installation is simple:

```bash
go get github.com/mylanconnolly/parallel
```

## Usage

Usage is similar to GNU parallel:

```bash
find /etc -type f | parallel md5sum            # Example to MD5 sum all files in /etc, using all cores
find /etc -type f | parallel -j 2 md5sum       # Example to MD5 sum all files in /etc, using 2 cores
find /etc -type f -print0 | parallel -0 md5sum # Like example 1, except using ASCII NULL delimiter, instead of newlines
parallel -a lines.txt md5sum                   # Like example 1, except using lines.txt as input, instead of stdin.
```

## TODO

The following items are still outstanding:

- [ ] Add tests
- [ ] Cover more use-cases that GNU parallel goes over
- [ ] Preserve output streams (stdout / stderr)
- [ ] Add shell completion
