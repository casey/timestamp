timestamp
---------

A very simple timestamp service.

A KEY matches `/[a-zA-Z0-9.-_]*/`.

A TIME is in seconds since the epoch and matches `/[0-9]+[.][0-9]+/`.

PUT /KEY establishes a timestamp for a key. Subsequent puts to a given key after the first are ignored.

GET /KEY retrieves the timestamp of a key.
