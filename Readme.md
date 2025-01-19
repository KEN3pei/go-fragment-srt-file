### 参照

https://github.com/asticode/go-astisub

上記パッケージの一部関数を拡張して、字幕ファイル（.srt）を指定した時間で分割して出力できるようにしたもの。

### 使い方

1. srtSubtitles.srtを用意する

2. go run main.go を実行

```shell
$ go run main.go
flangment: 0 ~ 197
flangment: 197 ~ 437
flangment: 437 ~ 703
flangment: 703 ~ 988
flangment: 988 ~ 1267
```

3. new_srtSubtitles_{num}.srt でファイルが分割されて出力される

### 仕様

デフォルトで10分単位で分割する。（変更したければmain.goを修正する）
