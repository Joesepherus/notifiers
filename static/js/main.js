function showSignUpForm() {
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