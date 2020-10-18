package opts

import (
	"errors"
	"flag"
	"fmt"
	"local/notorious/logging"
	"os"
	"regexp"
	"strings"
)

// raw CLI flags, only used to create an Opts during Parse()
var (
	after       = flag.Int("A", 0, "how many lines of context [A]fter the match to print")
	before      = flag.Int("B", 0, "how many lines of context [B]efore the match to print")
	context     = flag.Int("C", 0, "how many lines of [C]ontext around the match to print")
	ignoreCase  = flag.Bool("i", false, "[i]gnore case in matches")
	lineNumbers = flag.Bool("n", false, "whether to print the line [n]umbers")
	literal     = flag.Bool("e", false, "match using string lit[e]rals instead of regular expressions")
	posix       = flag.Bool("posix", false, "use [posix] regular expresisons")
	verbose     = flag.Bool("v", false, "[v]erbose debug info")
)

// Opts represent the parsed and validated options. Use these instead of the command-line flags directly
// You need to set the Matches func somehow; New() is a good way to do it
type Opts struct {
	// Lines of context to match. Default is 0, 0.
	Context struct{ Before, After int }
	// Whether to print line numbers. Default is 0.
	LineNumbers bool
	// A function to match a line of text.
	Matches func(string) bool
	// Print verbose debug output.
	Verbose bool
}

// create a new Opts with everything set to the default
func New(pattern string) (Opts, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return Opts{}, fmt.Errorf("could not compile regexp from %q: %v", pattern, err)
	}
	return Opts{Matches: re.MatchString}, nil
}

// these builder methods are mostly here to make testing over in main_test.go a little easier and cleaner, and you could easily omit them

// WithMatcher returns a copy with Matches set to f
func (o Opts) WithMatcher(f func(s string) bool) Opts {
	o.Matches = f
	return o
}

// WithContext returns a copy with Context.Before and Context.After set to before and after
func (o Opts) WithContext(before, after int) Opts {
	o.Context.Before, o.Context.After = before, after
	return o
}

// WithLineNumbers sets LineNumbers to b
func (o Opts) WithLineNumbers(b bool) Opts {
	o.LineNumbers = b
	return o
}

// WithVerbose sets verbose to b
func (o Opts) WithVerbose(b bool) Opts {
	o.Verbose = b
	return o
}

// parse and validate the command-line flags
func Parse() (o Opts, err error) {
	if flag.Parsed() {
		return Opts{}, errors.New("flags already parsed")
	}
	flag.Parse()
	logger := logging.Debug(os.Stderr, *verbose)

	pattern := flag.Arg(0)
	if pattern == "" {
		return Opts{}, errors.New("expected an argument PATTERN")
	}
	if len(flag.Args()) > 1 {
		return Opts{}, errors.New("notorious does not (yet) support more than one positional argument. If you set flags, they need to go before the positional arguments, not after")
	}
	var lineCtx struct{ Before, After int }
	switch {
	case *context < 0:
		return Opts{}, fmt.Errorf("flag -C must be nonnegative, but got %d", *context)
	case *before < 0:
		return Opts{}, fmt.Errorf("flag -B must be nonnegative, but got %d", before)
	case *after < 0:
		return Opts{}, fmt.Errorf("flag -A must be nonnegative, but got %d", after)
	case *context != 0 && *before != 0:
		return Opts{}, errors.New("flags -B and -C are mutually exclusive")
	case *context != 0 && *after != 0:
		return Opts{}, errors.New("flags -A and -C are mutually exclusive")
	case *context != 0:
		lineCtx = struct{ Before, After int }{*context, *context}
	default:
		lineCtx = struct{ Before, After int }{*before, *after}
	}
	logger.Printf("%#+v", lineCtx)

	var matcher func(s string) bool
	switch {
	case *ignoreCase && *literal:
		logger.Print("mode: case-insenstive literal")
		matcher = func(text string) bool { return strings.EqualFold(text, pattern) }
	case *literal:
		logger.Print("mode: literal")
		matcher = func(text string) bool { return text == pattern }
	case *ignoreCase && *posix:
		logger.Print("mode: case-insensitive posix")
		re, err := regexp.CompilePOSIX("(?i)" + pattern)
		if err != nil {
			return Opts{}, fmt.Errorf("could not compile case-insensitive POSIX regexp from pattern %q: %v", pattern, err)
		}
		matcher = re.MatchString
	case *posix:
		logger.Print("mode: posix")
		re, err := regexp.CompilePOSIX(pattern)
		if err != nil {
			return Opts{}, fmt.Errorf("could not compile POSIX regexp from pattern %q: %v", pattern, err)
		}
		matcher = re.MatchString
	case *ignoreCase:
		logger.Print("mode: case-insensitive regexp")
		re, err := regexp.CompilePOSIX("(?i)" + pattern)
		if err != nil {
			return Opts{}, fmt.Errorf("could not compile case-insensitive POSIX regexp from pattern %q: %v", pattern, err)
		}
		matcher = re.MatchString
	default:
		logger.Print("mode: regexp")
		re, err := regexp.Compile(pattern)
		if err != nil {
			return Opts{}, fmt.Errorf("could not compile regexp from %q: %v", pattern, err)
		}
		matcher = re.MatchString
	}
	o = Opts{
		Matches:     matcher,
		Context:     lineCtx,
		LineNumbers: *lineNumbers,
		Verbose:     *verbose,
	}
	return o, nil
}
