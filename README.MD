# Go Code Test Analyzer (GCTA)

Go Code Test Analyzer (GCTA) is a tool built with Go to analyze Go codebases. It traverses through your Go files, counting the number of functions and test functions, providing you with a clear picture of your codebase's structure - how many functions you have and how many tests you have for each file.

## Features

1. Analyzes each Go file in the codebase for function count and test function count.
2. Provides output in a neat tabular format in the terminal.
3. Supports output in markdown format for easy documentation.
4. Uses concurrency for efficient analysis of large codebases.

## Installation

To install GCTA, you need to have Go installed on your system. Once you have Go installed:

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/go-code-analyzer.git
   ```
2. Navigate to the cloned directory:
   ```
   cd go-code-analyzer
   ```
3. Build the project:
   ```
   go build -o GCTA
   ```
4. Move the `GCTA` binary to a directory in your `$PATH`:
   ```
   sudo mv GCTA /usr/local/bin/
   ```

## Usage

You can use GCTA by navigating to any directory containing Go files and running the command `GCTA`.

Optional flag:

- `-markdown`: Use this flag to output the results in markdown format.

Example:
```
gcta -markdown
```

## Contributing

Contributions to GCTA are welcome! Please see our [contributing guide](CONTRIBUTING.md) for more details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.