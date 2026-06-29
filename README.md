go-pinyin
=========

[![Build Status](https://github.com/xuanwolei/go-pinyin/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/xuanwolei/go-pinyin/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/xuanwolei/go-pinyin)](https://goreportcard.com/report/github.com/xuanwolei/go-pinyin)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/xuanwolei/go-pinyin)](https://pkg.go.dev/github.com/xuanwolei/go-pinyin)

汉语拼音转换工具 Go 版。

本项目 fork 自 [mozillazg/go-pinyin](https://github.com/mozillazg/go-pinyin)，在保持默认行为兼容的基础上，新增了可配置字典能力。


What This Fork Solves
---------------------

原项目使用 `pinyin-data` 的全量合并字典，适合通用拼音转换，但在开启多音字模式时会把古籍、通假、低频或已淘汰读音一起返回。

这会让姓名、员工搜索、拼音索引等业务场景产生过宽候选。例如：

```
默认字典：呼 -> hū,xiāo,xū,hè,xià
现代字典：呼 -> hū
```

本 fork 主要解决：

* 支持通过 `Args.Dict` 选择拼音字典。
* 保留默认 `PinyinDict`，兼容原项目行为。
* 新增 `ModernPinyinDict`，使用现代汉语来源收窄多音字候选。
* 新增 `StandardPinyinDict`，使用《通用规范汉字字典》来源进一步收窄极低频读音，适合姓名搜索等低噪声召回场景。
* 支持调用方传入自定义 `pinyin.Dict`，便于按业务继续扩展字典。
* 在现代字典中保留常见人名多音字，例如 `曾`、`乐`、`单`、`仇`、`区`。

相比原项目：

* 原项目只有一个全局默认字典；本 fork 可以按调用场景选择字典。
* 原项目开启多音字后会返回全量候选；本 fork 可以选择现代汉语候选，减少低价值读音。
* 原项目更偏通用转换；本 fork 更适合姓名搜索、员工搜索、拼音关键词索引等需要控制召回噪声的场景。


Installation
------------

```
go get github.com/xuanwolei/go-pinyin
```

install CLI tool:

```
# go version>=1.17
go install github.com/xuanwolei/go-pinyin/cli/pinyin@latest

# go version<1.17
go get -u github.com/xuanwolei/go-pinyin/cli/pinyin

$ pinyin 中国人
zhōng guó rén
```

注意：CLI 子模块本地开发时通过 `replace` 指向当前仓库根模块；如果要发布 `cli/pinyin`，请先为根模块发布对应版本 tag，或同步调整 `cli/pinyin/go.mod` 中依赖的根模块版本。


Documentation
--------------

API documentation can be found here:
https://pkg.go.dev/github.com/xuanwolei/go-pinyin


Usage
------

```go
package main

import (
	"fmt"
	"github.com/xuanwolei/go-pinyin"
)

func main() {
	hans := "中国人"

	// 默认
	a := pinyin.NewArgs()
	fmt.Println(pinyin.Pinyin(hans, a))
	// [[zhong] [guo] [ren]]

	// 包含声调
	a.Style = pinyin.Tone
	fmt.Println(pinyin.Pinyin(hans, a))
	// [[zhōng] [guó] [rén]]

	// 声调用数字表示
	a.Style = pinyin.Tone2
	fmt.Println(pinyin.Pinyin(hans, a))
	// [[zho1ng] [guo2] [re2n]]

	// 开启多音字模式
	a = pinyin.NewArgs()
	a.Heteronym = true
	fmt.Println(pinyin.Pinyin(hans, a))
	// [[zhong] [guo] [ren]]
	a.Style = pinyin.Tone2
	fmt.Println(pinyin.Pinyin(hans, a))
	// [[zho1ng zho4ng] [guo2] [re2n]]

	// 使用规范汉字字典，过滤古籍、通假等低频读音
	a = pinyin.NewArgs()
	a.Heteronym = true
	a.Dict = pinyin.StandardPinyinDict
	fmt.Println(pinyin.Pinyin("丁呼曾乐", a))
	// [[ding] [hu] [ceng zeng] [le yue]]

	fmt.Println(pinyin.LazyPinyin(hans, pinyin.NewArgs()))
	// [zhong guo ren]

	fmt.Println(pinyin.Convert(hans, nil))
	// [[zhong] [guo] [ren]]

	fmt.Println(pinyin.LazyConvert(hans, nil))
	// [zhong guo ren]
}
```

注意：

* 默认情况下会忽略没有拼音的字符（可以通过自定义 `Fallback` 参数的值来自定义如何处理没有拼音的字符，
  详见 [示例](https://pkg.go.dev/github.com/xuanwolei/go-pinyin#example-Pinyin--FallbackCustom1)）。
* 根据 [《汉语拼音方案》](http://www.moe.gov.cn/s78/A19/yxs_left/moe_810/s230/195802/t19580201_186000.html) y，w，ü (yu) 都不是声母，
  以及不是所有拼音都有声母，如果这不是你预期的话，你可能需要的是首字母风格 `FirstLetter`
  （ [详细信息](https://github.com/mozillazg/python-pinyin#%E4%B8%BA%E4%BB%80%E4%B9%88%E6%B2%A1%E6%9C%89-y-w-yu-%E5%87%A0%E4%B8%AA%E5%A3%B0%E6%AF%8D) ）。
* `Args.Dict` 为空时使用默认字典 `PinyinDict`；如需按业务场景收窄多音字候选，可以设置为
  `StandardPinyinDict`、`ModernPinyinDict` 或调用方自定义的 `pinyin.Dict`。


Related Projects
-----------------

* [hotoo/pinyin](https://github.com/hotoo/pinyin): 汉语拼音转换工具 Node.js/JavaScript 版。
* [mozillazg/python-pinyin](https://github.com/mozillazg/python-pinyin): 汉语拼音转换工具 Python 版。
* [mozillazg/rust-pinyin](https://github.com/mozillazg/rust-pinyin): 汉语拼音转换工具 Rust 版。


pinyin data
-----------------

* `PinyinDict` 使用 [pinyin-data](https://github.com/mozillazg/pinyin-data) 的 `pinyin.txt` 数据，保持默认行为兼容。
* `ModernPinyinDict` 以 `pinyin.txt` 为基础，并使用 `kXHC1983.txt` 和 `kTGHZ2013.txt` 的现代汉语读音覆盖对应字的候选读音。
  例如 `呼` 在默认字典中包含 `hū,xiāo,xū,hè,xià`，在现代汉语字典中只保留 `hū`。
* `StandardPinyinDict` 以 `pinyin.txt` 为基础，并使用 `kTGHZ2013.txt` 的规范汉字读音覆盖对应字的候选读音。
  例如 `丁` 在默认字典和现代汉语字典中包含 `dīng,zhēng`，在规范汉字字典中只保留 `dīng`。
* 重新生成内置字典：

```
git submodule update --init _tools/pinyin-data
make gen_pinyin_dict
```


License
---------

Under the MIT License.
