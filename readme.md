# CSV To Markdown Table Converter [![Go Reference](https://pkg.go.dev/badge/github.com/phamduylong/csv-to-md.svg)](https://pkg.go.dev/github.com/phamduylong/csv-to-md)

This is a utility tool used to convert [Comma-separated values (CSV)](https://en.wikipedia.org/wiki/Comma-separated_values) files to a [Markdown table](https://www.markdownguide.org/extended-syntax/#tables). You can test it out in the [Go playground](https://go.dev/play/p/p42h8-jG5dQ).

## Table Of Contents

- [CSV To Markdown Table Converter](#csv-to-markdown-table-converter-)
  - [Table Of Contents](#table-of-contents)
  - [Usage](#usage)
  - [Configuration Options](#configuration-options)
  - [Performance](#performance)

## Usage

An example conversion can look like this:

```go
package main

import (
  "fmt"

  csv2mdtable "github.com/phamduylong/csv-to-md"
)

func main() {
  var cfg csv2mdtable.Config
  cfg.Align = csv2mdtable.Left
  cfg.VerboseLogging = true
  
  csv := `Index,Customer Id,First Name,Last Name,Company,City,Country,Phone
1,DD37Cf93aecA6Dc,Sheryl,Baxter,Rasmussen Group,East Leonard,Chile,229.077.5154
2,1Ef7b82A4CAAD10,Preston,Lozano,Vega-Gentry,East Jimmychester,Djibouti,5153435776
3,6F94879bDAfE5a6,Roy,Berry,Murillo-Perry,Isabelborough,Antigua and Barbuda,+1-539-402-0259
4,5Cef8BFA16c5e3c,Linda,Olsen,"Dominguez, Mcmillan and Donovan",Bensonview,Dominican Republic,001-808-617-6467x12895
5,053d585Ab6b3159,Joanna,Bender,"Martin, Lang and Andrade",West Priscilla,Slovakia (Slovak Republic),001-234-203-0635x76146`

  res, convertErr := csv2mdtable.Convert(csv, cfg)

  if convertErr != nil {
    fmt.Println(convertErr)
  }

  fmt.Printf("Converted table:\n\n%s\n", res)
}
```

Output would look like this:

```console
2009/11/10 23:00:00 DEBUG Validating config 🤔
2009/11/10 23:00:00 DEBUG Config is valid ✅
Converted table:

| Index | Customer Id     | First Name | Last Name | Company                         | City              | Country                    | Phone                  |
| :---- | :-------------- | :--------- | :-------- | :------------------------------ | :---------------- | :------------------------- | :--------------------- |
| 1     | DD37Cf93aecA6Dc | Sheryl     | Baxter    | Rasmussen Group                 | East Leonard      | Chile                      | 229.077.5154           |
| 2     | 1Ef7b82A4CAAD10 | Preston    | Lozano    | Vega-Gentry                     | East Jimmychester | Djibouti                   | 5153435776             |
| 3     | 6F94879bDAfE5a6 | Roy        | Berry     | Murillo-Perry                   | Isabelborough     | Antigua and Barbuda        | +1-539-402-0259        |
| 4     | 5Cef8BFA16c5e3c | Linda      | Olsen     | Dominguez, Mcmillan and Donovan | Bensonview        | Dominican Republic         | 001-808-617-6467x12895 |
| 5     | 053d585Ab6b3159 | Joanna     | Bender    | Martin, Lang and Andrade        | West Priscilla    | Slovakia (Slovak Republic) | 001-234-203-0635x76146 |

Program exited.
```

## Configuration Options

The program offers a range of different configuration options to customize the tool to best fit your use case.

| Option                           | Type             | What does it do? |
| -------------------------------- | ---------------- | ---------------- |
| Align                            | Align            | Align the text on the rendered table. Visual feedback on the markdown syntax is also provided. |
| Caption                          | string           | Set a caption for the table (will be rendered as an HTML comment above the table). |
| Compact                          | bool             | Set whether the Markdown table be converted to compact syntax. |
| CSVReaderConfig                  | CSVReaderConfig  | Config options to be passed into CSV reader object. See [type Reader in the encoding/csv module](https://pkg.go.dev/encoding/csv#Reader). |
| CSVReaderConfig.Comma            | rune             | Set the delimiter of the CSV reader. |
| CSVReaderConfig.Comment          | rune             | Set the comment character for the CSV reader. |
| CSVReaderConfig.FieldsPerRecord  | int              | Set the amount of fields per CSV row. |
| CSVReaderConfig.LazyQuotes       | bool             | Set whether lazy quotes are allowed. If lazy quotes are allowed, a quote may appear in an unquoted field and a non-doubled quote may appear in a quoted field. |
| CSVReaderConfig.TrimLeadingSpace | bool             | Set whether leading space before the fields' values should be ignored. |
| CSVReaderConfig.ReuseRecord      | bool             | Set whether calls to Read may return a slice sharing the backing array of the previous call's returned slice for performance. By default, each call to Read returns newly allocated memory owned by the caller. |
| ExcludedColumns                  | []string         | Set the list of columns that should be ignored when constructing the table. |
| SortColumns                      | ColumnSortOption | Should the columns be sorted and how? |
| VerboseLogging                   | bool             | Log detailed diagnostic messages when running the program. |

## Performance

*I did not create a very proper setup to measure the performance. Ran it with my own PC so take it with a grain of salt.*

| Rows    | Columns | Average Execution Time (5 runs) |
| ------- | ------- | ------------------------------- |
| 100     | 12      | 1,8ms                           |
| 1.000   | 12      | 37ms                            |
| 10.000  | 12      | 2,4s                            |
| 100.000 | 12      | 249s                            |
