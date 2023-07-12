# Add Atmel/Microchip SAM Devices

Go to http://packs.download.atmel.com/ and download the
appropriate .atpack file into this directory.

Then run `go run gen-from-atpack FILE_NAME_HERE.atpack`.

This will produce extract the .svd files from the
atpack and run gen-device-svd.go to generate the various
constants which can then be used in `targets` to support
a new device.

These new generated files in mcus/sam/[devicename] should
be committed to the Git repository (assuming supported is being
added to goocd), but the .atpack or other manually downloaded
resources should not.
