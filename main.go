package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"math"
	"context"

	"github.com/xuri/excelize/v2"
	"golang.org/x/time/rate"
)

const (
	baseURL   = "https://api-sscasn.bkn.go.id/2024/portal/spf"
	maxRetries = 5
	initialDelay = 500 * time.Millisecond
)

var headers = map[string]string{
	"accept":             "application/json, text/plain, */*",
	"accept-encoding":    "gzip, deflate, br, zstd",
	"accept-language":    "en-US,en;q=0.9,id-ID;q=0.8,id;q=0.7",
	"connection":         "keep-alive",
	"host":               "api-sscasn.bkn.go.id",
	"origin":             "https://sscasn.bkn.go.id",
	"referer":            "https://sscasn.bkn.go.id/",
	"sec-ch-ua":          "\"Not)A;Brand\";v=\"99\", \"Google Chrome\";v=\"114\", \"Chromium\";v=\"114\"",
	"sec-ch-ua-mobile":   "?1",
	"sec-ch-ua-platform": "\"Android\"",
	"sec-fetch-dest":     "empty",
	"sec-fetch-mode":     "cors",
	"sec-fetch-site":     "same-site",
	"user-agent":         "Mozilla/5.0 (Linux; Android 13; Pixel 7 Build/TQ3A.230805.001) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.110 Mobile Safari/537.36",
}

type Response struct {
	Data struct {
		Meta struct {
			Total int `json:"total"`
		} `json:"meta"`
		Data []map[string]interface{} `json:"data"`
	} `json:"data"`
}

type DetailResponse struct {
	Data struct {
		JobDesc            string `json:"job_desc"`
		Keahlian           string `json:"keahlian"`
		LinkWebInstansi    string `json:"link_web_instansi"`
		CallCenterInstansi string `json:"call_center_instansi"`
		MedsosInstansi     string `json:"medsos_instansi"`
		HelpdeskInstansi   string `json:"helpdesk_instansi"`
		KualifikasiPendidikan string `json:"pendidikan_nm"`
		SyaratAdmin        []struct {
			Syarat      string `json:"syarat"`
			IsMandatory string `json:"is_mandatory"`
		} `json:"syarat_admin"`
	} `json:"data"`
}

type Config struct {
	KodeRefPend  string
	NamaJurusan  string
	FilterLokasi string
	PengadaanKd  int
	InstansiId   string
	Client       *http.Client
	Limiter      *rate.Limiter
}

func backoff(attempt int) time.Duration {
	return time.Duration(math.Pow(2, float64(attempt))) * initialDelay
}

func fetchData(cfg *Config, offset int) (*Response, error) {
    url := fmt.Sprintf("%s?kode_ref_pend=%s&offset=%d&pengadaan_kd=%d", baseURL, cfg.KodeRefPend, offset, cfg.PengadaanKd)
    if cfg.InstansiId != "" {
        url += fmt.Sprintf("&instansi_id=%s", cfg.InstansiId)
    }

    log.Println(url)
    
	var data *Response
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		data, err = fetchJSON[Response](cfg, url)

		if err == nil {
			return data, nil
		}

		log.Printf("Error fetching data at offset %d, attempt %d/%d: %v\n", offset, attempt, maxRetries, err)
		time.Sleep(backoff(attempt))
	}

	return nil, fmt.Errorf("max retries reached for fetching data")
}

func fetchDetailData(cfg *Config, formasiID string) (*DetailResponse, error) {
	var detail *DetailResponse
	var url string = baseURL + "/" + formasiID
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		detail, err = fetchJSON[DetailResponse](cfg, url)

		if err == nil {
			return detail, nil
		}

		log.Printf("Error fetching detail data for formasi_id %s, attempt %d/%d: %v\n", formasiID, attempt, maxRetries, err)
		time.Sleep(backoff(attempt))
	}

	return nil, fmt.Errorf("max retries reached for fetching detail data")
}

func fetchJSON[T Response | DetailResponse](cfg *Config, url string) (*T, error) {
	err := cfg.Limiter.Wait(context.Background())
	if err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := cfg.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var response T
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	return &response, nil
}

