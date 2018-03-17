# main

## word count

[`wc.go`](wc.go) reports the number of occurrences of
each word in its input. A word is any contiguous sequence of letters, as
determined by [unicode.IsLetter](https://golang.org/pkg/unicode/#IsLetter).

There are some input files with pathnames of the form `pg-*.txt`,
downloaded from [Project 
Gutenberg](https://www.gutenberg.org/ebooks/search/%3Fsort_order%3Ddownloads).
