// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path

import (
	"solod.dev/so/bytealg"
	"solod.dev/so/errors"
	"solod.dev/so/unicode/utf8"
)

// ErrBadPattern indicates a pattern was malformed.
var ErrBadPattern = errors.New("path: syntax error in pattern")

// Match reports whether name matches the shell pattern.
// The pattern syntax is:
//
//	pattern:
//		{ term }
//	term:
//		'*'         matches any sequence of non-/ characters
//		'?'         matches any single non-/ character
//		'[' [ '^' ] { character-range } ']'
//		            character class (must be non-empty)
//		c           matches character c (c != '*', '?', '\\', '[')
//		'\\' c      matches character c
//
//	character-range:
//		c           matches character c (c != '\\', '-', ']')
//		'\\' c      matches character c
//		lo '-' hi   matches character c for lo <= c <= hi
//
// Match requires pattern to match all of name, not just a substring.
// The only possible returned error is [ErrBadPattern], when pattern
// is malformed.
func Match(pattern, name string) (bool, error) {
	for len(pattern) > 0 {
		scan := scanChunk(pattern)
		star := scan.star
		chunk := scan.chunk
		pattern = scan.rest
		if star && chunk == "" {
			// Trailing * matches rest of string unless it has a /.
			matched := bytealg.IndexByteString(name, '/') < 0
			return matched, nil
		}
		// Look for match at current position.
		match := matchChunk(chunk, name)
		// if we're the last chunk, make sure we've exhausted the name
		// otherwise we'll give a false result even if we could still match
		// using the star
		if match.ok && (len(match.rest) == 0 || len(pattern) > 0) {
			name = match.rest
			continue
		}
		if match.err != nil {
			return false, match.err
		}
		if star {
			// Look for match skipping i+1 bytes.
			// Cannot skip /.
			matched := false
			for i := 0; i < len(name) && name[i] != '/'; i++ {
				match := matchChunk(chunk, name[i+1:])
				if match.ok {
					// if we're the last chunk, make sure we exhausted the name
					if len(pattern) == 0 && len(match.rest) > 0 {
						continue
					}
					name = match.rest
					matched = true
					break
				}
				if match.err != nil {
					return false, match.err
				}
			}
			if matched {
				continue
			}
		}
		// Before returning false with no error,
		// check that the remainder of the pattern is syntactically valid.
		for len(pattern) > 0 {
			scan := scanChunk(pattern)
			if match := matchChunk(scan.chunk, ""); match.err != nil {
				return false, match.err
			}
			pattern = scan.rest
		}
		return false, nil
	}
	return len(name) == 0, nil
}

type scanResult struct {
	star  bool
	chunk string
	rest  string
}

// scanChunk gets the next segment of pattern, which is a non-star string
// possibly preceded by a star.
func scanChunk(pattern string) scanResult {
	var res scanResult
	for len(pattern) > 0 && pattern[0] == '*' {
		pattern = pattern[1:]
		res.star = true
	}
	inrange := false
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '\\':
			// error check handled in matchChunk: bad pattern.
			if i+1 < len(pattern) {
				i++
			}
		case '[':
			inrange = true
		case ']':
			inrange = false
		case '*':
			if !inrange {
				res.chunk = pattern[:i]
				res.rest = pattern[i:]
				return res
			}
		}
	}
	res.chunk = pattern
	return res
}

type matchResult struct {
	rest string
	ok   bool
	err  error
}

// matchChunk checks whether chunk matches the beginning of s.
// If so, it returns the remainder of s (after the match).
// Chunk is all single-character operators: literals, char classes, and ?.
func matchChunk(chunk, s string) matchResult {
	// failed records whether the match has failed.
	// After the match fails, the loop continues on processing chunk,
	// checking that the pattern is well-formed but no longer reading s.
	failed := false
	for len(chunk) > 0 {
		failed = failed || len(s) == 0
		switch chunk[0] {
		case '[':
			// character class
			var r rune
			if !failed {
				var n int
				r, n = utf8.DecodeRuneInString(s)
				s = s[n:]
			}
			chunk = chunk[1:]
			// possibly negated
			negated := false
			if len(chunk) > 0 && chunk[0] == '^' {
				negated = true
				chunk = chunk[1:]
			}
			// parse all ranges
			match := false
			nrange := 0
			for {
				if len(chunk) > 0 && chunk[0] == ']' && nrange > 0 {
					chunk = chunk[1:]
					break
				}
				var lo, hi rune
				esc := getEsc(chunk)
				if esc.err != nil {
					return matchResult{"", false, esc.err}
				}
				lo = esc.r
				chunk = esc.nchunk
				hi = lo
				if len(chunk) > 0 && chunk[0] == '-' {
					esc := getEsc(chunk[1:])
					if esc.err != nil {
						return matchResult{"", false, esc.err}
					}
					hi = esc.r
					chunk = esc.nchunk
				}
				match = match || (lo <= r && r <= hi)
				nrange++
			}
			failed = failed || match == negated

		case '?':
			if !failed {
				failed = s[0] == '/'
				_, n := utf8.DecodeRuneInString(s)
				s = s[n:]
			}
			chunk = chunk[1:]

		case '\\':
			chunk = chunk[1:]
			if len(chunk) == 0 {
				return matchResult{"", false, ErrBadPattern}
			}

			if !failed {
				failed = chunk[0] != s[0]
				s = s[1:]
			}
			chunk = chunk[1:]

		default:
			if !failed {
				failed = chunk[0] != s[0]
				s = s[1:]
			}
			chunk = chunk[1:]
		}
	}
	if failed {
		return matchResult{"", false, nil}
	}
	return matchResult{s, true, nil}
}

type escResult struct {
	r      rune
	nchunk string
	err    error
}

// getEsc gets a possibly-escaped character from chunk, for a character class.
func getEsc(chunk string) escResult {
	var res escResult
	if len(chunk) == 0 || chunk[0] == '-' || chunk[0] == ']' {
		res.err = ErrBadPattern
		return res
	}
	if chunk[0] == '\\' {
		chunk = chunk[1:]
		if len(chunk) == 0 {
			res.err = ErrBadPattern
			return res
		}
	}
	r, n := utf8.DecodeRuneInString(chunk)
	if r == utf8.RuneError && n == 1 {
		res.err = ErrBadPattern
	}
	res.r = r
	res.nchunk = chunk[n:]
	if len(res.nchunk) == 0 {
		res.err = ErrBadPattern
	}
	return res
}
