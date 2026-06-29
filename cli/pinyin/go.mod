module github.com/xuanwolei/go-pinyin/cli/pinyin

go 1.17

require (
	github.com/mattn/go-isatty v0.0.18
	github.com/xuanwolei/go-pinyin v0.22.0
)

require golang.org/x/sys v0.6.0 // indirect

replace github.com/xuanwolei/go-pinyin => ../..
