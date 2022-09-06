# delim
Creates a visual line to aid in reading terminal screens

## Installation

You can install a few ways:

1. Download the binary for your OS from https://github.com/smford/delim/releases
1. or use `go install`
   ```
   go install -v github.com/smford/delim@latest
   ```
1. or clone the git repo and build
   ```
   git clone git@github.com:smford/delim.git
   cd delim
   go get -v
   go build
   ```
1. or by brew:
   ```
   brew install smford/tap/delim
   ```


## Running

Simply run the command `delim`!


## Configuration

No configuration is really needed, however if you wish to change the character used in the line from `=` to another character or even a string, this can be done a few ways.

Example methods to change the line from using `=` to `Z`:
1. By command line:
   
   `delim --char Z`
1. By configuration file:
   ```bash
   cat ~/.delim
   char: Z
   ```
1. By environment variable:
   
   `export DELIM_CHAR="Z"`

If you are feeling adventurous you can set it to being a string rather than a single char, for example `export DELIM_CHAR="carrot"`
