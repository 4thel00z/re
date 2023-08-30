# `re` - Quick File Content Replacer

A command-line utility to search and optionally replace patterns in files across directories. This tool provides a simple interface to make sweeping changes to files without the need for complicated commands or scripts.

---

## Features

- **Easy Search & Replace**: Find and replace patterns across multiple files and directories.
- **Include/Exclude Patterns**: Ability to specify which files to include or exclude using regular expressions.
- **Dry Run**: Preview changes before applying them.
- **Optimized**: Skips larger files and uses efficient search algorithms for faster performance.

---

## Installation

### Using `go install`

Run the following command to install `re`:

```sh
go install github.com/4thel00z/re@latest
```

Make sure your `$GOPATH/bin` is in your `$PATH` to use the `re` command anywhere.

---

## Usage

1. **Basic Search**:

   Replace the term "needle" with "replacement" in the current directory (dry run):

   ```sh
   re needle replacement
   ```

2. **Applying Changes**:

   To apply the replacements:

   ```sh
   re -f needle replacement
   ```

3. **Specifying Directories**:

   Search and replace in specific directories:

   ```sh
   re needle replacement dir1 dir2
   ```

4. **Include/Exclude Patterns**:

   Use `-i` to specify file patterns to include and `-e` to exclude:

   ```sh
   re -i "*.js,*.html" -e ".git" needle replacement
   ```

For additional options and usage details:

```sh
re --help
```

---

## Contributing

Contributions are always welcome! See the `CONTRIBUTING.md` file for guidelines.

---

## License

This project is licensed under the GPL-3 license.
