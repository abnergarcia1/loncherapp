package paypal

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/paypal_payments"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	auth "bitbucket.org/edgelabsolutions/loncherapp-core/auth/jwt"

	"bitbucket.org/edgelabsolutions/loncherapp-core/tools"
)

type Handler struct {
	rd         auth.AuthInterface
	tk         auth.TokenInterface
	tools      tools.Tools
	Service    *Service
	HTTPClient *http.Client
}

const (
	IntentCapture = "CAPTURE"
)

func NewHandler(rd auth.AuthInterface, tk auth.TokenInterface) *Handler {
	return &Handler{rd,
		tk,
		tools.Tools{},
		NewService(context.Background()),
		&http.Client{}}
}

func (h *Handler) PayPalGetToken() (*oauthPaypalResponse, error) {
	endpoint := os.Getenv("PAYPAL_OAUTH_API")

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	clientId := os.Getenv("PAYPAL_CLIENT_ID")
	clientSecret := os.Getenv("PAYPAL_CLIENT_SECRET")
	basicOauthEncoded := base64.StdEncoding.EncodeToString([]byte(clientId + ":" + clientSecret))

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		log.WithFields(log.Fields{
			"clientId":     clientId,
			"clientSecret": clientSecret,
		}).Errorf("Error when trying to create menu item in DB: %v", err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", string(basicOauthEncoded)))
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	response, err := h.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"clientId":     clientId,
			"clientSecret": clientSecret,
		}).Errorf("an error ocurred when try to create request to paypal oath: %v", response.Status)
		return nil, errors.New(response.Status)
	}

	body, _ := ioutil.ReadAll(response.Body)

	paypalResponse := &oauthPaypalResponse{}

	err = json.Unmarshal(body, paypalResponse)
	if err != nil {
		return nil, err
	}

	return paypalResponse, nil
}

func (h *Handler) CreatePaypalOrder(c *gin.Context) {
	profileID, _ := strconv.Atoi(c.Param("profileID"))
	endpoint := os.Getenv("PAYPAL_ORDER_API")

	subsPrice, err := h.Service.GetAppParameterValue("subscription_value")
	if err != nil {
		log.WithField("Parameter", "subscription_value").Errorf("Error when trying to get param value from DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error when getting subs price : %v", err)})
		return
	}
	subsCurrency, err := h.Service.GetAppParameterValue("subscription_currency")
	if err != nil {
		log.WithField("Parameter", "subscription_currency").Errorf("Error when trying to get param value from DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error when getting subs currency : %v", err)})
		return
	}

	paypalOrder := &OrderRequest{
		Intent: IntentCapture,
		PurchaseUnits: []PurchaseUnits{
			{Amount: Amount{
				CurrencyCode: subsCurrency.Value,
				Value:        subsPrice.Value},
			},
		},
	}

	payload, _ := json.Marshal(paypalOrder)
	authObj, err := h.PayPalGetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("an error ocurred when try to get paypal auth : %v", err)})
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		log.WithFields(log.Fields{
			"paypalOrder": paypalOrder,
			"profileID":   profileID,
		}).Errorf("Error when trying to create menu item in DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("an error ocurred when try to create paypal order : %v", err)})
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authObj.AccessToken))

	response, err := h.HTTPClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"paypalOrder": paypalOrder,
			"profileID":   profileID,
		}).Errorf("Error when trying to create order request to paypal: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("an error ocurred when try to create paypal order : %v", err)})
		return
	}

	if response.StatusCode != http.StatusCreated {
		log.WithFields(log.Fields{
			"paypalOrder": paypalOrder,
			"profileID":   profileID,
		}).Errorf("an error ocurred when try to create request to paypal: %v", response.Status)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("an error ocurred when try to create request to paypal: %v", response.Status)})
		return
	}

	body, _ := ioutil.ReadAll(response.Body)

	paypalResponse := &OrderResponse{}

	err = json.Unmarshal(body, paypalResponse)

	amount, _ := strconv.ParseFloat(paypalOrder.PurchaseUnits[0].Amount.Value, 32)

	loncherappPayment := models.Payment{
		ID:         paypalResponse.ID,
		LoncheraID: int32(profileID),
		Amount:     float32(amount),
		Currency:   paypalOrder.PurchaseUnits[0].Amount.CurrencyCode,
		Status:     paypalResponse.Status,
		Type:       "Subscription",
	}

	err = h.Service.CreateProfilePaymentOrder(loncherappPayment)
	if err != nil {
		log.WithFields(log.Fields{
			"paypalOrder":       paypalOrder,
			"profileID":         profileID,
			"paypalResponse":    paypalResponse,
			"loncherappPayment": loncherappPayment,
		}).Errorf("an error ocurred when try to create paypal payment in DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("an error ocurred when try to create paypal payment in DB : %v", err)})
		return
	}

	c.JSON(http.StatusCreated, paypalResponse)
}

