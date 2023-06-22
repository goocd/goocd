# goocd

GoOCD is some initial work on an "on chip debugging" in the style of OpenOCD or PyOCD, but written in Go.

In other words, it's a simple way, written in Go, to load code onto Microcontrollers.

## Why?

GoOCD attempts to address these concerns with OpenOCD (and PyOCD to some extent):

* **Cross platform support:** Compiling OpenOCD for different platforms is, much as any C program is, a pain.  Managing Python versions and package installation is only marginally better. The build workflow using Go is significantly simpler and more reliable, and runs on Windows, Linux, Mac.

* **Code Organization:** Managing a large code base that supports so many chips is difficult, and tends to result in a very non-obvious set of interrelated files that are difficult to understand and even more difficult to extend.  We wanted something where the functionality for a chip was just simple and obvious and in one place.  Take a look at main.go and the files in that directory and tell us what you think.

## Usage

Git clone, cd into the directroy and then `go install .`

Make sure `~/go/bin` is in your path. (Or you can customize the output path using the "go build" command instead.)

Run with `goocd -h` to get help.

Yes, that's it.

## Examples

TODO: show examples of loading code and whatnot.

## Supported Chips/Debuggers

TODO: show where to look for chips

Initial work is being done to support SAMD51/SAME51 and SAML10 via CMSIS-DAP (e.g. Atmel-ICE).

