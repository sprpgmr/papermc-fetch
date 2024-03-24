# papermc-fetch

papermc-fetch is a simple command line utility written in Go to download releases of the papermc minecraft server.

If you already have papermc downloaded, papermc-fetch will compare it's checksum to the checksum of the latest build found against your local version, and will skip the download if they're the same, so this can be used to efficiently check for updates on a schedule.

## Basic usage

```shell
# Download the latest stable build
./papermc-fetch

# Download the latest build, including experimental builds
./papermc-fetch --experimental

# Just check for builds, don't download
./papermc-fetch --skip-download

# Download the latest build for Minecraft 1.22.4
./papermc-fetch --prefix 1.22.4

# Download the latest build for Minecraft 1.22.X
./papermc-fetch --prefix 1.22

# Download the latest build and name the file papermc123.jar
./papermc-fetch --file papermc123.jar
```

## Sample Output:

Check for updates without downloading:
```text
./papermc-fetch --skip-download
Checking for latest version of paper...
Latest paper version is 1.20.4 - build #461
```

Download latest version:
```text
./papermc-fetch
Checking for latest version of paper...
Latest paper version is 1.20.4 - build #461
Downloading...
Finished downloading.
Verifying file integrity...
Download verified.
```

Download latest version (latest version already downloaded):
```text
./papermc-fetch
Checking for latest version of paper...
Latest paper version is 1.20.4 - build #461
You already have this version of paper.
```

## Compiling:
Make sure you have Go 1.21.5 or later installed, then run the commands below in the cloned repo:
```shell
go mod download
go build
```
