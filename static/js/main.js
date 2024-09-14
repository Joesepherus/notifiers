function showModal() {
    const modal = document.getElementById("signup-modal");
    modal.classList.add("show");
    window.onclick = function(event) {
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
    const resetPasswordForm = document.getElementById("reset-password-form");
    const setPasswordForm = document.getElementById("set-password-form");
    signupForm.classList.add("active-form");
    loginForm.classList.remove("active-form");
    resetPasswordForm.classList.remove("active-form");
    setPasswordForm.classList.remove("active-form");
}

function showLoginForm() {
    const signupForm = document.getElementById("signup-form");
    const loginForm = document.getElementById("login-form");
    const resetPasswordForm = document.getElementById("reset-password-form");
    const setPasswordForm = document.getElementById("set-password-form");
    loginForm.classList.add("active-form");
    signupForm.classList.remove("active-form");
    resetPasswordForm.classList.remove("active-form");
    setPasswordForm.classList.remove("active-form");
}

function showResetPasswordForm() {
    const signupForm = document.getElementById("signup-form");
    const loginForm = document.getElementById("login-form");
    const resetPasswordForm = document.getElementById("reset-password-form");
    const setPasswordForm = document.getElementById("set-password-form");

    loginForm.classList.remove("active-form");
    signupForm.classList.remove("active-form");
    resetPasswordForm.classList.add("active-form");
    setPasswordForm.classList.remove("active-form");
}

function showSetPasswordForm() {
    const signupForm = document.getElementById("signup-form");
    const loginForm = document.getElementById("login-form");
    const resetPasswordForm = document.getElementById("reset-password-form");
    const setPasswordForm = document.getElementById("set-password-form");

    loginForm.classList.remove("active-form");
    signupForm.classList.remove("active-form");
    resetPasswordForm.classList.remove("active-form");
    setPasswordForm.classList.add("active-form");
}

function openModalShowLoginForm() {
    showModal();
    showLoginForm();
}

function openModalShowSignUpForm() {
    showModal();
    showSignUpForm();
}

function openModalShowSetPasswordForm() {
    showModal();
    showSetPasswordForm();
}

function closeSubscribeForm() {
    const modal = document.getElementById("subscribe-modal");
    modal.classList.remove("show");
}

function showSubscribeModal() {
    const modal = document.getElementById("subscribe-modal");
    modal.classList.add("show");
    window.onclick = function(event) {
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
    await fetch("/api/customer-by-email", {
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

    fetch("/api/create-checkout-session", {
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

function showAlerts() {
    document.getElementById("alertsTable").classList.remove("none");
    document.getElementById("completedAlertsTable").classList.add("none");
    setActiveTab("Active Alerts");
}

function showCompletedAlerts() {
    document.getElementById("alertsTable").classList.add("none");
    document.getElementById("completedAlertsTable").classList.remove("none");
    setActiveTab("Completed Alerts");
}

function setActiveTab(tabName) {
    const tabs = document.querySelectorAll(".tab");
    tabs.forEach((tab) => {
        if (tab.innerText === tabName) {
            tab.classList.add("active");
        } else {
            tab.classList.remove("active");
        }
    });
}

window.onload = function() {
    const urlParams = new URLSearchParams(window.location.search);

    // Get the value of the 'token' parameter
    const token = urlParams.get("token");
    console.log("KOKOT")

    // Check if the token exists
    if (token) {
        console.log("Token:", token);
        document.getElementById("tokenInput").value = token;
        openModalShowSetPasswordForm();
        // Perform actions with the token, like sending it to the server
    } else {
        console.log("No token found in the URL");
    }

    const login = urlParams.get("login");
    if (login) {
        openModalShowLoginForm()
    }

};

window.onload = function() {
    if ('serviceWorker' in navigator) {
        navigator.serviceWorker.register('/sw.js').then(function(registration) {
            console.log('Service Worker registered with scope:', registration.scope);
        }).catch(function(error) {
            console.log('Service Worker registration failed:', error);
        });
    }

    const notifyBtn = document.getElementById('notifyBtn');
    notifyBtn.addEventListener('click', () => {
        Notification.requestPermission().then(permission => {
            if (permission === 'granted') {
                subscribeUser();
            }
        });
    });

    function subscribeUser() {
        navigator.serviceWorker.ready.then(function(registration) {
            registration.pushManager.subscribe({
                userVisibleOnly: true,
                applicationServerKey: urlBase64ToUint8Array('BJKWKNBvXlmmrd3yGMlKWOGau4ijiYOp3oP4TGqYbScnFQhK_5qs4x_LPyXltvQARznsg7kz4Wvmef2DluuREao')
            }).then(function(subscription) {
                console.log('User is subscribed:', subscription);

                // Send the subscription to the server to save it
                fetch('/api/subscribe', {
                    method: 'POST',
                    body: JSON.stringify(subscription),
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });
            }).catch(function(error) {
                console.log('Failed to subscribe the user:', error);
            });
        });
    }

    // Utility function to convert the VAPID key
    function urlBase64ToUint8Array(base64String) {
        const padding = '='.repeat((4 - base64String.length % 4) % 4);
        const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/');
        const rawData = atob(base64);
        return Uint8Array.from([...rawData].map(char => char.charCodeAt(0)));
    }

}
