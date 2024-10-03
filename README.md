# SSCASN API Scraper (CLI Based)
This project is a personal and educational tool designed to scrape job vacancy data from the SSCASN API (https://api-sscasn.bkn.go.id), which powers the SSCASN portal (https://sscasn.bkn.go.id). The program fetches data, processes it, and exports the results into an Excel file (.xlsx) stored in the "data/" directory. It leverages concurrent processing with goroutines, allowing up to 10 requests per second for efficient data retrieval.

## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [How to Run](#how-to-run)
- [Command-line Arguments](#command-line-arguments)
- [How to Obtain API Parameters](#how-to-obtain-api-parameters)

## Features
- **Concurrent Data Fetching**: Utilizes goroutines to handle multiple requests simultaneously, improving performance.
- **Rate Limiting**: Ensures compliance with API request limits to prevent server overload.
- **Excel Export**: Converts and saves the fetched data into a structured Excel file for easy analysis.
- **Customizable Filters**: Allows filtering by location (province) to tailor the data to specific needs.
- **Flexible API Querying**: Supports additional parameters such as `pengadaanKd` (type of procurement) and `instansiId` (institution) for more targeted data retrieval.
- **Enhanced Data Output**: Includes additional fields like `jumlah_ms` (count of passed verification) in the exported Excel file.

## Prerequisites
Before running the SSCASN Scraper, ensure you have [**Go**](https://go.dev) installed on your system.

Optional for running multiple `kodeRefPend`:

<details>
<summary>Click to expand</summary>
    
1. **jq**:
   - `jq` is a lightweight and flexible command-line JSON processor. It is used in the `run_all_programs.sh` script.
   - Install `jq` by following the instructions on the [official jq website](https://stedolan.github.io/jq/download/).
   - Verify the installation by running `jq --version` in your terminal.
2. **Permissions**:
   - Ensure that the `run_all_programs.sh` script is executable. You can make it executable by running `chmod +x run_all_programs.sh` in your terminal.
</details>

## How to Run
1. Clone the repository and navigate to the project directory:
    ```bash
    git clone https://github.com/rizkyilhampra/sscasn-scraper.git
    cd sscasn-scraper
    ```
2. Install the required dependencies:
    ```bash
    go mod tidy
    ```
3. Run the program with specific parameters. For example, to fetch data for "S1 Pendidikan Keagamaan Katolik":
    ```bash
    go run main.go -kodeRefPend=5102656 -namaJurusan="S1 Pendidikan Keagamaan Katolik" 
    ```
4. Optionally, filter results by province using the -provinsi flag:
    ```bash
    go run main.go -kodeRefPend=5102656 -namaJurusan="S1 Pendidikan Keagamaan Katolik" -provinsi="Jawa Tengah"
    ```
5. To process multiple `kodeRefPend` values, list them in `data.json` and execute:
    ```bash
    ./run_all_programs.sh
    ```

## Command-line Arguments
The program supports the following command-line arguments:

- `-kodeRefPend`: (Required) Code reference for education.
- `-namaJurusan`: (Required) Name of the major or field of study of `kodeRefPend`. It's stand for label or title uses for generate Excel filename.
- `-provinsi`: (Optional) Filter results by province. Example: `-provinsi="Jawa Timur"`.
- `-pengadaanKd`: (Optional) Procurement code. Default is 2.
- `-instansiId`: (Optional) Institution ID. Example: `-instansiId="A5EB03E23AFBF6A0E040640A040252AD"` (for Kementerian Lingkungan Hidup dan Kehutanan).

Example usage with all parameters:
```bash
go run main.go -kodeRefPend=5102656 -namaJurusan="S1 Pendidikan Keagamaan Katolik" -provinsi="Jawa Tengah" -pengadaanKd=2 -instansiId="A5EB03E23AFBF6A0E040640A040252AD"
```
## How to Obtain API Parameters

To use this scraper effectively, you'll need to obtain several parameters from the SSCASN website. Here's a general guide on how to find these parameters:

1. Open your browser and navigate to https://sscasn.bkn.go.id/.
2. Access the network tab in the browser's developer tools (usually accessible by pressing F12 or right-clicking and selecting "Inspect").
3. In the SSCASN website, perform a search for the desired major or use the filters available on the site.
4. In the network tab, look for requests to `api-sscasn.bkn.go.id`. These requests contain the parameters we need.
5. Examine the request URL and query parameters. You'll typically find:
   - `kode_ref_pend`: The education reference code
   - `pengadaan_kd`: The procurement code
   - `instansi_id`: The institution ID (if you're filtering by a specific institution)
6. Copy these values for use in the program.

Example of what you might see in the network tab:

```
https://api-sscasn.bkn.go.id/2024/portal/spf?kode_ref_pend=4480271&instansi_id=A5EB03E23AFBF6A0E040640A040252AD&pengadaan_kd=2&offset=0
```

In this URL:
- `kode_ref_pend=4480271`
- `instansi_id=A5EB03E23AFBF6A0E040640A040252AD`
- `pengadaan_kd=2`

You can use these values with the corresponding flags when running the scraper.

> [!NOTE] 
> The exact process might vary slightly depending on how you interact with the SSCASN website. Always ensure you're using the most recent and relevant parameters for your search.
