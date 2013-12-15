# argf

This is a small Go library that reads lines of text from either a file or files specified in `os.Args`, or, if
none are given, from stdin. It's similar to Ruby's `ARGF` or Perl's diamond operator.

This is mainly intended for tiny script-like CLI tools. Note that multiple goroutines should not call its
functions concurrently.
