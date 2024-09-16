<h2><strong>Trading Alerts</strong></h2>
<p></p>
<h3>Requirements</h3>
<ol>
  <li><p>✅ landing page</p></li>
  <li><p>✅ pricing page</p></li>
  <li><p>✅ payment mechanism/gateway</p></li>
  <li><p>✅ alerts dashboard</p></li>
  <li><p>✅ price change alerts</p></li>
  <li><p>✅ login/signup</p></li>
  <li><p>✅ load test</p></li>
</ol>
<p></p>
<h4>Alerts Dashboard</h4>
<ol>
  <li><p>✅ create alert</p></li>
  <li><p>✅ delete alert</p></li>
  <li><p>disable all alerts -- if you go on a holiday, etc.</p></li>
  <li><p>remove all alerts</p></li>
  <li><p>✅ show all alerts, active, disabled, completed</p></li>
</ol>
<p></p>
<h4>Price Change Alerts</h4>
<p>
  On a change of a particular stock price moved to a specified value then
  trigger an alert. For now just email alert. The alert check runs every second
  and checks for alerts that are active and if a certain stock changed its value
  to trigger the alert then send the email, make the alert status completed.
</p>
<p></p>
<h4>Pricing Page</h4>
<p><strong>Silver tier - Free</strong></p>
<ul>
  <li><p>5 alerts</p></li>
  <li><p>any stock</p></li>
  <li><p>email alert</p></li>
</ul>
<p></p>
<p><strong>Gold tier - $4.99</strong></p>
<ul>
  <li><p>100 alerts</p></li>
  <li><p>any stock</p></li>
  <li><p>email alert</p></li>
</ul>
<p></p>
<p><strong>Diamond tier - $9.99</strong></p>
<ul>
  <li><p>1000 alerts</p></li>
  <li><p>any stock</p></li>
  <li>
    <p>email alert</p>
    <p></p>
  </li>
</ul>
<p></p>
<h4>Login/SignUp</h4>
<p>
  Just a basic login and sign up that we already have implemented. Should send a
  token on login and then be able to call user specified endpoints.
</p>
<p></p>
<h4>Load test</h4>
<p>
  Make lots of users with hundreds of alerts and test if the app can handle it.
  If it slows down or no, if there need to be changes made, optimisations, etc.
</p>
<p></p>
<p><strong>Last touches:</strong></p>
<ol>
  <li><p>✅ send alerts based on the users email</p></li>
  <li><p>✅ add subscription success and cancel pages</p></li>
  <li><p>✅ add reset password</p></li>
  <li><p>✅ add created_at and completed_at to alerts</p></li>
  <li><p>✅ add cancel subscription</p></li>
  <li>
    <p>
      ✅ make the setup for users subscriptions not be from hardcoded users but
      from db
    </p>
  </li>
  <li>
    <p>
      ✅ On signup now createCustomer with the email in stripe and check his
      subscription status. so add him to UserSubscription. Also on signUp
      redirect to ?login=true and when that happens show the login modal form.
    </p>
  </li>
</ol>
<p><strong>More last touches:</strong></p>

<ol>
  <li>
    <p>
      ✅ <strong>REFACTOR a lot</strong> and see where things can be moved to
      and out of.
    </p>
  </li>
  <li>
    <p>✅ And then <strong>lots of TESTING.</strong></p>
  </li>
  <li>
    <p>
      Fix the way we abuse <strong>subscriptionUtils.Setup</strong>, where after
      every time someone subs we call the whole damn thing and refetch the
      subscription type for all users
    </p>
  </li>
  <li><p>✅ Add documentation</p></li>
  <li>
    <p>release <strong>Trading Alerts</strong> on the pi</p>
  </li>
  <li><p>finish documentation</p></li>
  <li><p>✅ Test everything, from registering, to subscribing and creating alerts.</p></li>
  <li><p>✅ Then promote your app on reddit and send out emails.</p></li>
  <li><p>✅ Try to make PWA and send notifications to the phone through it.</p></li>
  <li><p>✅ You can reset password for an email that's not registered. And you can even click on the link from emial and change password and no error is issued, actually says it reset the password. But when you try to login it says ofc the credentials are invalid.</p></li>
  <li><p>✅ Fix the invoice.success webhook</p></li>
  <li><p>✅ For some reason when you get redirected from the stripe payment wall then you don't see yourself as logged in anymore and have to refresh</p></li>
  <li><p>✅ Finish PWA and remove the notifications as it doesnt work on iphone and mac either. Its a shame really. Look into other alternatives for sending notifications.</p></li>


</ol>

<p></p>
