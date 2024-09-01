package payments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/customer"
	sub "github.com/stripe/stripe-go/v74/subscription"
	"github.com/stripe/stripe-go/v74/webhook"
)

func createCustomer(email string) (*stripe.Customer, error) {
	params := &stripe.CustomerParams{
		Email: stripe.String(email),
	}
	return customer.New(params)
}

type CheckoutSessionRequest struct {
	CustomerID string `json:"customer_id"`
	PriceID    string `json:"price_id"`
}

// createCheckoutSession handles the creation of a Stripe Checkout Session
func CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	var req CheckoutSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Use the parsed customer_id and price_id
	customerID := req.CustomerID
	priceID := req.PriceID

	params := &stripe.CheckoutSessionParams{
		Customer:           stripe.String(customerID),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String("subscription"),
		SuccessURL: stripe.String(os.Getenv("URL") + "/success"),
		CancelURL:  stripe.String(os.Getenv("URL") + "/cancel"),
	}

	s, err := session.New(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(s)
}

// handleWebhook processes Stripe webhook events
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), "whsec_YourWebhookSecret")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "invoice.payment_succeeded":
		log.Println("Payment succeeded")
		// Handle successful payment
	case "customer.subscription.deleted":
		log.Println("Subscription canceled")
		// Handle subscription cancellation
	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

func GetCustomerByEmail(email string) (*stripe.Customer, error) {
	params := &stripe.CustomerListParams{
		Email: stripe.String(email),
	}
	i := customer.List(params)
	for i.Next() {
		return i.Customer(), nil
	}
	return nil, fmt.Errorf("no customer found with email: %s", email)
}

func GetSubscriptionByCustomerAndProduct(customerID string, productID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(customerID),
	}
	i := sub.List(params)

	for i.Next() {
		subscription := i.Subscription()
		for _, item := range subscription.Items.Data {
			if item.Price.Product.ID == productID {
				return subscription, nil
			}
		}
	}

	return nil, fmt.Errorf("no subscription found for customer with product ID: %s", productID)
}

func test_createCustomer() {
	newCustomer, err := createCustomer("test@gmail.com")
	if err != nil {
		log.Fatalf("Failed to create customer: %v", err)
	}

	// Print the customer ID
	log.Println("Customer ID:", newCustomer.ID)
}

func test_getSubscriptionByUserEmail() {
	customerEmail := "test@gmail.com" // Replace with the actual user's email
	cust, err := GetCustomerByEmail(customerEmail)
	if err != nil {
		log.Printf("Error retrieving customer: %v", err)
	} else {
		log.Printf("Customer ID: %s\n", cust.ID)
	}
	productID := "prod_QkzhvwCenEWmDY"
	// Then, get the subscription for the specific product
	subscription, err := GetSubscriptionByCustomerAndProduct(cust.ID, productID)
	if err != nil {
		log.Fatalf("Error retrieving subscription: %v", err)
	}

	fmt.Printf("Subscription ID: %s, Status: %s\n", subscription.ID, subscription.Status)
}

func Setup() {
	// Set your Stripe secret key
	stripe.Key = os.Getenv("STRIPE_SECRET")
	// test_createCustomer()
	// test_getSubscriptionByUserEmail()
}
