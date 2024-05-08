package main

import (
	"fmt" // Untuk output format dan cetak
	"net" // Untuk operasi jaringan
	"os"  // Untuk operasi sistem seperti keluar program
)

// Struktur Server yang menyimpan daftar koneksi klien
type Server struct {
	clients []net.Conn // Daftar koneksi klien
}

// Metode untuk memulai server dan mendengarkan koneksi baru
func (s *Server) start() {
	// Buat TCP listener pada port 3000
	listener, err := net.Listen("tcp", ":3000")
	defer listener.Close() // Menutup listener saat fungsi selesai

	// Jika terjadi kesalahan saat membuat listener
	if err != nil {
		fmt.Println(err.Error()) // Cetak pesan kesalahan
		os.Exit(0)               // Keluar dari program
	}

	// Loop untuk terus-menerus menerima koneksi klien
	for {
		conn, err := listener.Accept() // Terima koneksi baru dari klien
		if err != nil {                // Jika terjadi kesalahan saat menerima koneksi
			fmt.Println(err.Error()) // Cetak pesan kesalahan
			continue                  // Lanjutkan ke iterasi berikutnya
		}

		fmt.Println("New connection") // Cetak pesan saat ada koneksi baru

		// Tambahkan koneksi baru ke daftar klien
		s.clients = append(s.clients, conn)

		// Jalankan loop pembacaan untuk koneksi baru dalam goroutine
		go s.readLoop(conn)
	}
}

// Metode untuk membaca data dari klien dan meneruskannya ke klien lain
func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close() // Tutup koneksi saat fungsi selesai

	buf := make([]byte, 2048) // Buffer untuk menyimpan data yang diterima

	// Loop untuk terus-menerus membaca data dari koneksi
	for {
		n, err := conn.Read(buf) // Baca data dari koneksi
		if err != nil {          // Jika terjadi kesalahan saat membaca
			fmt.Println(err.Error()) // Cetak pesan kesalahan
			os.Exit(1)               // Keluar dari program
		}

		// Cetak jumlah data yang diterima dan alamat klien
		fmt.Printf("Menerima %d bytes paket file dari %s \n", n, conn.RemoteAddr())

		// Kirim data yang diterima ke semua klien lainnya kecuali sumbernya
		for _, client := range s.clients {
			// Hindari mengirim ke sumber data
			if client.RemoteAddr() == conn.RemoteAddr() {
				continue
			}

			// Kirim data ke klien
			_, err := client.Write(buf[:n])
			if err != nil { // Jika terjadi kesalahan saat menulis ke klien
				fmt.Println(err.Error()) // Cetak pesan kesalahan
				continue                   // Lanjutkan ke klien berikutnya
			}

			// Cetak pesan saat meneruskan data ke klien lain
			fmt.Printf("Meneruskan  %d bytes file ke %s \n", n, client.RemoteAddr())
		}
	}
}

// Fungsi utama yang menjadi titik entri program
func main() {
	server := &Server{} // Buat instance Server
	server.start()      // Mulai server
}