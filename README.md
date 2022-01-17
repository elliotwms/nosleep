# nosleep

The world's spookiest linter

`nosleep` is a golang-ci compatible linter which checks for and fails if it detects usages of `time.Sleep`.

### Why did you do this

1. Why not!
2. To learn how to write a linter
3. To learn how [`go/analysis`](https://pkg.go.dev/golang.org/x/tools/go/analysis) works
4. To detect usages of `time.Sleep` (primarily in tests), and make sure they're _really_ necessary