package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	datatransfer "google.golang.org/api/admin/datatransfer/v1"
)

var (
	fromAddr string
	toAddr   string

	// calendarAppID   int64 = 435070579839
	docsAppID       int64 = 55656082996
	docTransfers          = 0
	docTransfersMax       = 5

	transferCmd = &cobra.Command{
		Use:   "transfer",
		Short: "Transfer user data",
		Run:   transferRunFunc,
	}
)

func init() {
	rootCmd.AddCommand(transferCmd)

	transferCmd.Flags().StringVarP(&fromAddr, "from", "f", "", "source email address for doc transfer")
	transferCmd.Flags().StringVarP(&toAddr, "to", "t", "", "destination email address for doc transfer")
}

func transferRunFunc(cmd *cobra.Command, args []string) {
	if fromAddr == "" || toAddr == "" {
		exitWithError("must provide --from and --to")
	}
	fmt.Printf("document transfer: %s --> %s\n", fromAddr, toAddr)
	fromID, toID := getUserIDs(fromAddr, toAddr)

	dtc, err := newDataTransferClient()
	if err != nil {
		exitWithError(fmt.Sprintf("unable to create client: %s", err))
	}

	// these should be stable, but if you need to look up appIDs this
	// is how it's done...
	// as := datatransfer.NewApplicationsService(dtc)
	// ar, err := as.List().Do()
	// if err != nil {
	// 	exitWithError(err.Error())
	// }
	// var appID int64
	// for _, i := range ar.Applications {
	// 	if i.Name == "Drive and Docs" {
	// 		appID = i.Id
	// 	}
	// }
	// if appID == 0 {
	// 	exitWithError("application id not found")
	// }

	transferDocs(dtc, fromID, toID)
	// maybe?
	//transferCalendar(...)
}

func getUserIDs(fromAddr, toAddr string) (string, string) {
	ac, err := newAdminClient()
	if err != nil {
		exitWithError(fmt.Sprintf("unable to create client: %s", err))
	}
	from, err := ac.Users.Get(fromAddr).Do()
	if err != nil {
		exitWithError(fmt.Sprintf("unable to get ID for %s: %v", fromAddr, err))
	}
	to, err := ac.Users.Get(toAddr).Do()
	if err != nil {
		exitWithError(fmt.Sprintf("unable to get ID for %s: %v", toAddr, err))
	}
	return from.Id, to.Id
}

func transferDocs(dtc *datatransfer.Service, fromID, toID string) {
	// https://developers.google.com/admin-sdk/data-transfer/v1/parameters
	p := datatransfer.ApplicationTransferParam{
		Key: "PRIVACY_LEVEL",
		Value: []string{
			"PRIVATE",
			"SHARED",
		},
	}
	adt := datatransfer.ApplicationDataTransfer{
		ApplicationId:             docsAppID,
		ApplicationTransferParams: []*datatransfer.ApplicationTransferParam{&p},
	}
	t := datatransfer.DataTransfer{
		ApplicationDataTransfers: []*datatransfer.ApplicationDataTransfer{&adt},
		NewOwnerUserId:           toID,
		OldOwnerUserId:           fromID,
	}
	ts := datatransfer.NewTransfersService(dtc)
	tr, err := ts.Insert(&t).Do()
	if err != nil {
		if docTransfers == docTransfersMax {
			exitWithError(err.Error())
		}
		fmt.Printf("retry %v/%v\n", docTransfers+1, docTransfersMax)
		docTransfers = docTransfers + 1
		time.Sleep(5 * time.Second)
		transferDocs(dtc, fromID, toID)
	}
	count := 1
	for {
		time.Sleep(5 * time.Second)
		res, _ := ts.Get(tr.Id).Do()
		if count == 5 {
			if res != nil && res.OverallTransferStatusCode == "inProgress" {
				fmt.Println("transfer running long")
				break
			} else {
				if docTransfers == docTransfersMax {
					fmt.Printf("transfer failed (status: %v)\n", res.OverallTransferStatusCode)
					break
				}
				fmt.Printf("retry %v/%v\n", docTransfers+1, docTransfersMax)
				docTransfers = docTransfers + 1
				transferDocs(dtc, fromID, toID)
			}
		}
		if res != nil && res.OverallTransferStatusCode == "completed" {
			fmt.Println("transfer complete")
			break
		}
		count = count + 1
	}
}
