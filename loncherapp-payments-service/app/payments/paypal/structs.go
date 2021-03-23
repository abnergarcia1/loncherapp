package paypal

type oauthPaypalResponse struct {
	Scope       string `json:"scope"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	AppID       string `json:"app_id"`
	ExpiresIn   int32  `json:"expires_in"`
	Nonce       string `json:"nonce"`
}

//OrderRequest
type OrderRequest struct {
	Intent        string          `json:"intent"`
	PurchaseUnits []PurchaseUnits `json:"purchase_units"`
}

type PurchaseUnits struct {
	Amount Amount `json:"amount"`
}

type Amount struct {
	CurrencyCode string `json:"currency_code" bson:"currency_code"`
	Value        string `json:"value" bson:"value"`
}

type OrderResponse struct {
	ID     string       `json:"id"`
	Links  []OrderLinks `json:"links"`
	Status string       `json:"status"`
}

type OrderLinks struct {
	Href   string `json:"href" bson:"href"`
	Method string `json:"method" bson:"method"`
	Rel    string `json:"rel" bson:"rel"`
}
type Payer struct {
	EmailAddress string `json:"email_address"`
	Name         struct {
		GivenName string `json:"given_name"`
		Surname   string `json:"surname"`
	} `json:"name"`
	PayerID string `json:"payer_id"`
}

type OrderCapture struct {
	ID            string       `json:"id"`
	Links         []OrderLinks `json:"links"`
	Payer         Payer        `json:"payer"`
	PurchaseUnits []struct {
		Payments struct {
			Captures []struct {
				Amount           Amount       `json:"amount"`
				CreateTime       string       `json:"create_time"`
				DisbursementMode string       `json:"disbursement_mode"`
				FinalCapture     bool         `json:"final_capture"`
				ID               string       `json:"id"`
				Links            []OrderLinks `json:"links"`
				SellerProtection struct {
					DisputeCategories []string `json:"dispute_categories"`
					Status            string   `json:"status"`
				} `json:"seller_protection"`
				SellerReceivableBreakdown struct {
					GrossAmount Amount `json:"gross_amount" bson:"gross_amount"`
					NetAmount   Amount `json:"net_amount" bson:"net_amount"`
					PaypalFee   Amount `json:"paypal_fee" bson:"paypal_fee"`
				} `json:"seller_receivable_breakdown" bson:"seller_receivable_breakdown"`
				Status     string `json:"status" bson:"status"`
				UpdateTime string `json:"update_time" bson:"update_time"`
			} `json:"captures" bson:"captures"`
		} `json:"payments" bson:"payments"`
		ReferenceID string `json:"reference_id" bson:"reference_id"`
		Shipping    struct {
			Address struct {
				AddressLine1 string `json:"address_line_1" bson:"address_line_1"`
				AddressLine2 string `json:"address_line_2" bson:"address_line_2"`
				AdminArea1   string `json:"admin_area_1" bson:"admin_area_1"`
				AdminArea2   string `json:"admin_area_2" bson:"admin_area_2"`
				CountryCode  string `json:"country_code" bson:"country_code"`
				PostalCode   string `json:"postal_code" bson:"postal_code"`
			} `json:"address" bson:"address"`
		} `json:"shipping" bson:"shipping"`
	} `json:"purchase_units" bson:"purchase_units"`
	Status string `json:"status" bson:"status"`
}

type WebhookResponse struct {
	CreateTime   string       `json:"create_time" bson:"create_time"`
	EventType    string       `json:"event_type" bson:"event_type"`
	EventVersion string       `json:"event_version" bson:"event_version"`
	ID           string       `json:"id" bson:"pid"`
	Links        []OrderLinks `json:"links" bson:"links"`
	Resource     struct {
		Amount           Amount       `json:"amount" bson:"amount"`
		CreateTime       string       `json:"create_time" bson:"create_time"`
		FinalCapture     bool         `json:"final_capture" bson:"final_capture"`
		ID               string       `json:"id" bson:"resource_id"`
		Links            []OrderLinks `json:"links" bson:"links"`
		SellerProtection struct {
			DisputeCategories []string `json:"dispute_categories" bson:"dispute_categories"`
			Status            string   `json:"status" bson:"status"`
		} `json:"seller_protection" bson:"seller_protection"`
		SellerReceivableBreakdown struct {
			GrossAmount Amount `json:"gross_amount" bson:"gross_amount"`
			NetAmount   Amount `json:"net_amount" bson:"net_amount"`
			PaypalFee   Amount `json:"paypal_fee" bson:"paypal_fee"`
		} `json:"seller_receivable_breakdown" bson:"seller_receivable_breakdown"`
		Status     string `json:"status" bson:"status"`
		UpdateTime string `json:"update_time" bson:"update_time"`
	} `json:"resource" bson:"resource"`
	ResourceType    string `json:"resource_type" bson:"resource_type"`
	ResourceVersion string `json:"resource_version" bson:"resource_version"`
	Summary         string `json:"summary" bson:"summary"`
}
