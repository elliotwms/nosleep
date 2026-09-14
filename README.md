# nosleep

The world's spookiest linter

`nosleep` is a Go linter which reports usages of `time.Sleep`.

Sleeping to wait for something to happen makes a test slow when the sleep is too
long and flaky when it is too short. `nosleep` reports every `time.Sleep` so that
each one has to be justified, rather than reached for out of habit.

## Install

```sh
go install github.com/elliotwms/nosleep/cmd/nosleep@latest
```

## Usage

```sh
nosleep ./...
```

By default only test files are checked, since sleeping in production code is
often legitimate — backoff, rate limiting and polling hardware are all real
reasons to sleep. Pass `-all-files` to check every file:

```sh
nosleep -all-files ./...
```

## Allowing a sleep

Sometimes a sleep really is necessary. Say so with a `//nosleep:allow` directive
and a reason, either on the same line as the call or on the line above it:

```go
func TestReconnect(t *testing.T) {
	time.Sleep(time.Second) //nosleep:allow the fixture server has no readiness endpoint
}
```

```go
func TestDebounce(t *testing.T) {
	//nosleep:allow waiting out the debounce interval is the behaviour under test
	time.Sleep(500 * time.Millisecond)
}
```

**The reason is mandatory.** A bare `//nosleep:allow` is itself reported and does
not suppress anything, so the directive cannot be used to wave a sleep through
without explanation. That requirement is the point of the linter: the goal is not
to ban sleeping, but to make sure every sleep is _really_ necessary.

## What is detected

`nosleep` resolves calls through the type checker rather than matching on the
text of the source, so aliased and dot imports are caught and unrelated methods
named `Sleep` are not:

```go
import clock "time"

clock.Sleep(clock.Second) // reported
```

```go
type worker struct{}

func (worker) Sleep(d time.Duration) {}

func f() {
	time := worker{}
	time.Sleep(0) // not reported: nothing to do with the time package
}
```

One gap is known and deliberate: a call made through a function value, such as
`var sleep = time.Sleep` followed by `sleep(d)`, is not reported. Tracking that
would mean following assignments for little benefit.

## golangci-lint

Not supported yet. A [module plugin][module-plugins] is the intended route; see
the issue tracker.

[module-plugins]: https://golangci-lint.run/docs/plugins/module-plugins/

## Why did you do this

1. Why not!
2. To learn how to write a linter
3. To learn how [`go/analysis`](https://pkg.go.dev/golang.org/x/tools/go/analysis) works
4. To detect usages of `time.Sleep` (primarily in tests), and make sure they're _really_ necessary
