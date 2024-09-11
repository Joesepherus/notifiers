package subscriptionUtils

import (
	"fmt"
	"log"
	"os"
	"tradingalerts/services/alertsService"
	"tradingalerts/services/userService"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	sub "github.com/stripe/stripe-go/v74/subscription"
)

// TODO: add logic for when user is subscribed, so he can have
// more than 5 active alerts
var Gold_productID string = "prod_QkzhvwCenEWmDY"
var Diamond_productID string = "prod_QlltE9sAx7aY9z"

type UserAlertInfo struct {
	CanAddAlert      bool
	SubscriptionType string
}

var UserSubscription = make(map[string]UserAlertInfo)

var SUBSCRIPTION_LIMITS = map[string]int{
	"silver":  5,
	"gold":    100,
	"diamond": 1000,
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

func CheckToAddAlert(userID int, email string) (bool, string) {
	alerts, _ := alertsService.GetAlertsByUserID(userID)

	cust, err := GetCustomerByEmail(email)
	if err != nil {
		log.Printf("Error retrieving customer: %v", err)
		return false, ""
	}
	gold_subscription, err := GetSubscriptionByCustomerAndProduct(cust.ID, Gold_productID)
	diamond_subscription, err2 := GetSubscriptionByCustomerAndProduct(cust.ID, Diamond_productID)

	if err == nil && gold_subscription.Status == "active" {
		if len(alerts) >= SUBSCRIPTION_LIMITS["gold"] {
			return false, ""
		} else {
			return true, "gold"
		}
	}

	if err2 == nil && diamond_subscription.Status == "active" {
		if len(alerts) >= SUBSCRIPTION_LIMITS["diamond"] {
			return false, ""
		} else {
			return true, "diamond"
		}
	}
	if len(alerts) >= SUBSCRIPTION_LIMITS["silver"] {
		return false, "silver"
	}
	return true, "silver"
}

func SetupKey() {
	// Set your Stripe secret key
	stripe.Key = os.Getenv("STRIPE_SECRET")
}

func Setup() {
	SetupKey()
	users, _ := userService.GetUsers()
	for _, user := range users {
		canAddAlert, subscriptionType := CheckToAddAlert(user.ID, user.Email)
		UserSubscription[user.Email] = UserAlertInfo{
			CanAddAlert:      canAddAlert,
			SubscriptionType: subscriptionType,
		}
	}
}
