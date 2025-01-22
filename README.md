# imgzip

A simple and efficient command-line tool for image compression. It supports common image formats and provides easy-to-use quality control.

## Features

- Simple command-line interface
- Supports common image formats (JPG, JPEG, PNG)
- Automatic file size optimization (targets around 1MB)
- Quality control from 1-100
- Preserves original aspect ratio
- Automatically generates output filename with timestamp

## Installation

Install imgzip using Homebrew:

```bash
brew tap owo-network/brew
brew update
brew install imgzip --cask
```

## Usage

Basic syntax:
```bash
imgzip [filename] [quality]
```

- `filename`: Path to the image file you want to compress
- `quality`: Optional. Compression quality (1-100). Default is 80

Examples:
```bash
# Compress with default quality (80)
imgzip photo.jpg

# Compress with specific quality
imgzip photo.jpg 60

# Compress PNG file
imgzip screenshot.png 70
```

Output will be saved in the current directory with format: `originalname_compressed_HHMMSS.ext`

Example output:
```
Compression complete!
Original size: 2048 KB
Compressed size: 512 KB
Compression ratio: 25.00%
Output file: photo_compressed_153042.jpg
```

## Build from Source

Requirements:
- Go 1.20 or higher
- Dependencies: github.com/disintegration/imaging

```bash
git clone https://github.com/missuo/imgzip.git
cd imgzip
go build
```

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Author

**Vincent Yang**
- Telegram: [@missuo](https://t.me/missuo)
- GitHub: [@missuo](https://github.com/missuo)

## Contributing

Contributions, issues, and feature requests are welcome! Feel free to check [issues page](https://github.com/missuo/imgzip/issues).

## Show your support

Give a ⭐️ if this project helped you!

---
Copyright © 2025 by Vincent. All Rights Reserved.