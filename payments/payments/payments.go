package payments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"notifiers/controllers/alertsController"
	"notifiers/middlewares/authMiddleware"
	"notifiers/services/userService"
	subscriptionUtils "notifiers/utils/subscription"

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
	log.Printf("kokot")
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), "whsec_df27cdea0a382282ba590e627ddb527358a1be26da6cbbe7819588df7f40e573")
	log.Printf("event", event)

	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	switch event.Type {
	case "invoice.payment_succeeded":
		log.Println("Payment succeeded")
		alertsController.Setup()
		// Handle successful payment
	case "customer.subscription.deleted":
		log.Println("Subscription canceled")
		// Handle subscription cancellation
	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

	w.WriteHeader(http.StatusOK)
}

type GetCustomerByEmailRequest struct {
	Email string `json:"email"`
}

func HandleGetCustomerByEmail(w http.ResponseWriter, r *http.Request) {
	var req GetCustomerByEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	email := req.Email
	log.Printf("email", email)
	// Fetch customer details
	customer, err := subscriptionUtils.GetCustomerByEmail(email)
	if err != nil {
		http.Error(w, "Customer not found", http.StatusNotFound)
		return
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Return the customer details as JSON
	if err := json.NewEncoder(w).Encode(customer); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
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
	cust, err := subscriptionUtils.GetCustomerByEmail(customerEmail)
	if err != nil {
		log.Printf("Error retrieving customer: %v", err)
	} else {
		log.Printf("Customer ID: %s\n", cust.ID)
	}
	productID := "prod_QkzhvwCenEWmDY"
	// Then, get the subscription for the specific product
	subscription, err := subscriptionUtils.GetSubscriptionByCustomerAndProduct(cust.ID, productID)
	if err != nil {
		log.Fatalf("Error retrieving subscription: %v", err)
	}

	fmt.Printf("Subscription ID: %s, Status: %s\n", subscription.ID, subscription.Status)
}

func CancelSubscription(w http.ResponseWriter, r *http.Request) {
	email, _ := r.Context().Value(authMiddleware.UserEmailKey).(string)
	cust, err := subscriptionUtils.GetCustomerByEmail(email)
	if err != nil {
		log.Printf("Error retrieving customer: %v", err)
	}
	UserSubscription := subscriptionUtils.UserSubscription[email]
	var subscription *stripe.Subscription
	if UserSubscription.SubscriptionType == "gold" {
		subscription, err = subscriptionUtils.GetSubscriptionByCustomerAndProduct(cust.ID, subscriptionUtils.Gold_productID)
		if err != nil {
			log.Fatalf("Error retrieving subscription: %v", err)
		}
	} else if UserSubscription.SubscriptionType == "diamond" {
		subscription, err = subscriptionUtils.GetSubscriptionByCustomerAndProduct(cust.ID, subscriptionUtils.Diamond_productID)
		if err != nil {
			log.Fatalf("Error retrieving subscription: %v", err)
		}
	}

	// Assuming the subscription ID is passed as a URL parameter
	subscriptionID := subscription.ID
	if subscriptionID == "" {
		http.Error(w, "Subscription ID is required", http.StatusBadRequest)
		return
	}

	params := &stripe.SubscriptionCancelParams{}
	_, err = sub.Cancel(subscriptionID, params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to cancel subscription: %v", err), http.StatusInternalServerError)
		return
	}
	user, err := userService.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	canAddAlert, subscriptionType := subscriptionUtils.CheckToAddAlert(user.ID, email)
	subscriptionUtils.UserSubscription[email] = subscriptionUtils.UserAlertInfo{
		CanAddAlert:      canAddAlert,
		SubscriptionType: subscriptionType,
	}

	w.Write([]byte("Subscription canceled successfully."))
}

func Setup() {
	// Set your Stripe secret key
	stripe.Key = os.Getenv("STRIPE_SECRET")
	// test_createCustomer()
	// test_getSubscriptionByUserEmail()
}
