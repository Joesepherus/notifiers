function showModal() {
  const modal = document.getElementById("signup-modal");
  modal.classList.add("show");
  window.onclick = function (event) {
    if (event.target.classList.contains("modal")) {
      closeSignUpForm();
    }
  };
}

function closeSignUpForm() {
  const modal = document.getElementById("signup-modal");
  modal.classList.remove("show");
}

function showSignUpForm() {
  const signupForm = document.getElementById("signup-form");
  const loginForm = document.getElementById("login-form");
  signupForm.classList.add("active-form");
  loginForm.classList.remove("active-form");
}

function showLoginForm() {
  const signupForm = document.getElementById("signup-form");
  const loginForm = document.getElementById("login-form");
  loginForm.classList.add("active-form");
  signupForm.classList.remove("active-form");
}

function openModalShowLoginForm() {
  showModal();
  showLoginForm();
}

function openModalShowSignUpForm() {
  showModal();
  showSignUpForm();
}

function closeSubscribeForm() {
  const modal = document.getElementById("subscribe-modal");
  modal.classList.remove("show");
}

function showSubscribeModal() {
  const modal = document.getElementById("subscribe-modal");
  modal.classList.add("show");
  window.onclick = function (event) {
    if (event.target.classList.contains("modal")) {
      closeSubscribeForm();
    }
  };
}

const gold_price = "price_1PtThgL7saX4DlUzele6xOaA";
const diamond_price = "price_1PuEM6L7saX4DlUzJgkTzYGy";

async function selectPlan(type, email) {
  console.log("type: ", type);
  console.log("email: ", email);
  let customer;
  await fetch("/customer-by-email", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ email }),
  })
    .then((response) => {
      if (response.ok) {
        return response.json(); // or response.json() if expecting JSON response
      }
      throw new Error("Network response was not ok.");
    })
    .then((data) => {
      customer = data;
    })
    .catch((error) => {
      console.error("There was a problem with your fetch operation:", error);
    });

  console.log("customer: ", customer);
  const priceID = type === "gold" ? gold_price : diamond_price;
  console.log("priceID: ", priceID);

  const data = {
    customer_id: customer.id,
    price_id: priceID,
  };

  fetch("/create-checkout-session", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  })
    .then((response) => {
      if (response.ok) {
        return response.json(); // or response.json() if expecting JSON response
      }
      throw new Error("Network response was not ok.");
    })
    .then((data) => {
      window.location.href = data.url;
      console.log("Checkout session created:", data);
      console.log("data: ", data.url);
    })
    .catch((error) => {
      console.error("There was a problem with your fetch operation:", error);
    });
}

function deleteAlert(id) {
  console.log("deleteAlert", id);

  fetch(`/api/delete-alert?id=${id}`, {
    method: "DELETE",
  })
    .then((response) => {
      if (response.ok) {
        console.log("Alert deleted successfully");
        // Optionally, remove the alert from the DOM or refresh the alerts list
        document.getElementById(`alert-${id}`).remove();
      } else {
        console.error("Failed to delete alert");
      }
    })
    .catch((error) => {
      console.error("Error:", error);
    });
}
