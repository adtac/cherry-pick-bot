# cherry-pick-bot

A bot to help you cherry-pick PRs that are ready to merge but need a rebase.

To get this running in your server, please fill the `src/config/config.go`
file with your credentials.

### Installation

Dependencies: [`oauth2`](https://godoc.org/golang.org/x/oauth2), [`go-github`](https://godoc.org/github.com/google/go-github/github), and [`context`](https://godoc.org/golang.org/x/net/context). To install these:

```bash
$ go get -u golang.org/x/oauth2
$ go get -u github.com/google/go-github
$ go get -u golang.org/x/net/context
```

### License

```
Copyright 2017 Adhityaa Chandrasekar

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
