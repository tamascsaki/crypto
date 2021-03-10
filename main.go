package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println()
	one, two, three := time.Now().Add(24*time.Hour), time.Now().Add(48*time.Hour), time.Now().Add(72*time.Hour)
	fmt.Println("           " + one.Format("02") + "      " + two.Format("02") + "      " + three.Format("02"))
	fmt.Println()
	for i, current := range currencies {
		currencies[i].Prediction, currencies[i].Balance = getPrediction(current)
	}
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Balance > currencies[j].Balance
	})
	for i, current := range currencies {
		if current.Balance < 1 && currencies[i-1].Balance > 0 {
			fmt.Println()
		}
		fmt.Print(current.Name + " - ")
		if current.Balance == 0 {
			fmt.Print("  ")
		} else if current.Balance > 0 {
			fmt.Print("▲ ")
		} else {
			fmt.Print("▼ ")
		}
		for i, percent := range current.Prediction {
			if percent >= 10 {
				printPercent(1, percent, current.Balance, i == 2)
			} else if percent >= 0 {
				printPercent(2, percent, current.Balance, i == 2)
			} else if percent <= -10 {
				printPercent(0, percent, current.Balance, i == 2)
			} else {
				printPercent(1, percent, current.Balance, i == 2)
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func getPrediction(current currency) (prediction []int, balance int) {
	resp, err := http.Get(current.WalletInvestor)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	source := string(html)
	var momentPrice float64
	if split := strings.Split(source, "Current price today: <span class=\"number\"><span class=\"green bignum\">▲</span>"); len(split) > 1 {
		momentPrice, err = strconv.ParseFloat(split[1][:current.Chars], 64)
	} else if split = strings.Split(source, "Current price today: <span class=\"number\"><span class=\"red bignum\">▼</span>"); len(split) > 1 {
		momentPrice, err = strconv.ParseFloat(split[1][:current.Chars], 64)
	} else {
		fmt.Println(current.Name)
		os.Exit(1)
	}
	if err != nil {
		fmt.Println(current.Name)
		panic(err)
	}
	first := time.Now().Add(24*time.Hour).Format("2006-01-02") + "</td><td class=\"w0\" data-col-seq=\"1\"><span class=\"mobileshow\">Price: </span><i class=\"fa fa-usd\" aria-hidden=\"true\"></i> "
	firstPrice, err := strconv.ParseFloat(strings.Split(source, first)[1][:current.Chars], 64)
	if err != nil {
		fmt.Println(current.Name)
		panic(err)
	}
	second := time.Now().Add(48*time.Hour).Format("2006-01-02") + "</td><td class=\"w0\" data-col-seq=\"1\"><span class=\"mobileshow\">Price: </span><i class=\"fa fa-usd\" aria-hidden=\"true\"></i> "
	secondPrice, err := strconv.ParseFloat(strings.Split(source, second)[1][:current.Chars], 64)
	if err != nil {
		fmt.Println(current.Name)
		panic(err)
	}
	third := time.Now().Add(72*time.Hour).Format("2006-01-02") + "</td><td class=\"w0\" data-col-seq=\"1\"><span class=\"mobileshow\">Price: </span><i class=\"fa fa-usd\" aria-hidden=\"true\"></i> "
	thirdPrice, err := strconv.ParseFloat(strings.Split(source, third)[1][:current.Chars], 64)
	if err != nil {
		fmt.Println(current.Name)
		panic(err)
	}
	prediction = []int{int(math.Floor(firstPrice/momentPrice*100) - 100), int(math.Floor(secondPrice/momentPrice*100) - 100), int(math.Floor(thirdPrice/momentPrice*100) - 100)}
	return prediction, (prediction[0] + prediction[1] + prediction[2]) / 3
}

func printPercent(spaces, percent, balance int, showBalance bool) {
	for i := 0; i < spaces; i++ {
		fmt.Print(" ")
	}
	if showBalance {
		fmt.Print(strconv.Itoa(percent) + "%  (" + strconv.Itoa(balance) + "%)")
	} else {
		fmt.Print(strconv.Itoa(percent) + "%  | ")
	}
}

type currency struct {
	Name           string
	WalletInvestor string
	Chars          int
	Prediction     []int
	Balance        int
}

var currencies = []currency{
	{
		Name:           "btc ",
		WalletInvestor: "https://walletinvestor.com/forecast/bitcoin-prediction",
		Chars:          8,
	},
	{
		Name:           "eth ",
		WalletInvestor: "https://walletinvestor.com/forecast/ethereum-prediction",
		Chars:          8,
	},
	{
		Name:           "ltc ",
		WalletInvestor: "https://walletinvestor.com/forecast/litecoin-prediction",
		Chars:          7,
	},
	{
		Name:           "bch ",
		WalletInvestor: "https://walletinvestor.com/forecast/bitcoin-cash-prediction",
		Chars:          7,
	},
	{
		Name:           "eos ",
		WalletInvestor: "https://walletinvestor.com/forecast/eos-prediction",
		Chars:          5,
	},
	{
		Name:           "xlm ",
		WalletInvestor: "https://walletinvestor.com/forecast/stellar-lumens-prediction",
		Chars:          5,
	},
	{
		Name:           "atom",
		WalletInvestor: "https://walletinvestor.com/forecast/cosmos-prediction",
		Chars:          5,
	},
	{
		Name:           "link",
		WalletInvestor: "https://walletinvestor.com/forecast/link-prediction",
		Chars:          6,
	},
	{
		Name:           "xtz ",
		WalletInvestor: "https://walletinvestor.com/forecast/tezos-prediction",
		Chars:          5,
	},
	{
		Name:           "vet ",
		WalletInvestor: "https://walletinvestor.com/forecast/vechain-prediction",
		Chars:          6,
	},
	{
		Name:           "icx ",
		WalletInvestor: "https://walletinvestor.com/forecast/icon-prediction",
		Chars:          5,
	},
	{
		Name:           "ada ",
		WalletInvestor: "https://walletinvestor.com/forecast/cardano-prediction",
		Chars:          5,
	},
	{
		Name:           "algo",
		WalletInvestor: "https://walletinvestor.com/forecast/algorand-prediction",
		Chars:          5,
	},
	{
		Name:           "neo ",
		WalletInvestor: "https://walletinvestor.com/forecast/neo-prediction",
		Chars:          6,
	},
	{
		Name:           "omg ",
		WalletInvestor: "https://walletinvestor.com/forecast/omg-network-prediction",
		Chars:          5,
	},
	{
		Name:           "mana",
		WalletInvestor: "https://walletinvestor.com/forecast/decentraland-prediction",
		Chars:          5,
	},
}
