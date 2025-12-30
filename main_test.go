package csv2mdtable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/* STRING FUNCTION */
const STRINGS_SHOULD_BE_THE_SAME = "The two strings should be the same"

func TestPadStart(t *testing.T) {
	originalString := "start"
	expected := "     start"
	res, err := padStart(originalString, 10, ' ')

	assert.Nil(t, err, "padStart should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestPadEnd(t *testing.T) {
	originalString := "end"
	expected := "end       "
	res, err := padEnd(originalString, 10, ' ')

	assert.Nil(t, err, "padEnd should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestPadCenterEven(t *testing.T) {
	originalString := "eleven"
	expected := "  eleven  "
	res, err := padCenter(originalString, 10, ' ')

	assert.Nil(t, err, "padCenter should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestPadCenterOdd(t *testing.T) {
	originalString := "eight"
	expected := "  eight   "
	res, err := padCenter(originalString, 10, ' ')

	assert.Nil(t, err, "padCenter should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

/* Conversion */
const csvString = `First name,Last name,Email,Phone
Jane,Smith,jane.smith@email.com,555-555-1212
John,Doe,john.doe@email.com,555-555-3434
Alice,Wonder,alice@wonderland.com,555-555-5656`

const csvStringWithNarrowColumn = `#,first name,last name,email,gender
1,Herman,Gribbin,hgribbin0@deliciousdays.com,Male
2,Bing,Langthorne,blangthorne1@a8.net,Male
3,Keith,Hansford,khansford2@reference.com,Male`

const csvStringWithSemiColon = `First name;Last name;Email;Phone
Jane;Smith;jane.smith@email.com;555-555-1212
John;Doe;john.doe@email.com;555-555-3434
Alice;Wonder;alice@wonderland.com;555-555-5656`

const csvStringWithCommentLines = `First name,Last name,Email,Phone
Jane,Smith,jane.smith@email.com,555-555-1212
- this is a commented line    Jane,Smith,jane.smith@email.com,555-555-1212
John,Doe,john.doe@email.com,555-555-3434
- this is a commented line John,Doe,john.doe@email.com,555-555-3434
Alice,Wonder,alice@wonderland.com,555-555-5656
- this is a commented line  Alice,Wonder,alice@wonderland.com,555-555-5656`

const csvStringWithWhiteSpaces = `First name,   Last name,  Email,Phone
   Jane,   Smith,jane.smith@email.com,555-555-1212
John,Doe,  john.doe@email.com,555-555-3434
Alice,Wonder,    alice@wonderland.com,    555-555-5656`

const csvStringWithPipeCharacters = `ID,Expression,Description
1,A || B,Logical OR using pipe
2,foo | bar | baz,Chained pipe values
3,cmd1 | cmd2,Unix-style pipe between commands
4,x | y == z,Comparison involving a pipe operator`

func TestConvertGeneric(t *testing.T) {
	cfg := createGenericConfig()

	expected := `| First name | Last name |        Email         |    Phone     |
| :--------: | :-------: | :------------------: | :----------: |
|    Jane    |   Smith   | jane.smith@email.com | 555-555-1212 |
|    John    |    Doe    |  john.doe@email.com  | 555-555-3434 |
|   Alice    |  Wonder   | alice@wonderland.com | 555-555-5656 |`

	res, err := Convert(csvString, cfg)

	assert.Nil(t, err, "Convert should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)

}

func TestConvertWithNarrowColumnCenterAlign(t *testing.T) {
	cfg := createGenericConfig()
	cfg.Align = Center

	expected := `|  #  | first name | last name  |            email            | gender |
| :-: | :--------: | :--------: | :-------------------------: | :----: |
|  1  |   Herman   |  Gribbin   | hgribbin0@deliciousdays.com |  Male  |
|  2  |    Bing    | Langthorne |     blangthorne1@a8.net     |  Male  |
|  3  |   Keith    |  Hansford  |  khansford2@reference.com   |  Male  |`

	res, err := Convert(csvStringWithNarrowColumn, cfg)

	assert.Nil(t, err, "Convert with narrow column align center should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestConvertWithNarrowColumnLeftAlign(t *testing.T) {
	cfg := createGenericConfig()
	cfg.Align = Left

	expected := `| #  | first name | last name  | email                       | gender |
| :- | :--------- | :--------- | :-------------------------- | :----- |
| 1  | Herman     | Gribbin    | hgribbin0@deliciousdays.com | Male   |
| 2  | Bing       | Langthorne | blangthorne1@a8.net         | Male   |
| 3  | Keith      | Hansford   | khansford2@reference.com    | Male   |`

	res, err := Convert(csvStringWithNarrowColumn, cfg)

	assert.Nil(t, err, "Convert with narrow column align left should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestConvertWithNarrowColumnRightAlign(t *testing.T) {
	cfg := createGenericConfig()
	cfg.Align = Right

	expected := `|  # | first name |  last name |                       email | gender |
| -: | ---------: | ---------: | --------------------------: | -----: |
|  1 |     Herman |    Gribbin | hgribbin0@deliciousdays.com |   Male |
|  2 |       Bing | Langthorne |         blangthorne1@a8.net |   Male |
|  3 |      Keith |   Hansford |    khansford2@reference.com |   Male |`

	res, err := Convert(csvStringWithNarrowColumn, cfg)

	assert.Nil(t, err, "Convert with narrow column align right should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestLeftAlign(t *testing.T) {
	cfg := createGenericConfig()
	cfg.Align = Left

	expected := `| First name | Last name | Email                | Phone        |
| :--------- | :-------- | :------------------- | :----------- |
| Jane       | Smith     | jane.smith@email.com | 555-555-1212 |
| John       | Doe       | john.doe@email.com   | 555-555-3434 |
| Alice      | Wonder    | alice@wonderland.com | 555-555-5656 |`

	res, err := Convert(csvString, cfg)

	assert.Nil(t, err, "Convert with left align should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestRightAlign(t *testing.T) {
	cfg := createGenericConfig()
	cfg.Align = Right

	expected := `| First name | Last name |                Email |        Phone |
| ---------: | --------: | -------------------: | -----------: |
|       Jane |     Smith | jane.smith@email.com | 555-555-1212 |
|       John |       Doe |   john.doe@email.com | 555-555-3434 |
|      Alice |    Wonder | alice@wonderland.com | 555-555-5656 |`

	res, err := Convert(csvString, cfg)

	assert.Nil(t, err, "Convert with right align should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestWithCustomDelimiter(t *testing.T) {
	cfg := createGenericConfig()
	cfg.CSVReaderConfig.Comma = ';'
	expected := `| First name | Last name |        Email         |    Phone     |
| :--------: | :-------: | :------------------: | :----------: |
|    Jane    |   Smith   | jane.smith@email.com | 555-555-1212 |
|    John    |    Doe    |  john.doe@email.com  | 555-555-3434 |
|   Alice    |  Wonder   | alice@wonderland.com | 555-555-5656 |`

	res, err := Convert(csvStringWithSemiColon, cfg)

	assert.Nil(t, err, "Convert with custom delimiter should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)

}

func TestWithCustomComment(t *testing.T) {
	cfg := createGenericConfig()
	cfg.CSVReaderConfig.Comment = '-'
	expected := `| First name | Last name |        Email         |    Phone     |
| :--------: | :-------: | :------------------: | :----------: |
|    Jane    |   Smith   | jane.smith@email.com | 555-555-1212 |
|    John    |    Doe    |  john.doe@email.com  | 555-555-3434 |
|   Alice    |  Wonder   | alice@wonderland.com | 555-555-5656 |`

	res, err := Convert(csvStringWithCommentLines, cfg)

	assert.Nil(t, err, "Convert with comment lines should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestWithCorrectCustomFieldsPerRecord(t *testing.T) {
	cfg := createGenericConfig()
	cfg.CSVReaderConfig.FieldsPerRecord = 4
	expected := `| First name | Last name |        Email         |    Phone     |
| :--------: | :-------: | :------------------: | :----------: |
|    Jane    |   Smith   | jane.smith@email.com | 555-555-1212 |
|    John    |    Doe    |  john.doe@email.com  | 555-555-3434 |
|   Alice    |  Wonder   | alice@wonderland.com | 555-555-5656 |`

	res, err := Convert(csvString, cfg)

	assert.Nil(t, err, "Convert with correct fields per record setting should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestWithMismatchCustomFieldsPerRecord(t *testing.T) {
	cfg := createGenericConfig()
	cfg.CSVReaderConfig.FieldsPerRecord = 17
	expected := `| First name | Last name |        Email         |    Phone     |
| :--------: | :-------: | :------------------: | :----------: |
|    Jane    |   Smith   | jane.smith@email.com | 555-555-1212 |
|    John    |    Doe    |  john.doe@email.com  | 555-555-3434 |
|   Alice    |  Wonder   | alice@wonderland.com | 555-555-5656 |`

	res, err := Convert(csvString, cfg)

	assert.NotNil(t, err, "Convert with mismatch fields per record setting should return an error")

	assert.NotEqual(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestWithTrimLeadingWhiteSpace(t *testing.T) {
	cfg := createGenericConfig()
	cfg.CSVReaderConfig.TrimLeadingSpace = true
	expected := `| First name | Last name |        Email         |    Phone     |
| :--------: | :-------: | :------------------: | :----------: |
|    Jane    |   Smith   | jane.smith@email.com | 555-555-1212 |
|    John    |    Doe    |  john.doe@email.com  | 555-555-3434 |
|   Alice    |  Wonder   | alice@wonderland.com | 555-555-5656 |`

	res, err := Convert(csvStringWithWhiteSpaces, cfg)

	assert.Nil(t, err, "Convert with trim leading whitespace setting should not return an error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestWithCaption(t *testing.T) {
	cfg := createGenericConfig()
	cfg.Caption = "Table 2: Customers who are United fans"
	expected := `<!-- Table 2: Customers who are United fans -->
| First name | Last name |        Email         |    Phone     |
| :--------: | :-------: | :------------------: | :----------: |
|    Jane    |   Smith   | jane.smith@email.com | 555-555-1212 |
|    John    |    Doe    |  john.doe@email.com  | 555-555-3434 |
|   Alice    |  Wonder   | alice@wonderland.com | 555-555-5656 |`

	res, err := Convert(csvString, cfg)

	assert.Nil(t, err, "Convert with caption setting should not return an error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestCompactConvertGeneric(t *testing.T) {
	cfg := createGenericConfig()
	cfg.Compact = true

	expected := `|First name|Last name|Email|Phone|
|:-:|:-:|:-:|:-:|
|Jane|Smith|jane.smith@email.com|555-555-1212|
|John|Doe|john.doe@email.com|555-555-3434|
|Alice|Wonder|alice@wonderland.com|555-555-5656|`

	res, err := Convert(csvString, cfg)

	assert.Nil(t, err, "Convert compact should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestCompactConvertGenericLeftAlign(t *testing.T) {
	cfg := createGenericConfig()
	cfg.Compact = true
	cfg.Align = Left

	expected := `|First name|Last name|Email|Phone|
|:-|:-|:-|:-|
|Jane|Smith|jane.smith@email.com|555-555-1212|
|John|Doe|john.doe@email.com|555-555-3434|
|Alice|Wonder|alice@wonderland.com|555-555-5656|`

	res, err := Convert(csvString, cfg)

	assert.Nil(t, err, "Convert compact left align should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestCompactConvertGenericRightAlign(t *testing.T) {
	cfg := createGenericConfig()
	cfg.Compact = true
	cfg.Align = Right

	expected := `|First name|Last name|Email|Phone|
|-:|-:|-:|-:|
|Jane|Smith|jane.smith@email.com|555-555-1212|
|John|Doe|john.doe@email.com|555-555-3434|
|Alice|Wonder|alice@wonderland.com|555-555-5656|`

	res, err := Convert(csvString, cfg)

	assert.Nil(t, err, "Convert compact right align should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestEscapePipeCharacter(t *testing.T) {
	cfg := createGenericConfig()
	cfg.Align = Left

	expected := `| ID | Expression        | Description                          |
| :- | :---------------- | :----------------------------------- |
| 1  | A \|\| B          | Logical OR using pipe                |
| 2  | foo \| bar \| baz | Chained pipe values                  |
| 3  | cmd1 \| cmd2      | Unix-style pipe between commands     |
| 4  | x \| y == z       | Comparison involving a pipe operator |`

	res, err := Convert(csvStringWithPipeCharacters, cfg)

	assert.Nil(t, err, "Convert with pipe characters should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestConvertExcludeAllColumnsButOne(t *testing.T) {
	cfg := createGenericConfig()
	cfg.ExcludedColumns = []string{"Email", "First name", "Phone"}

	expected := `| Last name |
| :-------: |
|   Smith   |
|    Doe    |
|  Wonder   |`

	res, err := Convert(csvString, cfg)

	assert.Nil(t, err, "Convert while excluded all columns but one should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestConvertExcludeSomeColumns(t *testing.T) {
	cfg := createGenericConfig()
	cfg.ExcludedColumns = []string{"Email", "First name"}

	expected := `| Last name |    Phone     |
| :-------: | :----------: |
|   Smith   | 555-555-1212 |
|    Doe    | 555-555-3434 |
|  Wonder   | 555-555-5656 |`

	res, err := Convert(csvString, cfg)

	assert.Nil(t, err, "Convert with excluded columns should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func TestConvertExcludeAllColumns(t *testing.T) {
	cfg := createGenericConfig()
	cfg.ExcludedColumns = []string{"Email", "Last name", "First name", "Phone"}

	res, err := Convert(csvString, cfg)

	assert.Nil(t, err, "Convert with all excluded columns should not return a non-nil error")

	assert.Empty(t, res, "String should be empty")
}

func TestConvertExcludeNoColumn(t *testing.T) {
	cfg := createGenericConfig()
	cfg.ExcludedColumns = []string{}

	expected := `| First name | Last name |        Email         |    Phone     |
| :--------: | :-------: | :------------------: | :----------: |
|    Jane    |   Smith   | jane.smith@email.com | 555-555-1212 |
|    John    |    Doe    |  john.doe@email.com  | 555-555-3434 |
|   Alice    |  Wonder   | alice@wonderland.com | 555-555-5656 |`

	res, err := Convert(csvString, cfg)

	assert.Nil(t, err, "Convert with empty list of excluded columns should not return a non-nil error")

	assert.Equal(t, expected, res, STRINGS_SHOULD_BE_THE_SAME)
}

func createGenericConfig() Config {
	var cfg Config
	return cfg
}
