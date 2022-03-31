package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

var (
	args = map[int][]DealArgs{}

	maddr        string
	inputdataURL string
	piecesize    int
	index        int
)

func init() {
	flag.StringVar(&maddr, "maddr", "f0127896", "miner address on-chain")
	flag.IntVar(&piecesize, "piecesize", 2, "piece size in GB")
	flag.IntVar(&index, "index", 0, "file index")
	flag.StringVar(&inputdataURL, "inputdata-url", "https://anton-public-bucket-boost.s3.eu-central-1.amazonaws.com/spx-notes.json", "input data (fixtures)")
}

type DealArgs struct {
	URL        string
	CommP      string
	PieceSize  uint64
	CarSize    uint64
	PayloadCID string
}

func readInputData() {
	resp, err := http.Get(inputdataURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &args)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	readInputData()

	d := args[piecesize][index]

	out, err := exec.Command("boost", "deal",
		"--verified=false",
		"--quiet",
		fmt.Sprintf("--provider=%s", maddr),
		fmt.Sprintf("--http-url=%s", d.URL),
		fmt.Sprintf("--commp=%s", d.CommP),
		fmt.Sprintf("--car-size=%d", d.CarSize),
		fmt.Sprintf("--piece-size=%d", d.PieceSize),
		fmt.Sprintf("--payload-cid=%s", d.PayloadCID),
	).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
