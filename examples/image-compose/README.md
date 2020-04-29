# HOW TO RUN THIS EXAMPLE

### In case there's no flow binary in this folder

```bash
go build -o flow ../../*.go
```

### Run this example

```bash
cat config.json | ./flow flow.def
```

## Configuration

`config.json` can be adjusted. The following constants can be used for blend mode and output format:

### Blend modes:

- SOURCE_OVER
- DESTINATION_OVER
- MULTIPLY
- SCREEN
- OVERLAY
- DARKEN
- LIGHTEN
- HARDLIGHT
- DIFFERENCE
- EXCLUSION

### Output formats:

- PNG
- JPEG
- BMP
