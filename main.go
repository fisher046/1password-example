package main

import (
	"time"

	"github.com/1Password/connect-sdk-go/connect"
	"github.com/1Password/connect-sdk-go/onepassword"
	"k8s.io/klog/v2"
)

const vaultTitle = "huan-test"
const itemTitle = "autotest-1"

func main() {
	// Create a 1password client to connect to the connector running on AKS.
	// TODO: Investigate if it's thread safe.
	client, err := connect.NewClientFromEnvironment()
	if err != nil {
		klog.Error(err)
	}

	// Get Vault by name.
	vaults, err := client.GetVaultsByTitle(vaultTitle)
	if err != nil {
		klog.Error(err)
	}

	klog.Infof("vault: %v", vaults[0])

	item := &onepassword.Item{
		Title: itemTitle,

		// This is the type of Vault item.
		// TODO: Investigate the difference between Login and Password.
		Category: onepassword.Login,
		Fields: []*onepassword.ItemField{

			// concealed means password type which cannot be seen from the web console directly.
			// But we still can see the content by clicking the Reveal button...
			{Value: "O3ah~.i227KdjOJz~.24gqY6jOo1.xxxxx", Type: "concealed", Label: "client_secret"},
			{Value: "bc3c8e5d-266d-4277-ae7f-28b0b06f5xxx", Type: "STRING", Label: "client_id"},
		},

		// It means the item was created by a connector running on AKS.
		Tags: []string{"1password-connector"},
	}

	_, err = client.CreateItem(item, vaults[0].ID)
	if err != nil {
		klog.Error(err)
	}

	// The creation is async and we must wait for a few seconds.
	time.Sleep(5 * time.Second)

	item, err = client.GetItemByTitle(item.Title, vaults[0].ID)
	if err != nil {
		klog.Error(err)
	}

	klog.Infof("item: %v", item)
}
