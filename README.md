# mdi

English | [中文](README_zh-CN.md)

mdi(Markdown indexer) is a command line tool used to recursively generate markdown indexes in directories.

## Installation

```bash
go install github.com/poneding/mdi@latest
```

## Usage

Generate markdown index:

```bash
mdi
```

- `-d` or `--workdir`: Specify the directory to generate markdown index.
- `-t` or `--index-title`: Specify the title of markdown index, default is title of markdown index file or current directory name.
- `-f` or `--index-file`: Specify the markdown index file, default is `index.md`.
- `--exclude`: Exclude directories or files, separated by commas.
- `--override`: Override markdown existing index file, default is `false`.
- `-r` or `--recursive`: Recursively generate markdown index in subdirectories, default is `false`.
- `--nav`: Generate navigation in markdown file, default is `false`.
- `-v` or `--verbose`: Show verbose log, default is `false`.

Other commands:

```bash
# Print version
mdi version

# Print help
mdi help

# Auto Complete
mdi completion
# Example:
# source <(mdi completion zsh)
```

## Screenshots

Markdown folder:

![20231124165206](https://images.poneding.com/2023/11/20231124165206.png)

**Generate markdown index**：

```bash
mdi -f README.md -t "My Notes"
```

![20231124170746](https://images.poneding.com/2023/11/20231124170746.png)

As the same time, the index file will be generated in the subdirectory recursively.

**Generate nav in markdown file**:

```bash
mdi -f README.md -t "My Notes" --nav
```

![20231124170902](https://images.poneding.com/2023/11/20231124170902.png)

**Customize sub index title**:

You can customize the sub index title by modifying the first-level title of the sub markdown index file generated in the subdirectory.
