#!/bin/bash

# Check if jq is installed
if ! command -v jq &> /dev/null
then
    echo "jq is not installed. Please install it using: sudo apt-get install jq"
    exit 1
fi

# Check if the data.json file exists
if [ ! -f "data.json" ]; then
    echo "data.json file not found"
    exit 1
fi

# Read the JSON file and run the Go command for each entry
jq -r 'to_entries[] | "\(.key)|\(.value)"' data.json | while IFS='|' read -r nama_jurusan kode_ref_pend
do
    echo "Running for: $nama_jurusan (Code: $kode_ref_pend)"
    go run main-with-detail-v5.go -kodeRefPend="$kode_ref_pend" -namaJurusan="$nama_jurusan"
 
    
    # Optional: Add a short pause between runs if needed
    # sleep 1
    
    echo "Completed: $nama_jurusan"
    echo "-----------------------------------"
done

echo "All programs completed"
