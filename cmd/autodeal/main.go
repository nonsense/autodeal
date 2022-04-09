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

	wallet       string
	maddr        string
	inputdataURL string
	piecesize    int
	index        int
	verified     bool
	price        int64
)

func init() {
	flag.StringVar(&wallet, "wallet", "", "wallet to be used for the deal")
	flag.StringVar(&maddr, "maddr", "f0127896", "miner address on-chain")
	flag.IntVar(&piecesize, "piecesize", 2, "piece size in GB")
	flag.IntVar(&index, "index", 0, "file index")
	flag.BoolVar(&verified, "verified", false, "")
	flag.StringVar(&inputdataURL, "inputdata-url", "https://anton-public-bucket-boost.s3.eu-central-1.amazonaws.com/spx-notes.json", "input data (fixtures)")
	flag.Int64Var(&price, "price-per-epoch", 1, "price-per-epoch for deal")
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

	args := []string{
		"deal",
		fmt.Sprintf("--verified=%t", verified),
		fmt.Sprintf("--provider=%s", maddr),
		fmt.Sprintf("--http-url=%s", d.URL),
		fmt.Sprintf("--commp=%s", d.CommP),
		fmt.Sprintf("--car-size=%d", d.CarSize),
		fmt.Sprintf("--piece-size=%d", d.PieceSize),
		fmt.Sprintf("--payload-cid=%s", d.PayloadCID),
		fmt.Sprintf("--storage-price-per-epoch=%d", price),
	}

	if wallet != "" {
		args = append(args,
			fmt.Sprintf("--wallet=%s", wallet),
		)
	}

	out, err := exec.Command("boost", args...).CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
