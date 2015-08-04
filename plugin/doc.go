// Package plugin contains status bar widget plugins.
package plugin

// We put side effect imports in this non-main, non-test package,
// because we expect people will copy-paste the main command in
// github.com/tv42/quobar/cmd/quobar, and it's less error prone if
// that package consists of just one file. This would cause golint
// complaints, but gen-imports.go works around that. It's not
// necessarily pretty, but an easy-to-customize main is more
// important.
//
//go:generate go run ../task/gen-imports.go -o imports.gen.go github.com/tv42/quobar/plugin/...
