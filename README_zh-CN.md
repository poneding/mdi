# mdi

[English](README.md) ｜ 中文

mdi(Markdown indexer) 是一个命令行工具，用于在目录下递归地生成 Markdown 索引。

## 安装

```bash
go install github.com/poneding/mdi@latest
```

## 使用

生成 Markdown 索引：

```bash
mdi
```

- `-d` or `--workdir`: 指定要生成 Markdown 索引的目录
- `-f` or `--index-file`: 指定输出 Markdown 索引文件，默认为 `index.md`
- `-t` or `--index-title`: 指定 Markdown 索引标题，默认为 Markdown 索引文件的一级标题或当前目录名
- `--exclude`: 排除目录或文件，用逗号分隔
- `--override`: 覆盖现有的 Markdown 索引文件，默认为 `false`
- `-r` or `--recursive`: 递归在子目录中生成 Markdown 索引，默认为 `true`
- `--nav`: 在 Markdown 文件中生成导航，默认为 `false`
- `-v` or `--verbose`: 显示详细日志，默认为 `false`

其他命令：

```bash
# 打印版本
mdi version

# 打印帮助
mdi help

# 自动补全
mdi completion
# 示例:
# source <(mdi completion zsh)
```
