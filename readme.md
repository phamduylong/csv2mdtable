# CSV To Markdown Table Converter

This is a utility tool used to convert [Comma-separated values (CSV)](https://en.wikipedia.org/wiki/Comma-separated_values) files to a [Markdown table](https://www.markdownguide.org/extended-syntax/#tables). The tool was written with [Go](https://go.dev/), a programming language known for being fantastic with string processing, great performance and a collection of other great features.

## Table Of Contents

- [CSV To Markdown Table Converter](#csv-to-markdown-table-converter)
  - [Table Of Contents](#table-of-contents)
  - [Usage](#usage)
  - [Configuration Options](#configuration-options)
    - [Align](#align)
    - [Delimiter](#delimiter)
    - [Input File](#input-file)
    - [URL](#url)
    - [Verbose Logging](#verbose-logging)
  - [Performance](#performance)

## Usage

The easiest and recommended way to use this tool is through the command line. Download the [.exe executable file](https://github.com/phamduylong/csv-to-md/releases) and run the tool from a terminal of your choice.

An example command to run the tool can look like this:

```go
import (
  "github.com/phamduylong/csv-to-md"
  "fmt"
  "errors"
)

func convertMyCSV() (string, error) {
  var cfg Config
  cfg.URL = "https://raw.githubusercontent.com/askmedia/datalogue/master/olympics/winter_olympics_medals.csv"
  cfg.Align = Center
  cfg.VerboseLogging = true
  table, err := Convert(cfg)
  
  fmt.Println(table)
}
```

## Configuration Options

The program offers a range of different configuration options to customize the tool to best fit your use case.

### Align

<sub>*(Optional)*</sub>

This option decides should text be aligned in the table. Available choices:
- 0 - Center (default)
- 1 - Left
- 2 - Right


### Delimiter

<sub>*(Optional)*</sub>

Set the delimiter character for the CSV parser. If not given, the comma `,` character will be used. This can come in handy when your CSV file use different delimiters than the popular comma.

### Input File

The path to the input CSV file to be converted. This can be either [absolute or relative path](https://en.wikipedia.org/wiki/Path_(computing)#).

**NOTE**: This should be used mutual exclusively with [URL](#url) (meaning only one of them should be used at a time)

### URL

The URL from which the CSV data can be fetched from. For this to work, the HTTPS response sent by the server must set the [Media Type](https://en.wikipedia.org/wiki/Media_type#Common_examples) to `text/csv; charset=utf-8` and the [return status code](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#Standard_codes) must be in range of [success (200-299)](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#2xx_success) or [redirection (300-399)](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#3xx_redirection). Last but not least, the [response body](https://en.wikipedia.org/wiki/HTTP#Body) should contain only the CSV data and nothing else.

**NOTE**: This should be used mutual exclusively with [Input File](#input-file) (meaning only one of them should be used at a time).

### Verbose Logging

<sub>*(Optional)*</sub>

Should detailed diagnostic messages be logged? By default, these messages are kept out of the process to avoid polluting the program output and prevent unnecessary information from confusing users. It could be helpful to track down the problem when the program execution went wrong.

## Performance

*I did not create a very proper setup to measure the performance. Ran it with my own PC so take it with a grain of salt.*

| Rows    | Columns | Average Execution Time (5 runs) |
| ------- | ------- | ------------------------------- |
| 1.000   | 10      | 29ms                            |
| 10.000  | 10      | 2,1s                            |
| 50.000  | 10      | 48,2s                           |
| 100.000 | 10      | 3 minutes                       |
