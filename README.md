# Usage

When logging into the commandline using single-sign-on, you typically need to
click the copy button, open your credentials file, paste the contents
in the appropriate location, and then potentially change the profile name.

pcreds is a simple commandline utility to do all of these things
with one alias-able command.

pcreds:

* Reads a credentials profile from the clipboard
* Parses the profile name and looks up an alias in pcreds.yaml (stored in your .aws folder)
* Parses your credentials file, updates the corresponding profile and saves an updated version

# Installation

Install Go from the [official website](https://go.dev/)

clone this repository and build the executable. Then move it to your bin folder

```bash
git clone https://github.com/DaraDadachanji/pcreds.git
cd pcreds
go build
mv ./pcreds /usr/local/bin/pcreds
```
