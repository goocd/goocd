# goocd

GoOCD is some initial work on an "on chip debugging" in the style of OpenOCD or PyOCD, but written in Go.

In other words, it's a simple way, written in Go, to load code onto Microcontrollers.

## Why?

GoOCD attempts to address these concerns with OpenOCD (and PyOCD to some extent):

* **Cross platform support:** Compiling OpenOCD for different platforms is, much as any C program is, a pain.  Managing Python versions and package installation is only marginally better. The build workflow using Go is significantly simpler and more reliable, and runs on Windows, Linux, Mac.

* **Each target is a simple function** that imports and executes exactly what is needed for that chip+debug probe.

* A corollary is **we don't have big abstractions that try to support e.g. every possible flash module (or other general feature) in the world,** instead we have specific packages that handle a particular thing, which are then imported by the specific targets that need them and are used directly.  This difference in approach allows us to make much more natural and appropriate abstractions when we observe them, rather than trying to make and maintain a general abstraction that tries to support every possible scenario.  If two chips have similar-enough NVM functionality, their implementation can be shared.  If another chip has totally different NVM functionality, then it's NVM support can be implemented completely independently and used just for that other target.

* You can **easily extend and customize by making a new program (main)** and copying the bits you need (see Make Your own Program below).  No scripting languages, just Go code.

* **General code organization concerns.**  While this is certainly subjective, we feel that the above changes in approach from what other projects have done make it a lot easier to see what's going on and do updates, fixes, add more chips/probes, etc.  Take a look at the `targets` directory and tell us if you like how we've done it. 

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

## Make Your own Program

TODO: explain how targets package is organized and that someone can just copy the bits they need out of a specific target (mention that you could copy main.go and do that too, it just depends on how much you need to reuse the target logic from various targets vs just making something specific to whatever chip+probe you're using).