func processData(cfg *Config, totalData int) ([]map[string]interface{}, error) {
	var completeData []map[string]interface{}
	var wg sync.WaitGroup

	numWorkers := 10
	jobs := make(chan int, totalData)
	results := make(chan map[string]interface{}, totalData)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go worker(cfg, &wg, jobs, results)
	}

	for offset := 0; offset < totalData; offset += 10 {
		jobs <- offset
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		completeData = append(completeData, result)
	}

	return completeData, nil
}

func worker(cfg *Config, wg *sync.WaitGroup, jobs <-chan int, results chan<- map[string]interface{}) {
	defer wg.Done()
	for offset := range jobs {
		data, err := fetchData(cfg, offset)
		if err != nil {
			log.Printf("Failed to fetch data at offset %d after retries: %v\n", offset, err)
			continue
		}

		for _, record := range data.Data.Data {
			if cfg.FilterLokasi == "" || strings.Contains(strings.ToLower(record["lokasi_nm"].(string)), strings.ToLower(cfg.FilterLokasi)) {
				formasiID := fmt.Sprintf("%v", record["formasi_id"])
				detailData, err := fetchDetailData(cfg, formasiID)
				if err != nil {
					log.Printf("Skipping record with formasi_id %s due to failed detail fetching: %v\n", formasiID, err)
					continue
				}

				record["job_desc"] = detailData.Data.JobDesc
				record["keahlian"] = detailData.Data.Keahlian
				record["link_web_instansi"] = detailData.Data.LinkWebInstansi
				record["call_center_instansi"] = detailData.Data.CallCenterInstansi
				record["medsos_instansi"] = detailData.Data.MedsosInstansi
				record["helpdesk_instansi"] = detailData.Data.HelpdeskInstansi
				record["syarat_admin"] = detailData.Data.SyaratAdmin
				record["kualifikasi_pendidikan"] = detailData.Data.KualifikasiPendidikan

				results <- record
			}
		}
	}
}

func writeToExcel(cfg *Config, completeData []map[string]interface{}, excelOutputFile string) error {
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{
		"ins_nm", "jp_nm", "formasi_nm", "jabatan_nm", "lokasi_nm", "jumlah_formasi", "jumlah_ms",
		"gaji_min", "gaji_max", "pengumuman", "job_desc", "keahlian", "link_web_instansi",
		"call_center_instansi", "medsos_instansi", "helpdesk_instansi", "syarat_admin", "kualifikasi_pendidikan",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 4)
		f.SetCellValue(sheet, cell, header)
	}

	f.SetCellValue(sheet, "A1", "updated_at")
	f.SetCellValue(sheet, "B1", time.Now().Format("2006-01-02 15:04:05"))
	f.SetCellValue(sheet, "A2", "updated_by")
	f.SetCellValue(sheet, "B2", "rizkyilhampra")

	for i, record := range completeData {
		gajiMin, _ := strconv.ParseFloat(record["gaji_min"].(string), 64)
		gajiMax, _ := strconv.ParseFloat(record["gaji_max"].(string), 64)

		f.SetCellValue(sheet, fmt.Sprintf("A%d", i+5), record["ins_nm"])
		f.SetCellValue(sheet, fmt.Sprintf("B%d", i+5), record["jp_nama"])
		f.SetCellValue(sheet, fmt.Sprintf("C%d", i+5), record["formasi_nm"])
		f.SetCellValue(sheet, fmt.Sprintf("D%d", i+5), record["jabatan_nm"])
		f.SetCellValue(sheet, fmt.Sprintf("E%d", i+5), record["lokasi_nm"])
		f.SetCellValue(sheet, fmt.Sprintf("F%d", i+5), record["jumlah_formasi"])
		f.SetCellValue(sheet, fmt.Sprintf("G%d", i+5), record["jumlah_ms"])
		f.SetCellValue(sheet, fmt.Sprintf("H%d", i+5), gajiMin)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", i+5), gajiMax)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", i+5), fmt.Sprintf("https://sscasn.bkn.go.id/detailformasi/%v", record["formasi_id"]))
		f.SetCellValue(sheet, fmt.Sprintf("K%d", i+5), record["job_desc"])
		f.SetCellValue(sheet, fmt.Sprintf("L%d", i+5), record["keahlian"])
		f.SetCellValue(sheet, fmt.Sprintf("M%d", i+5), record["link_web_instansi"])
		f.SetCellValue(sheet, fmt.Sprintf("N%d", i+5), record["call_center_instansi"])
		f.SetCellValue(sheet, fmt.Sprintf("O%d", i+5), record["medsos_instansi"])
		f.SetCellValue(sheet, fmt.Sprintf("P%d", i+5), record["helpdesk_instansi"])

		var syaratAdmin []string
		if syaratAdminData, ok := record["syarat_admin"].([]struct {
			Syarat      string `json:"syarat"`
			IsMandatory string `json:"is_mandatory"`
		}); ok {
			for _, syarat := range syaratAdminData {
				syaratAdmin = append(syaratAdmin, syarat.Syarat)
			}
		}

		f.SetCellValue(sheet, fmt.Sprintf("Q%d", i+5), strings.Join(syaratAdmin, ", "))
		f.SetCellValue(sheet, fmt.Sprintf("R%d", i+5), record["kualifikasi_pendidikan"])
	}

	for i := 1; i <= len(headers); i++ {
		col, _ := excelize.ColumnNumberToName(i)
		f.SetColWidth(sheet, col, col, 30)
	}

	return f.SaveAs(excelOutputFile)
}

