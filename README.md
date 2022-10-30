# Castle
*So far it is absolutely nothing*

<img src="graphics/logo.png" width="300" alt="Logo">

## Dependencies
`GNU make, go ^1.19`

## Building castle
`make` - Builds castle for local OS

`make windows` - Builds castle for windows

`make install` - Builds castle and copies to `/usr/local/bin/castle`. Linux only.
Requires sudo.

`make clean` - Removes dist directory and `/usr/local/bin/castle`.

## Testing
`make test`