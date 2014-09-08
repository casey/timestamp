timestamp
---------

A very simple timestamp service.

[Test instance here.](http://rodarmor-timestamp.appspot.com)

A KEY matches `/[a-zA-Z0-9.-_]*/`.
A TIME matches `/[0-9]+[.][0-9]+/` and is in seconds since the epoch.

PUT /KEY establishes a timestamp for a key. Subsequent puts after the first are ignored.
GET /KEY retrieves the timestamp of a key.