func main() {
	kodeRefPend := flag.String("kodeRefPend", "", "Kode referensi pendidikan")
	namaJurusan := flag.String("namaJurusan", "", "Nama jurusan")
	filterLokasi := flag.String("provinsi", "", "Provinsi yang diinginkan. Contoh: -provinsi=\"Jawa Timur\"")
	pengadaanKd := flag.Int("pengadaanKd", 2, "Kode pengadaan")
	instansiId := flag.String("instansiId", "", "ID instansi. Contoh: -instansiID=\"A5EB03E23AFBF6A0E040640A040252AD\" (untuk Kementerian Lingkungan Hidup dan Kehutanan)")
	flag.Parse()

	if *kodeRefPend == "" || *namaJurusan == "" {
		log.Fatal("Mohon masukkan kodeRefPend dan namaJurusan")
	}

	cfg := &Config{
		KodeRefPend:  *kodeRefPend,
		NamaJurusan:  *namaJurusan,
		FilterLokasi: *filterLokasi,
		PengadaanKd:  *pengadaanKd,
		InstansiId:   *instansiId,
		Client:       &http.Client{ Timeout: 30 * time.Second, },
		Limiter:      rate.NewLimiter(rate.Every(100*time.Millisecond), 1), // 10 requests per second
	}

	log.Println("Memulai proses pengambilan data...")

	initialData, err := fetchData(cfg, 0)
	if err != nil {
		log.Fatal("Gagal mengambil data awal:", err)
	}

	totalData := initialData.Data.Meta.Total
	log.Printf("Total data ditemukan: %d\n", totalData)

	timestamp := time.Now().Format("20060102_150405")
	dataDir := "data"
	excelOutputFile := filepath.Join(dataDir, fmt.Sprintf("sscasn_data_%s.xlsx", strings.ReplaceAll(*namaJurusan, " ", "_")+"_"+timestamp))

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal("Error creating data directory:", err)
	}

	filteredData, err := processData(cfg, totalData)
	if err != nil {
		log.Fatal("Error processing data:", err)
	}

	log.Println("Membuat file Excel...")

	if err := writeToExcel(cfg, filteredData, excelOutputFile); err != nil {
		log.Fatal("Error writing to Excel:", err)
	}

	log.Printf("Berhasil mengambil data untuk jurusan: %s, kode referensi: %s\n", *namaJurusan, *kodeRefPend)
	if *filterLokasi != "" {
		log.Printf("Filter provinsi: %s\n", *filterLokasi)
	}
	log.Printf("Jumlah total data: %d\n", len(filteredData))
	log.Printf("Proses selesai! Data berhasil disimpan dalam file %s\n", excelOutputFile)
}
