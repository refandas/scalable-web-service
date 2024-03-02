package main

import (
	"fmt"
	"os"
	"strconv"
)

type Trainee struct {
	Name    string
	Address string
	Job     string
	Reason  string
}

func (trainee *Trainee) String() string {
	return fmt.Sprintf("Nama: %s\nAlamat : %s\nPekerjaan: %s\nAlasan: %s",
		trainee.Name, trainee.Address, trainee.Job, trainee.Reason)
}

func showTraineeData(number int) {
	traineeData := []Trainee{
		{
			Name:    "Irfan ",
			Address: "Trenggalek, Jawa Tengah",
			Job:     "Software Engineer",
			Reason:  "Meningkatkan keahlian dalam membangun aplikasi web yang dapat menangani beban kerja besar dengan Go.",
		},
		{
			Name:    "Anisa",
			Address: "Surabaya, Jawa Timur",
			Job:     "Web Developer",
			Reason:  "Ingin memperdalam pengetahuan dalam mengembangkan aplikasi web yang scalable menggunakan Go.",
		},
		{
			Name:    "Budi",
			Address: "Yogyakarta, DI Yogyakarta",
			Job:     "IT Manager",
			Reason:  "Belajar strategi terbaik dalam merancang arsitektur web yang scalable dengan Go, untuk meningkatkan efisiensi sistem di perusahaan.",
		},
	}

	if number < 0 || number >= len(traineeData) {
		fmt.Println("Tidak ada siswa dengan nomer absen tersebut.")
		return
	}

	fmt.Println(traineeData[number].String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Argument tidak tepat.\nPenggunaan: go run main.go <number>")
		return
	}

	number, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Argument tidak tepat, masukkan nilai integer.")
		return
	}

	showTraineeData(number)
}
