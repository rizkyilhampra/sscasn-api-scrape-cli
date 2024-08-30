# SSCASN SCRAPER
This project is used for personal and educational purposes. It is designed to scrape data from an API of https://api-sscasn.bkn.go.id which is used for https://sscasn.bkn.go.id, parse it into an Excel format (.xlsx), and save it under the "data/" directory. It also implements job/worker and concurrent processing using goroutines which will run 10 request per second.

## How to run

```bash
git clone https://github.com/rizkyilhampra/sscasn-scraper.git
cd sscasn-scraper
go mod tidy
```

Example getting data from API for "S1 Pendidikan Keagamaan Katolik"

```bash
go run mod.go -kodeRefPend=5102656 -namaJurusan="S1 Pendidikan Keagamaan Katolik" 
```

You can also add flag in example `-provinsi="Jawa Tengah"` if you want filter by *Instansi* where is contain the string.

Or maybe you want to run with multiple `kodeRefPend`. Put your desired `kodeRefPend` and `namaJurusan` in `data.json`, then run `./run_all_programs.sh`. It will run one by one synchrounously. Ensure that you have `jq` installed in your system and make it script is executable

## How to get kodeRefPend (Bahasa)

1. Buka browser dan masuk ke halaman https://sscasn.bkn.go.id/
2. Buka network tab di inspect element
2. Cari nama jurusan yang ingin diambil kode ref pendidikannya pada kolom search web tersebut 
3. Lihat request yang terjadi, cari request yang mengandung `kode_ref_pend` pada pathnya dan dari host `api-sscasn.bkn.go.id`
4. Copy kode ref pendidikan tersebut