func (h *Handler) CapturePayPalOrder(c *gin.Context) {
	endpoint := os.Getenv("PAYPAL_ORDER_API")
	profileID, _ := strconv.Atoi(c.Param("profileID"))
	orderID := c.Param("orderID")

	endpoint = endpoint + "/" + orderID + "/capture"

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"orderID":   orderID,
			"profileID": profileID,
		}).Errorf("Error when trying to create request for capture paypal order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error when trying to create request for capture paypal order: %v", err)})
		return
	}

	authObj, err := h.PayPalGetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("an error ocurred when try to get paypal auth : %v", err)})
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authObj.AccessToken))

	response, err := h.HTTPClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"orderID":   orderID,
			"profileID": profileID,
		}).Errorf("Error when trying to create capture order request to paypal: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error when trying to create capture order request to paypal: %v", err)})
		return
	}

	if response.StatusCode != http.StatusCreated {
		log.WithFields(log.Fields{
			"orderID":   orderID,
			"profileID": profileID,
		}).Errorf("an error ocurred when try to post request to paypal: %v", response.Status)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("an error ocurred when try to post request to paypal: %v", response.Status)})
		return
	}

	body, _ := ioutil.ReadAll(response.Body)

	paypalResponse := &OrderCapture{}

	err = json.Unmarshal(body, paypalResponse)

	amount, _ := strconv.ParseFloat(paypalResponse.PurchaseUnits[0].Payments.Captures[0].Amount.Value, 32)

	loncherappPayment := models.Payment{
		ID:         paypalResponse.ID,
		LoncheraID: int32(profileID),
		Amount:     float32(amount),
		Currency:   paypalResponse.PurchaseUnits[0].Payments.Captures[0].Amount.CurrencyCode,
		Status:     paypalResponse.Status,
		Type:       "Subscription",
	}

	err = h.Service.CreateProfilePaymentOrder(loncherappPayment)
	if err != nil {
		log.WithFields(log.Fields{
			"profileID":         profileID,
			"paypalResponse":    paypalResponse,
			"loncherappPayment": loncherappPayment,
		}).Errorf("an error ocurred when try to create paypal payment in DB: %v", err)
		//c.JSON(http.StatusInternalServerError, fmt.Sprintf("an error ocurred when try to create paypal payment in DB : %v", err))
		return
	}

	err = h.Service.ActivateProfile(profileID)
	if err != nil {
		log.WithFields(log.Fields{
			"profileID":         profileID,
			"paypalResponse":    paypalResponse,
			"loncherappPayment": loncherappPayment,
		}).Errorf("an error ocurred when try to activate profile: %v", err)
	}

	c.JSON(http.StatusCreated, paypalResponse)
}

func (h *Handler) HandleWebhookRequest(c *gin.Context) {
	wbResponse := &WebhookResponse{}
	if err := c.ShouldBindJSON(&wbResponse); err != nil {
		payload, _ := ioutil.ReadAll(c.Request.Body)
		log.WithFields(log.Fields{
			"payload": payload,
		}).Errorf("an error ocurred when try to parse paypal webhook response: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.StoreWebhookResponseMongo(*wbResponse)
	if err != nil {
		log.WithFields(log.Fields{
			"whResponse": wbResponse,
		}).Errorf("an error ocurred when try to store to monto webhook response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, nil)
}
