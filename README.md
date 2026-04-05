# Go LibTAI

`go-libtai` is a Go implementation of the TAI64 and TAIA64 time formats 
originally specified by D. J. Bernstein in [libtai](https://cr.yp.to/libtai.html).

TAI (International Atomic Time) provides a linear time scale without the 
discontinuities of leap seconds, making it ideal for high-precision 
timestamping and interval calculations.

## Contents

- `lib/tai`: Core TAI64 implementation (second precision).
- `lib/taia`: TAIA64 implementation (attosecond precision).
- `cmd/tai_now`: Print the current TAI64 time.
- `cmd/taia_now`: Print the current TAIA64 time.
- `cmd/tai64n`: Filter to prepend TAI64N timestamps to each line of input.

## Building

A `Makefile` is provided for common tasks:

```bash
make all        # Build all command-line tools into bin/
make lib-test   # Run library unit tests
make check      # Run linting and formatting checks
```

## License

Like the original `libtai`, this project is released into the Public Domain.
