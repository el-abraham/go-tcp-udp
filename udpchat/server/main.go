package main

import (
	"fmt" // Untuk fungsi format dan cetak
	"net" // Untuk operasi jaringan
)

// Struktur Client yang menyimpan informasi tentang alamat dan nama klien
type Client struct {
	addr *net.UDPAddr // Alamat UDP klien
	name string       // Nama klien
}

func main() {
	members := make(map[string]*Client) // Map/Dict yang menyimpan anggota obrolan dengan alamat sebagai kunci

	// Resolusi alamat UDP pada port 3001
	udpAddr, err := net.ResolveUDPAddr("udp4", ":3001")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Mendengarkan koneksi UDP pada alamat yang telah ditemukan
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer conn.Close() // Menutup koneksi UDP saat fungsi berakhir

	// Fungsi untuk mengirim pesan ke klien tertentu
	msg := func(addr *net.UDPAddr, message string) {
		conn.WriteToUDP([]byte(message), addr) // Mengirim pesan ke alamat tertentu
	}

	// Fungsi untuk mengirim pesan ke semua klien kecuali klien saat ini
	broadcast := func(currentMember *Client, message string) {
		for memberAddr, member := range members { // Iterasi melalui semua anggota
			if memberAddr != currentMember.addr.String() { // Kecuali klien saat ini
				msg(member.addr, message) // Kirim pesan
			}
		}
	}

	// Loop untuk terus-menerus menerima pesan dari klien
	for {
		buf := make([]byte, 1024)             // Buffer untuk menyimpan pesan yang diterima
		n, addr, err := conn.ReadFromUDP(buf) // Baca data dari koneksi UDP
		if err != nil {                       // Jika terjadi kesalahan saat membaca
			fmt.Println(err.Error()) // Cetak kesalahan
			continue                 // Lanjut ke iterasi berikutnya
		}

		// Periksa apakah klien sudah ada dalam member
		currentMember, ok := members[addr.String()]
		if ok {
			go broadcast(currentMember, fmt.Sprintf("[%s] : %s \n", currentMember.name, string(buf[:n]))) // Broadcast pesan
		} else {
			members[addr.String()] = &Client{addr: addr, name: string(buf[:n])}                                  // Tambahkan klien ke member
			fmt.Printf("%s join from %s \n", string(buf[:n]), addr.String())                                     // Cetak pesan bergabungnya klien
			go broadcast(members[addr.String()], fmt.Sprintf("[Server] : %s join the chat \n", string(buf[:n]))) // Siarkan bahwa klien telah bergabung
			msg(addr, fmt.Sprintf("[Server] : Welcome %s!  \n", string(buf[:n])))
		}
	}

}
