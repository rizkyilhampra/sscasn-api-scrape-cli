This project is used for personal and educational purposes. It is designed to scrape data from an API, parse it into an Excel format (.xlsx), and save it under the "data/" directory. Additionally, it implements concurrent processing using goroutines, so on my local computer, it can retrieve 500-600 records per minute, followed by detailed data from another endpoint.

## How to run

```bash
git clone https://github.com/rizkyilhampra/sscasn-scraper.git
cd sscasn-scraper
go mod tidy
```

```bash
go run mod.go -kodeRefPend=5102656 -namaJurusan="S1 Pendidikan Keagamaan Katolik" 
```

you can also add flag `-provinsi="Jawa Tengah"` if you want filter by *Instansi* contain the string.

or you maybe want to run with multiple kode ref pendidikan. Please take a look of `run_all_program.sh` and `data.json`. Ensure you have `jq` before you can run this.

## How to get KodeRefPend (Bahasa)

1. Buka browser dan masuk ke halaman https://sscasn.bkn.go.id/
2. Buka network tab di inspect element
2. Cari nama jurusan yang ingin diambil kode ref pendidikannya pada kolom search web tersebut 
3. Lihat request yang terjadi, cari request yang mengandung `kode_ref_pend` pada pathnya dan dari host `api-sscasn.bkn.go.id`
4. Copy kode ref pendidikan tersebut

