# notorious

`notorious` searches `stdin` for lines of text matching a pattern and outputs them to stdout.

## Installation

```
$ go install notorious.go
```

## Usage

`notorious` runs the program with the specified options.

```
$ notorious [OPTIONS] PATTERN
```

### get help

```  
$ notorious --help
```

### filter the output of stderr stream by appending the lines that contain ERROR to a file

```
$ myprogram> notorious ERROR >> myprogram_errors.log
```

### search for the word "log" in notorious.go, case-insensitive

```
$ notorious -i log < notorious.go
```
