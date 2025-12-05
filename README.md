# letterboxdgo

A library and collection of tools to gather information from [Letterboxd](https://letterboxd.com/).

## Requirements

This library does not require an API key to gather diary information. However, to gather additional film information from [TMDB](https://www.themoviedb.org), you must provide an [API key]. To provide a TMDB API key, set the `TMDB_KEY` environment variable:

```
$ export TMDB_KEY=<TMDB API KEY>
```

## Get stats

Gather diary information and write it to a CSV file. The folloing command gathers diary information for the user set by the `-username` flag. The `-year` flag sets the max diary history, and the `-out` flag sets the filename to write the stats for. The following example gathers information for films watched since 2024 and stores the CSV in a file named `stats.csv`

```
$ go run cmd/get/main.go -username <USERNAME> -year 2024 -out stats.csv
```