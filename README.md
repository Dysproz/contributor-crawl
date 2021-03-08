# Contributiors Crawl

This is a simple script that allows to easily from CLI gather CSV file with all contributions across organization in specific date range.

## Usage

Clone this rpeository and then using `go build` command build binary of the code or download already compiled binary from the release files.
After that make sure that binary is executable `chmod +x ./contributors-crawl`.
Finally run `./contributors-crawl --help` to receive help instructions of the tool.

## Flags

These are possible flags for the contributors-crawl CLI tool:

- `-organization`: This is name of the organization that you want to crawl for contribution data
- `-repository`: This is the name of specific repository that you are interested in. This flag is optional and if left empty then script will automatically deect all repositories in organization and crawl through them.
- `-start-date`: This is a start date from which you want to count contributions. It has to be in format YYYY-MM-DD
- `-end-date`: This is an end date to which you want to count contributions. It has to be in format YYYY-MM-DD
- `-out-file`: This is path to the file where data should be written. If left empty then it defaults to *contribution-crawl-out.csv*.
- `-oauth`: This is API token for GitHub account. If left empty then it defaults to use unauthorized client. Unauthorized GitHub clients can make only 60 requests per hour which is usually not enought for big organizations. Authorized users can make up to 15k requests per hour which should be enought for most use cases. API token may be created [here](https://github.com/settings/tokens). Please notice that you don't need to specify any permissions for the API token as it will usually use publicly available read requests to GitHub.

## Example

`./contributors-crawl -organization tungstenfabric -start-date 2020-12-01 -end-date 2022-07-20 -out-file 3months-contribs.csv -oauth AAAAAAAAAAAAAAAAAAAAAAAAAAAAA`