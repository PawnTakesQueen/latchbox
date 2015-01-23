// +build !windows

/*
 * Copyright (C) 2015 Vi Grey. All rights reserved.
 * Copyright (C) 2013-2015 Yasuhiro Matsumoto,
 * http://mattn.kaoriya.net <mattn.jp@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the “Software”), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package runewidth

import (
  "os"
  "regexp"
  "strings"
)

var reLoc = regexp.MustCompile(`^[a-z][a-z][a-z]?(?:_[A-Z][A-Z])?\.(.+)`)

func IsEastAsian() bool {
  locale := os.Getenv("LC_CTYPE")
  if locale == "" {
    locale = os.Getenv("LANG")
  }

  /* ignore C locale */
  if locale == "POSIX" || locale == "C" {
    return false
  }
  if len(locale) > 1 && locale[0] == 'C' && (locale[1] == '.' ||
      locale[1] == '-') {
    return false
  }

  charset := strings.ToLower(locale)
  r := reLoc.FindStringSubmatch(locale)
  if len(r) == 2 {
    charset = strings.ToLower(r[1])
  }

  if strings.HasSuffix(charset, "@cjk_narrow") {
    return false
  }

  for pos, b := range []byte(charset) {
    if b == '@' {
      charset = charset[:pos]
      break
    }
  }

  mbc_max := 1
  switch charset {
  case "utf-8", "utf8":
    mbc_max = 6
  case "jis":
    mbc_max = 8
  case "eucjp":
    mbc_max = 3
  case "euckr", "euccn":
    mbc_max = 2
  case "sjis", "cp932", "cp51932", "cp936", "cp949", "cp950":
    mbc_max = 2
  case "big5":
    mbc_max = 2
  case "gbk", "gb2312":
    mbc_max = 2
  }

  if mbc_max > 1 && (charset[0] != 'u' ||
    strings.HasPrefix(locale, "ja") ||
    strings.HasPrefix(locale, "ko") ||
    strings.HasPrefix(locale, "zh")) {
    return true
  }
  return false
}
