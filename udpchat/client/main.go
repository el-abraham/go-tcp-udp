package main

import (
	"bufio"   // Untuk pembacaan buffered
	"flag"    // Untuk argumen baris perintah
	"fmt"     // Untuk format dan output
	"net"     // Untuk operasi jaringan
	"os"      // Untuk operasi sistem seperti keluar program
	"strings" // Untuk manipulasi string
)

func main() {
	nick := flag.String("nick", "Matt", "a nickname") // Definisikan flag "-nick" untuk nama panggilan
	flag.Parse()                                      // Parse flag dari baris perintah

	udpAddr, err := net.ResolveUDPAddr("udp", ":3001") // Cari alamat UDP di port 3001
	if err != nil {                                    // Jika ada kesalahan
		fmt.Println(err.Error()) // Cetak pesan kesalahan
		os.Exit(1)               // Keluar dari program dengan kode kesalahan
	}

	conn, err := net.DialUDP("udp", nil, udpAddr) // Buat koneksi UDP ke alamat yang ditemukan
	if err != nil {                               // Jika ada kesalahan
		fmt.Println(err.Error()) // Cetak pesan kesalahan
		os.Exit(1)               // Keluar dari program dengan kode kesalahan
	}

	_, err = conn.Write([]byte(*nick)) // Kirim nama panggilan ke koneksi UDP
	if err != nil {                    // Jika terjadi kesalahan saat menulis ke koneksi
		fmt.Println(err.Error()) // Cetak pesan kesalahan
		os.Exit(1)               // Keluar dari program
	}

	go handleReceiveMsg(conn) // Jalankan fungsi penerima pesan dalam goroutine

	// Loop untuk terus-menerus membaca input pengguna
	for {
		reader := bufio.NewReader(os.Stdin) // Pembaca buffered dari input standar
		text, _ := reader.ReadString('\n')  // Baca input hingga newline

		fmt.Printf("\033[1A\033[K")    // Kode ANSI untuk naik satu baris dan menghapusnya
		fmt.Printf("[You] : %s", text) // Tampilkan teks dengan prefiks "[You]"

		_, err := conn.Write([]byte(strings.Trim(text, "\r\n"))) // Kirim teks setelah menghapus karakter newline
		if err != nil {                                          // Jika ada kesalahan saat mengirim pesan
			fmt.Println(err) // Cetak kesalahan
			break            // Keluar dari loop
		}
	}
}

// Fungsi untuk menangani penerimaan pesan dari koneksi UDP
func handleReceiveMsg(conn *net.UDPConn) {
	// Loop untuk terus-menerus menerima pesan
	for {
		buf := make([]byte, 1024)          // Buffer untuk menyimpan data yang diterima
		n, _, err := conn.ReadFromUDP(buf) // Baca data dari koneksi UDP
		if err != nil {                    // Jika terjadi kesalahan saat membaca
			fmt.Println(err.Error()) // Cetak pesan kesalahan
			break                    // Keluar dari loop
		}

		fmt.Print(string(buf[:n])) // Cetak pesan yang diterima
	}
}
