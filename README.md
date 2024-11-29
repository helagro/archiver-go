Simple utility to move files to trash. Tested and built for macOS. Example settings.yaml file:

```yaml
rules:
  - path: Downloads
    pattern: .+
    days: 10
  - path: Developer
    pattern: .+
    days: 14
  - path: Pictures
    pattern: .+\.jpg
    days: 28
  - path: Desktop
    pattern: .+
    days: -1 # delete immediately
exclude:
  - .*/\.DS_Store
  - .*/\.localized
root: /Users/h
trash: /Users/h/.Trash
```

Simple chained build command:

```zsh
go build . && rm -f build/macOS && mv archiver-go build/macOS
```