# TAI and the TAI64/TAIA64 Formats

Temps Atomique International (TAI) is a high-precision, linear time scale 
based on the vibrations of cesium atoms. Unlike UTC (Coordinated 
Universal Time), TAI does not include leap seconds. This makes it 
ideal for interval arithmetic and high-performance applications where 
time must never flow backwards or experience "stops."

## The TAI64 Format

TAI64 is a 64-bit format designed by D. J. Bernstein to represent a 
moment in time with second-level precision.

- **Storage:** 8-byte big-endian unsigned integer.
- **Epoch:** The TAI64 epoch is 2^62 seconds before 1970-01-01 00:00:10 TAI.
- **Range:** It covers a span of hundreds of billions of years, far 
  exceeding the expected lifetime of the Earth.

A TAI64 value of `2^62` corresponds to 1970-01-01 00:00:00 TAI.

## High-Precision Extensions

The `libtai` specifications include extensions for sub-second precision:

### TAI64N (Nanosecond Precision)
Consists of 12 bytes:
- 8 bytes: TAI64 seconds.
- 4 bytes: Nanoseconds (0 to 999,999,999).

### TAIA64 (Attosecond Precision)
Consists of 16 bytes:
- 8 bytes: TAI64 seconds.
- 4 bytes: Nanoseconds.
- 4 bytes: Attoseconds (0 to 999,999,999).

One attosecond is 10^-18 seconds.

## Leap Seconds and Conversion

The difference between TAI and UTC changes whenever a leap second is 
inserted. As of January 2026, TAI is 37 seconds ahead of UTC.

Conversion between TAI and Unix time (UTC-based) requires knowing the 
history of leap seconds. The `go-libtai` library provides basic 
offset-based conversion but does not currently include a full leap 
second table.

## References

The TAI64, TAI64N, and TAIA64 formats were originally specified by 
D. J. Bernstein. Further details can be found at:

- https://cr.yp.to/libtai/tai64.html
- https://cr.yp.to/libtai.html
