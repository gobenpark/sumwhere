package models

type Receipt struct {
	Receipt struct {
		ReceiptType                string `json:"receipt_type"`
		AdamID                     int    `json:"adam_id"`
		AppItemID                  int    `json:"app_item_id"`
		BundleID                   string `json:"bundle_id"`
		ApplicationVersion         string `json:"application_version"`
		DownloadID                 int    `json:"download_id"`
		VersionExternalIdentifier  int    `json:"version_external_identifier"`
		ReceiptCreationDate        string `json:"receipt_creation_date"`
		ReceiptCreationDateMs      string `json:"receipt_creation_date_ms"`
		ReceiptCreationDatePst     string `json:"receipt_creation_date_pst"`
		RequestDate                string `json:"request_date"`
		RequestDateMs              string `json:"request_date_ms"`
		RequestDatePst             string `json:"request_date_pst"`
		OriginalPurchaseDate       string `json:"original_purchase_date"`
		OriginalPurchaseDateMs     string `json:"original_purchase_date_ms"`
		OriginalPurchaseDatePst    string `json:"original_purchase_date_pst"`
		OriginalApplicationVersion string `json:"original_application_version"`
		InApp                      []struct {
			Quantity                string `json:"quantity"`
			ProductID               string `json:"product_id"`
			TransactionID           string `json:"transaction_id"`
			OriginalTransactionID   string `json:"original_transaction_id"`
			PurchaseDate            string `json:"purchase_date"`
			PurchaseDateMs          string `json:"purchase_date_ms"`
			PurchaseDatePst         string `json:"purchase_date_pst"`
			OriginalPurchaseDate    string `json:"original_purchase_date"`
			OriginalPurchaseDateMs  string `json:"original_purchase_date_ms"`
			OriginalPurchaseDatePst string `json:"original_purchase_date_pst"`
			IsTrialPeriod           string `json:"is_trial_period"`
		} `json:"in_app"`
	} `json:"receipt"`
	Status      int    `json:"status"`
	Environment string `json:"environment"`
}
