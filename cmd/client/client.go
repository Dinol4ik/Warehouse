package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/rpc"
)

type (
	Req struct {
		Data string
	}

	Result struct {
		Data string
	}

	ReserveItemRequest struct {
		Method string `json:"method"`
		Items  []Item `json:"items"`
	}

	Item struct {
		Code   string `json:"code"`
		Amount int    `json:"amount"`
	}

	UnreservedItemRequest struct {
		Method string `json:"method"`
		Items  []Item `json:"items"`
	}

	FetchWarehouseItemsRequest struct {
		Method string `json:"method"`
		ID     string `json:"id"`
	}

	FetchWarehouseItemsResponse struct {
		Items []ItemInfo `json:"items"`
	}

	ItemInfo struct {
		Name     string `json:"name"`
		Size     string `json:"size"`
		Code     string `json:"code"`
		Amount   int    `json:"amount"`
		Reserved int    `json:"reserved"`
	}
	FetchItemsByCodesRequest struct {
		Codes []string `json:"codes"`
	}
)

func main() {
	client, err := rpc.DialHTTPPath("tcp", "127.0.0.1:8080", "/rpc")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ReserveItemHandlerTest(client)
	FetchWarehouseItemsHandlerTest(client)
	UnReserveItemHandlerTest(client)
	FetchItemsByCodesHandlerTest(client)
}

func ReserveItemHandlerTest(client *rpc.Client) {
	marshal, err := json.Marshal(ReserveItemRequest{
		Items: []Item{
			{"MP002XW0YADJ", 10},
		},
		Method: "ReserveItemHandler",
	})
	args1 := &Req{string(marshal)}
	var reply *Result
	err = client.Call("WarehouseService.ReserveItemHandler", args1, &reply)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*reply)
}

func FetchWarehouseItemsHandlerTest(client *rpc.Client) {
	marshal, err := json.Marshal(FetchWarehouseItemsRequest{
		Method: "FetchWarehouseItemsHandler",
		ID:     "c8575d67-7cfb-48ff-b8ed-6f455a18cf05",
	})
	args1 := &Req{string(marshal)}
	var reply *Result
	err = client.Call("WarehouseService.FetchWarehouseItemsHandler", args1, &reply)

	if err != nil {
		log.Fatal(err)
	}

	unmarshal := FetchWarehouseItemsResponse{}
	err = json.Unmarshal([]byte(reply.Data), &unmarshal)
	if err != nil {
		return
	}
	fmt.Println(*reply)
}

func UnReserveItemHandlerTest(client *rpc.Client) {
	marshal, err := json.Marshal(UnreservedItemRequest{
		Items: []Item{
			{"MP002XW0L57W", 8},
		},
		Method: "UnReserveItemHandler",
	})

	args1 := &Req{string(marshal)}
	var reply *Result
	err = client.Call("WarehouseService.UnReserveItemHandler", args1, &reply)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*reply)
}
func FetchItemsByCodesHandlerTest(client *rpc.Client) {
	marshal, err := json.Marshal(FetchItemsByCodesRequest{
		Codes: []string{"RTLACI608901"},
	})
	args1 := &Req{string(marshal)}
	var reply *Result
	err = client.Call("WarehouseService.FetchItemsByCodesHandler", args1, &reply)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n%v\n", *reply)
}
