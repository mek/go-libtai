# TODO

- Adding leapsecs table for accurate TAI <-> UTC conversion.
- Support for reading leapseconds from standard system files (e.g., `/usr/share/zoneinfo/leap-seconds.list`).
- API for converting between `time.Time` and `tai.TAI` / `taia.TAIA` with leap second awareness.
- Create simple `tai` log script (second-level precision only).
- Create decoder scripts (`tai64local`, `tai64nlocal`) to convert timestamps back to local time.
