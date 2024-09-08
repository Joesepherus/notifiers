package subscriptionUtils

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

const TEST_EMAIL = "joes@joesexperiences.com"
const TEST_CUSTOMER_ID = "cus_QmbO6AUESvXbYE"

func TestGetCustomerByEmail(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
	SetupKey()
	// Replace with a test email address from your Stripe test environment
	testEmail := TEST_EMAIL

	customer, err := GetCustomerByEmail(testEmail)
	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, testEmail, customer.Email)
}

func TestGetSubscriptionByCustomerAndProduct(t *testing.T) {
	// Replace with test customer ID and product ID from your Stripe test environment
	testCustomerID := TEST_CUSTOMER_ID
	testProductID := Diamond_productID

	subscription, err := GetSubscriptionByCustomerAndProduct(testCustomerID, testProductID)
	assert.NoError(t, err)
	assert.NotNil(t, subscription)
	assert.Equal(t, testProductID, subscription.Items.Data[0].Price.Product.ID)
}

func TestCheckToAddAlert(t *testing.T) {
	// Replace with test user ID and email
	testUserID := 1
	testEmail := TEST_EMAIL

	canAddAlert, subscriptionType := CheckToAddAlert(testUserID, testEmail)
	assert.True(t, canAddAlert)
	assert.Equal(t, "silver", subscriptionType)
}
