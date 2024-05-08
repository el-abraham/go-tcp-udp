package main

import (
	"bufio" // Untuk buffered input/output
	"bytes" // Untuk manipulasi byte
	"encoding/gob" // Untuk encoding/decoding objek Go ke dalam format biner
	"flag" // Untuk argumen baris perintah
	"fmt" // Untuk output format dan cetak
	"net" // Untuk operasi jaringan
	"os" // Untuk operasi sistem seperti file dan keluar program
	"time" // Untuk operasi waktu
)

// Struktur File untuk mentransfer file melalui jaringan
type File struct {
	Name string // Nama file
	Buf  []byte // Isi file dalam bentuk byte
}

// Fungsi utama yang menjadi titik entri program
func main() {
	// Definisikan flag untuk mengirim dan menerima file
	send := flag.String("send", "", "a file") // Flag untuk mengirim file
	receiver := flag.Bool("receiver", false, "a receiver") // Flag untuk menjadi penerima

	flag.Parse() // Parse argumen baris perintah

	// Buat koneksi TCP ke server di port 3000
	conn, err := net.Dial("tcp", ":3000")
	if err != nil { // Jika terjadi kesalahan saat membuat koneksi
		fmt.Println(err.Error()) // Cetak kesalahan
	}

	defer conn.Close() // Tutup koneksi saat selesai

	// Jika flag receiver diaktifkan, terima file
	if *receiver {
		receiveFile(conn) // Panggil fungsi untuk menerima file
	} else {
		sendFile(conn, *send) // Panggil fungsi untuk mengirim file dengan nama file yang diberikan
	}
}

// Fungsi untuk menerima file melalui koneksi jaringan
func receiveFile(conn net.Conn) {
	defer conn.Close() // Tutup koneksi saat selesai

	fmt.Println("Menunggu file...") // Pesan untuk menunggu file

	buf := make([]byte, 2048) // Buffer untuk membaca data dari koneksi
	var combineBuf []byte // Buffer gabungan untuk menyimpan semua data yang diterima

	// Loop untuk terus-menerus menerima data
	for {
		n, err := conn.Read(buf) // Baca data dari koneksi
		if err != nil { // Jika terjadi kesalahan saat membaca
			fmt.Println(err.Error()) // Cetak kesalahan
			break // Keluar dari loop
		}

		// Tampilkan jumlah byte yang diterima
		fmt.Printf("Paket %d byte diterima \n", n)
		combineBuf = append(combineBuf, buf[:n]...) // Gabungkan data yang diterima

		dec := gob.NewDecoder(bytes.NewReader(combineBuf)) // Buat decoder GOB dari data yang diterima

		var file File // Objek file untuk menampung data yang didecode
		decodeErr := dec.Decode(&file) // Decode data menjadi objek file

		if decodeErr == nil { // Jika decoding berhasil
			// Buat file baru dengan nama unik berdasarkan waktu
			cFile, err := os.Create(fmt.Sprintf("%d%s", time.Now().UnixMilli(), file.Name))
			if err != nil { // Jika terjadi kesalahan saat membuat file
				fmt.Println(err.Error()) // Cetak kesalahan
				os.Exit(0) // Keluar dari program
			}

			// Tulis isi file ke dalam file yang baru dibuat
			cw, err := cFile.Write(file.Buf) 
			if err != nil { // Jika terjadi kesalahan saat menulis file
				fmt.Println(err.Error()) // Cetak kesalahan
				os.Exit(0) // Keluar dari program
			}

			// Tampilkan pesan keberhasilan
			fmt.Printf("File berhasil diterima dengan ukuran %d byte \n", cw)

			cFile.Sync() // Sinkronkan file
		}

	}
}

// Fungsi untuk mengirim file melalui koneksi jaringan
func sendFile(conn net.Conn, filename string) {
	defer conn.Close() // Tutup koneksi saat selesai

	file, err := os.Open(filename) // Buka file yang akan dikirim
	defer file.Close() // Tutup file saat selesai
	if err != nil { // Jika terjadi kesalahan saat membuka file
		fmt.Println(err.Error()) // Cetak kesalahan
		os.Exit(1) // Keluar dari program
	}

	// Ambil informasi tentang file (ukuran, nama, dll.)
	fileInfo, err := file.Stat()
	if err != nil { // Jika terjadi kesalahan saat mengambil informasi file
		fmt.Println(err.Error()) // Cetak kesalahan
		os.Exit(1) // Keluar dari program
	}

	buf := make([]byte, fileInfo.Size()) // Buffer untuk membaca isi file
	n, err := file.Read(buf) // Baca isi file
	if err != nil { // Jika terjadi kesalahan saat membaca file
		fmt.Println(err.Error()) // Cetak kesalahan
		os.Exit(1) // Keluar dari program
	}

	var encodebuf bytes.Buffer // Buffer untuk encoding data
	enc := gob.NewEncoder(&encodebuf) // Buat encoder GOB

	// Encode informasi file dan isi file ke dalam buffer
	encErr := enc.Encode(File{Name: filename, Buf: buf[:n]})
	if encErr != nil { // Jika terjadi kesalahan saat encoding
		fmt.Println(encErr.Error()) // Cetak kesalahan
		os.Exit(1) // Keluar dari program
	}

	// Kirim data yang di-encode ke koneksi jaringan
	wn, err := conn.Write(encodebuf.Bytes())
	if err != nil { // Jika terjadi kesalahan saat menulis ke koneksi
		fmt.Println(err.Error()) // Cetak kesalahan
		os.Exit(1) // Keluar dari program
	}

	// Pesan keberhasilan saat mengirim file
	fmt.Printf("Package berhasil dikirim! %d byte \n", wn)

	// Menunggu input pengguna untuk mencegah program keluar terlalu cepat
	bufio.NewReader(os.Stdin).ReadBytes('\n') 
}