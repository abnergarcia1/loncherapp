package paypal

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	models "bitbucket.org/edgelabsolutions/loncherapp-core/models/paypal_payments"

	sharedModels "bitbucket.org/edgelabsolutions/loncherapp-core/models/shared"

	log "github.com/sirupsen/logrus"

	"bitbucket.org/edgelabsolutions/loncherapp-core/db/sql"
)

//PaymentsService
type Service struct {
	ctx context.Context
	db  sql.StorageDB
}

func NewService(context context.Context) *Service {
	return &Service{
		ctx: context,
		db:  sql.NewClient(),
	}
}

func (s *Service) GetAppParameterValue(param string) (*sharedModels.AppConfigParam, error) {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return nil, err
	}
	defer func() {
		s.db.Disconnect()
	}()

	var paramResult sharedModels.AppConfigParam

	err := ses.QueryOne(&paramResult, `SELECT * FROM App_Configuration WHERE Parameter=?`, param)
	if err != nil {
		log.WithField("Parameter", param).Errorf("Error when trying to get param value from DB: %v", err)
		return nil, err
	}

	return &paramResult, nil
}

func (s *Service) CreateProfilePaymentOrder(payment models.Payment) error {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return err
	}
	defer func() {
		s.db.Disconnect()
	}()

	_, err := ses.Execute("INSERT INTO App_Payments(ID, Lonchera_ID, Amount, Currency, Status, Type, Created_At) VALUES(?,?,?,?,?,?,CURRENT_TIMESTAMP) ",
		payment.ID, payment.LoncheraID, payment.Amount, payment.Currency, payment.Status, payment.Type)
	if err != nil {
		log.WithFields(log.Fields{
			"payment": payment,
		}).Errorf("Error when trying to create payment order in DB: %v", err)
		return err
	}

	return nil
}

func (s *Service) ActivateProfile(profileID int) error {
	var ses = sql.NewClient()
	if err := ses.Connect(); err != nil {
		return err
	}
	defer func() {
		s.db.Disconnect()
	}()

	_, err := ses.Execute("UPDATE Loncheras SET Active=1, Membership_Due_Date=ADDDATE(Membership_Due_Date, INTERVAL 1 YEAR)  WHERE ID=?",
		profileID)
	if err != nil {
		log.WithFields(log.Fields{
			"profileID": profileID,
		}).Errorf("Error when trying to activate profile DB: %v", err)
		return err
	}

	return nil
}

func (s *Service) StoreWebhookResponseMongo(response WebhookResponse) error {
	mongoURI := os.Getenv("LONCHERAPP_MONGO_DB_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	loncherappDatabase := client.Database("loncherapp")
	paymentsCollection := loncherappDatabase.Collection("payments")

	_, err = paymentsCollection.InsertOne(ctx, response)
	if err != nil {
		return err
	}
	return nil
}
