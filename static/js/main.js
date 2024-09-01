function showModal() {
    const modal = document.getElementById('signup-modal');
    modal.classList.add('show');
    window.onclick = function(event) {
        if (event.target == modal) {
          closeSignUpForm();
        }
      };
  }

  function closeSignUpForm() {
    const modal = document.getElementById('signup-modal');
    modal.classList.remove('show');
  }

  function showSignUpForm() {
    const signupForm = document.getElementById('signup-form');
    const loginForm = document.getElementById('login-form');
    signupForm.classList.add('active-form');
    loginForm.classList.remove('active-form');
  }

  function showLoginForm() {
    const signupForm = document.getElementById('signup-form');
    const loginForm = document.getElementById('login-form');
    loginForm.classList.add('active-form');
    signupForm.classList.remove('active-form');
  }