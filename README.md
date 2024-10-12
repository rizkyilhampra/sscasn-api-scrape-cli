# SSCASN API Scraper (CLI Based)

An personal and educational tool designed to scrape job vacancy data from the [SSCASN API](https://api-sscasn.bkn.go.id), which powers the [SSCASN portal](https://sscasn.bkn.go.id). The program fetches data, processes it, and exports the results into an Excel file `.xlsx` stored in the `data` directory. It leverages concurrent processing with goroutines, allowing up to 10 requests per second for efficient data retrieval.

> [!NOTE]
> This project is my first learning experience to explore Go and concurrency features. 

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
- **Enhanced Data Output**: Includes additional fields like `jumlah_ms` (count of passed verification) and also detail information of each institution in the exported Excel file.

## Prerequisites

Before running the SSCASN Scraper, ensure you have [**Go**](https://go.dev) installed and setup on your system.

## How to Run

1. Clone the repository:
   ```bash
   git clone https://github.com/rizkyilhampra/sscasn-scraper.git
   ```
2. Navigate to the project directory:
   ```bash
   cd sscasn-scraper
   ```
3. Install the required dependencies:
   ```bash
   go mod tidy
   ```
4. Run the program with specific parameters. For example, to fetch data for "S1 Pendidikan Keagamaan Katolik":
   ```bash
   go run main.go -kodeRefPend=5102656 -namaJurusan="S1 Pendidikan Keagamaan Katolik"
   ```
5. Optionally, filter results by province using the -provinsi flag:
   ```bash
   go run main.go -kodeRefPend=5102656 -namaJurusan="S1 Pendidikan Keagamaan Katolik" -provinsi="Jawa Tengah"
   ```

## Command-line Arguments

The program supports the following command-line arguments:

- `-kodeRefPend`: (Required) Code reference for education.
- `-namaJurusan`: (Required) Name of the major or field of study of `kodeRefPend`. It's stand for label or title uses for generate Excel filename.
- `-provinsi`: (Optional) Filter results by province. Example: `-provinsi="Jawa Timur"`.
- `-pengadaanKd`: (Optional) Procurement code. Default is 2.
- `-instansiId`: (Optional) Institution ID. Example: `-instansiId="A5EB03E23AFBF6A0E040640A040252AD"` for "Kementerian Lingkungan Hidup dan Kehutanan".

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
