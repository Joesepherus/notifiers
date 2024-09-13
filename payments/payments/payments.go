package payments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"tradingalerts/middlewares/authMiddleware"
	"tradingalerts/services/loggingService"
	"tradingalerts/services/userService"
	"tradingalerts/utils/subscriptionUtils"

	"os"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
	"github.com/stripe/stripe-go/v79/customer"
	sub "github.com/stripe/stripe-go/v79/subscription"
	"github.com/stripe/stripe-go/v79/webhook"
)

func CreateCustomer(email string) (*stripe.Customer, error) {
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
		log.Println("Invalid request payload")
		loggingService.LogToDB("ERROR", "Invalid request payload", r)
		http.Redirect(w, r, "/error?message=Invalid+request+payload", http.StatusSeeOther)
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
		SuccessURL: stripe.String(os.Getenv("URL") + "/subscription-success"),
		CancelURL:  stripe.String(os.Getenv("URL") + "/subscription-cancel"),
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

	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), os.Getenv("STRIPE_WEBHOOK_SECRET"))
	log.Print("event", event, err)

	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	switch event.Type {
	case "invoice.payment_succeeded":
		log.Println("Payment succeeded")

		var rawData map[string]interface{}
		err := json.Unmarshal(event.Data.Raw, &rawData)
		if err != nil {
			log.Fatalf("Failed to unmarshal event data: %v", err)
		}

		// Access the nested data
		// Access the email directly from rawData
		email, ok := rawData["customer_email"].(string)
		if !ok {
			log.Fatalf("Failed to get email from raw data")
		}
		log.Println("Email:", email)

		user, err := userService.GetUserByEmail(email)
		if err != nil {
			log.Print("Error finding user.")
		}
		canAddAlert, subscriptionType := subscriptionUtils.CheckToAddAlert(user.ID, user.Email)
		subscriptionUtils.UserSubscription[user.Email] = subscriptionUtils.UserAlertInfo{
			CanAddAlert:      canAddAlert,
			SubscriptionType: subscriptionType,
		}
		log.Print("canAddAlert", canAddAlert)
		log.Print("subscriptionType", subscriptionType)

		// Handle successful payment
	case "customer.subscription.deleted":
		log.Println("Subscription canceled")
		// Handle subscription cancellation
	default:
		log.Printf("Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

type GetCustomerByEmailRequest struct {
	Email string `json:"email"`
}

func HandleGetCustomerByEmail(w http.ResponseWriter, r *http.Request) {
	var req GetCustomerByEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Invalid request payload")
		loggingService.LogToDB("ERROR", "Invalid request payload", r)
		http.Redirect(w, r, "/error?message=Invalid+request+payload", http.StatusSeeOther)
		return
	}

	email := req.Email
	log.Print("email", email)
	// Fetch customer details
	customer, err := subscriptionUtils.GetCustomerByEmail(email)
	if err != nil {
		log.Println("Customer not found")
		loggingService.LogToDB("ERROR", "Customer not found", r)
		http.Redirect(w, r, "/error?message=Customer+not+found", http.StatusSeeOther)
		return
	}

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Return the customer details as JSON
	if err := json.NewEncoder(w).Encode(customer); err != nil {
		log.Println("Failed to encode response")
		loggingService.LogToDB("ERROR", "Failed to encode response", r)
		http.Redirect(w, r, "/error?message=Failed+to+encode+response", http.StatusSeeOther)
	}
}

func test_createCustomer() {
	newCustomer, err := CreateCustomer("test@gmail.com")
	if err != nil {
		log.Printf("Failed to create customer: %v", err)
		return
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
		log.Printf("Error retrieving subscription: %v", err)
	}

	fmt.Printf("Subscription ID: %s, Status: %s\n", subscription.ID, subscription.Status)
}

func CancelSubscription(w http.ResponseWriter, r *http.Request) {
	email, _ := r.Context().Value(authMiddleware.UserEmailKey).(string)
	cust, err := subscriptionUtils.GetCustomerByEmail(email)
	if err != nil {
		log.Printf("Error retrieving customer")
		loggingService.LogToDB("ERROR", "Error retrieving customer", r)
		http.Redirect(w, r, "/error?message=Error+retrieving+customer", http.StatusSeeOther)
		return
	}
	UserSubscription := subscriptionUtils.UserSubscription[email]
	var subscription *stripe.Subscription
	if UserSubscription.SubscriptionType == "gold" {
		subscription, err = subscriptionUtils.GetSubscriptionByCustomerAndProduct(cust.ID, subscriptionUtils.Gold_productID)
		if err != nil {
			log.Printf("Error retrieving subscription")
			loggingService.LogToDB("ERROR", "Error retrieving subscription", r)
			http.Redirect(w, r, "/error?message=Error+retrieving+subscription", http.StatusSeeOther)
			return
		}
	} else if UserSubscription.SubscriptionType == "diamond" {
		subscription, err = subscriptionUtils.GetSubscriptionByCustomerAndProduct(cust.ID, subscriptionUtils.Diamond_productID)
		if err != nil {
			log.Printf("Error retrieving subscription")
			loggingService.LogToDB("ERROR", "Error retrieving subscription", r)
			http.Redirect(w, r, "/error?message=Error+retrieving+subscription", http.StatusSeeOther)
			return
		}
	}

	// Assuming the subscription ID is passed as a URL parameter
	subscriptionID := subscription.ID
	if subscriptionID == "" {
		log.Printf("Subscription ID is required")
		loggingService.LogToDB("ERROR", "Subscription ID is required", r)
		http.Redirect(w, r, "/error?message=Subscription+ID+is+required", http.StatusSeeOther)
		return
	}

	params := &stripe.SubscriptionCancelParams{}
	_, err = sub.Cancel(subscriptionID, params)
	if err != nil {
		log.Printf("Failed to cancel subscription")
		loggingService.LogToDB("ERROR", "Failed to cancel subscription", r)
		http.Redirect(w, r, "/error?message=Failed+to+cancel+subscription", http.StatusSeeOther)
		return
	}
	user, err := userService.GetUserByEmail(email)
	if err != nil {
		log.Printf("User not found")
		loggingService.LogToDB("ERROR", "User not found", r)
		http.Redirect(w, r, "/error?message=User+not+found", http.StatusSeeOther)
		return
	}
	canAddAlert, subscriptionType := subscriptionUtils.CheckToAddAlert(user.ID, email)
	subscriptionUtils.UserSubscription[email] = subscriptionUtils.UserAlertInfo{
		CanAddAlert:      canAddAlert,
		SubscriptionType: subscriptionType,
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)

	w.Write([]byte("Subscription canceled successfully."))
}
