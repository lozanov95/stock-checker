package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

type CheckResponse struct {
	Available bool
	Price     float64
}

func HandleRequest() {
	oc := NewOzoneChecker()

	ozoneItems := []string{
		"https://www.ozone.bg/product/geyming-monitor-samsung-odyssey-g3-24g30a-24-va-144hz-1ms/",
		"https://www.ozone.bg/product/monitor-samsung-27g500-27-ips-led-165-hz-1-ms-gtg-2560x1440/",
	}

	for _, url := range ozoneItems {
		resp, err := oc.Check(url)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("=---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------==================================")
		if resp.Available {
			fmt.Printf("ITEM AVAILABLE!\n\t%s\n\t%.2f lv\n", url, resp.Price)
		} else {
			fmt.Printf("NOT AVAILABLE!\n\t%s", url)
		}
	}
}

func main() {
	lambda.Start(HandleRequest)
}
