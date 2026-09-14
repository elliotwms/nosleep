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

Prebuilt binaries for Linux, macOS and Windows are attached to each
[release](https://github.com/elliotwms/nosleep/releases), with a
`checksums.txt` to verify them against.

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

**Stale directives are reported too.** A `//nosleep:allow` which governs no
`time.Sleep` is flagged, so a justification does not outlive the sleep it was
written for and sit there looking load-bearing:

```go
func TestReady(t *testing.T) {
	//nosleep:allow the fixture server has no readiness endpoint
	waitForReady() // the sleep went away; the excuse did not
}
```

Only files which are actually being checked are considered, so a directive in a
non-test file is left alone unless you pass `-all-files`.

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

`nosleep` is a golangci-lint [module plugin][module-plugins]. golangci-lint
v2 is required, and the Go toolchain must be installed, since a module plugin
is compiled into a custom golangci-lint binary rather than loaded at runtime.

### 1. Build a golangci-lint which includes nosleep

Add a `.custom-gcl.yml` next to your `.golangci.yml`:

```yaml
# The golangci-lint release to build. Must be a v2 tag.
version: v2.13.2
plugins:
  - module: github.com/elliotwms/nosleep
    # A nosleep release tag; see the releases page for the latest.
    version: v1.2.3
```

then build it:

```sh
golangci-lint custom
```

This produces a `custom-gcl` binary in the current directory which is
golangci-lint with `nosleep` compiled in. Add `custom-gcl` to your `.gitignore`.

`version` is handed straight to `go get`, so anything `go get` accepts works
there. Before the first tagged release, or to try an unreleased change, pin a
commit instead:

```yaml
  - module: github.com/elliotwms/nosleep
    version: 5f6dae9
```

To hack on `nosleep` itself, point at a local checkout with `path` in place of
`version`:

```yaml
  - module: github.com/elliotwms/nosleep
    path: ../nosleep
```

### 2. Enable it

Custom linters are not enabled by default, so both the `enable` entry and the
`custom` settings block are needed:

```yaml
version: "2"
linters:
  enable:
    - nosleep
  settings:
    custom:
      nosleep:
        type: module
        description: Reports time.Sleep calls that have not been justified.
        original-url: github.com/elliotwms/nosleep
        settings:
          all-files: false # the default; true checks every file
```

The keys under `settings` are the same as the flags of the standalone binary.

### 3. Run it

Use `custom-gcl` wherever you would use `golangci-lint`:

```sh
./custom-gcl run ./...
```

### 4. Run it in CI

With GitHub Actions, install golangci-lint with the official action, build the
custom binary, and run that instead:

```yaml
- uses: actions/setup-go@v5
  with:
    go-version: stable

- uses: golangci/golangci-lint-action@v9
  with:
    version: v2.13.2
    install-only: true

- run: golangci-lint custom

- run: ./custom-gcl run ./...
```

Keep the `version` here in step with the one in `.custom-gcl.yml`.

### Allowing a sleep under golangci-lint

`//nosleep:allow <reason>` works exactly as it does standalone. golangci-lint's
own `//nolint:nosleep` is also honoured, but it does not insist on a reason, so
prefer `//nosleep:allow` if the justification is the point.

[module-plugins]: https://golangci-lint.run/docs/plugins/module-plugins/

## Why did you do this

1. Why not!
2. To learn how to write a linter
3. To learn how [`go/analysis`](https://pkg.go.dev/golang.org/x/tools/go/analysis) works
4. To detect usages of `time.Sleep` (primarily in tests), and make sure they're _really_ necessary

## License

[MIT](LICENSE)
