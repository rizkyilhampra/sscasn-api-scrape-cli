# SSCASN SCRAPER
This project is a personal and educational tool designed to scrape job vacancy data from the SSCASN API (https://api-sscasn.bkn.go.id), which powers the SSCASN portal (https://sscasn.bkn.go.id). The program fetches data, processes it, and exports the results into an Excel file (.xlsx) stored in the "data/" directory. It leverages concurrent processing with goroutines, allowing up to 10 requests per second for efficient data retrieval.

## Table of Contents
- [SSCASN Scraper](#sscasn-scraper)
  - [Features](#features)
  - [Prerequisites](#prerequisites)
  - [How to Run](#how-to-run)
  - [How to Obtain `kodeRefPend` (Bahasa)](#how-to-obtain-koderefpend-bahasa)

## Features

- **Concurrent Data Fetching**: Utilizes goroutines to handle multiple requests simultaneously, improving performance.
- **Rate Limiting**: Ensures compliance with API request limits to prevent server overload.
- **Excel Export**: Converts and saves the fetched data into a structured Excel file for easy analysis.
- **Customizable Filters**: Allows filtering by location (province) to tailor the data to specific needs.

## Prerequisites

Before running the SSCASN Scraper, ensure you have the following installed on your system:

1. **Go Programming Language**: 
   - Make sure you have Go installed. You can download it from the [official Go website](https://golang.org/dl/).
   - Verify the installation by running `go version` in your terminal.

2. **Git**:
   - Git is required to clone the repository. Download it from the [official Git website](https://git-scm.com/).
   - Verify the installation by running `git --version` in your terminal.

Optional for running multiple `kodeRefPend`:

1. **jq**:
   - `jq` is a lightweight and flexible command-line JSON processor. It is used in the `run_all_programs.sh` script.
   - Install `jq` by following the instructions on the [official jq website](https://stedolan.github.io/jq/download/).
   - Verify the installation by running `jq --version` in your terminal.

2. **Permissions**:
   - Ensure that the `run_all_programs.sh` script is executable. You can make it executable by running `chmod +x run_all_programs.sh` in your terminal.

## How to Run

1. Clone the repository and navigate to the project directory:
    ```bash
    git clone https://github.com/rizkyilhampra/sscasn-scraper.git
    cd sscasn-scraper
    go mod tidy
    ```
2. Run the program with specific parameters. For example, to fetch data for "S1 Pendidikan Keagamaan Katolik":
    ```bash
    go run mod.go -kodeRefPend=5102656 -namaJurusan="S1 Pendidikan Keagamaan Katolik" 
    ```
3. Optionally, filter results by province using the -provinsi flag:
    ```bash
    go run mod.go -kodeRefPend=5102656 -namaJurusan="S1 Pendidikan Keagamaan Katolik" -provinsi="Jawa Tengah"
    ```
4. To process multiple `kodeRefPend` values, list them in `data.json` and execute:
    ```bash
    ./run_all_programs.sh
    ```

## How to Obtain kodeRefPend (Bahasa)
1. Open your browser and navigate to https://sscasn.bkn.go.id/.
2. Access the network tab in the browser's developer tools (Inspect Element).
3. Search for the desired major in the website's search bar.
4. Identify the network request containing `kode_ref_pend` in its path from the host `api-sscasn.bkn.go.id`.
5. Copy the `kode_ref_pend` value for use in the program.
