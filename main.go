package main

import (
	// "github.com/dli-invest/finreddit/pkg/discord"
	"github.com/piquette/finance-go/quote"
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"os"
	"bytes"
	"strings"
)

func get_content() ([] string) {
    url := "https://raw.githubusercontent.com/FriendlyUser/cad_tickers_list/main/static/latest/raw_tickers.json"

    res, err := http.Get(url)
	if err != nil {

	}
    // perror(err)
    defer res.Body.Close()

    decoder := json.NewDecoder(res.Body)
    var symbols []string
    err = decoder.Decode(&symbols)
	if err != nil {

	}
    // if err != nil {
    //     fmt.Printf("%T\n%s\n%#v\n",err, err, err)
    //     switch v := err.(type){
    //         case *json.SyntaxError:
    //             fmt.Println(string(body[v.Offset-40:v.Offset]))
    //     }
    // }
    return symbols
}

func main() {
	var symbols []string = get_content()
    // getJson(, symbols)
	// data, _ := json.Marshal([]*Bird{pigeon, pigeon})
	iter := quote.List(symbols)
	var stock_data [][]string
	var used_symbols []string
	// columns := [8]string{"Last Price", 
	// "Change", "Volume", "Avg Vol (3 Month)", "Vol Ratio", 
	// "Dollar", "Market", "Exchange"}
	stock_data = append(stock_data, []string{"Symbol", "Last Price", 
	"Change",  "Vol Ratio"})
	for iter.Next() {
		q := iter.Quote()
		volume_ratio := float64(q.RegularMarketVolume) / float64(q.AverageDailyVolume3Month)
		used_symbols = append(used_symbols, q.Symbol)
		stock_data = append(stock_data, []string{ 
			q.Symbol,
			fmt.Sprintf("%2.2f", q.RegularMarketPrice),   
			fmt.Sprintf("%2.2f", q.RegularMarketChangePercent), 
			fmt.Sprintf("%2.2f", volume_ratio),
		})

		// fmt.Println(stock_data)
	}
	// send header
	var send_str string = "```"
	for i, s := range stock_data {
		temp_str := strings.Join(s, "\t")
		send_str = send_str + "\n" + temp_str
		if i % 10 == 0 && i != 0 {
			send_str = send_str + "```"
			resp, err := SendWebhook(send_str)
			fmt.Println(resp)
			fmt.Println(err)
			// send to discord
			send_str = "```"
		}
	}
	// result1 := strings.Join(stock_data, " ")
}

type DiscordPayload struct {
	Content string `json:"content,omitempty" xml:"content,omitempty" form:"content,omitempty" query:"content,omitempty"`
	// Embeds []DiscordEmbed `json:"embeds,omitempty" xml:"embeds,omitempty" form:"embeds,omitempty" query:"embeds,omitempty"`
}

func SendWebhook(stockData string) (*http.Response, error){
	discordUrl := os.Getenv("DISCORD_WEBHOOK")
	if discordUrl == "" {
		log.Fatal("DISCORD_WEBHOOK not set")
	}
	payload := &DiscordPayload{Content: stockData}
	webhookData, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(discordUrl, "application/json", bytes.NewBuffer(webhookData))
	return resp, err
}
