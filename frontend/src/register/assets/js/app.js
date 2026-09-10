// Because Angular does not play nice with the css on the register page
// We do this custom VUE app. Some day we should fix this and redo the CSS/HTML
var app = new Vue({
  el: "#app",

  // Data useed in this component
  data: {
    isLoading: false,
    turnstileToken: "",
    turnstileWidget: null,
    securityError: "",
    first: "",
    last: "",
    email: "",
    password: "",
    passwordConfirmed: "",
    company: "",
    errorMsg: "",
    token: "",
  },

  // Method used in this component
  methods: {
    // Load the public runtime key, then load Cloudflare only after Vue has mounted.
    loadTurnstile: function () {
      const vm = this;
      // Fetch the site key from the same backend that verifies registration requests.
      axios.get(this.getBaseUrl() + "/registration-config").then(function (response) {
        // Install the callback before loading the asynchronous Turnstile script.
        window.skyclerkTurnstileReady = function () {
          // Render explicitly so Vue never replaces an implicitly rendered widget.
          vm.turnstileWidget = window.turnstile.render("#register-turnstile", {
            sitekey: response.data.site_key,
            action: "register",
            size: "flexible",
            // Retain a successful token only until expiry or a submission attempt.
            callback: function (token) {
              vm.turnstileToken = token;
              vm.securityError = "";
            },
            // Expired challenges must be solved again before another submission.
            "expired-callback": function () {
              vm.resetTurnstile();
            },
            // Keep submission blocked while Cloudflare retries a failed challenge.
            "error-callback": function () {
              vm.turnstileToken = "";
              vm.securityError = "The security check could not load. Please retry or refresh the page.";
            }
          });
        };
        // Load the official script directly; it must not be bundled or proxied.
        const script = document.createElement("script");
        script.src = "https://challenges.cloudflare.com/turnstile/v0/api.js?onload=skyclerkTurnstileReady&render=explicit";
        script.async = true;
        // Explain script blocking or network failures without permitting signup.
        script.onerror = function () {
          vm.securityError = "The security check could not load. Please refresh the page to try again.";
        };
        // Begin downloading only after the callback and widget container exist.
        document.head.appendChild(script);
      }).catch(function () {
        vm.securityError = "Registration is temporarily unavailable. Please refresh the page to try again.";
      });
    },

    // Discard consumed or expired tokens and request a fresh challenge for retries.
    resetTurnstile: function () {
      this.turnstileToken = "";
      if (window.turnstile && this.turnstileWidget !== null) {
        // Reset this specific widget after every unsuccessful registration attempt.
        window.turnstile.reset(this.turnstileWidget);
      }
    },

    // Return the base URL
    getBaseUrl: function () {
      if (location.origin.indexOf("localhost") >= 0) {
        return "http://localhost:9090";
      }

      return "https://app.skyclerk.com";
    },

    // Return the client_id
    getClientId: function () {
      if (location.origin.indexOf("localhost") >= 0) {
        return "abc123";
      }

      return "p9KgPZ50YGrcTOb";
    },

    // Submit register form.
    submit: function () {
      const vm = this;

      // Prevent duplicate requests and direct form submission without a challenge.
      if (vm.isLoading) return;
      if (!vm.turnstileToken) {
        vm.securityError = "Please complete the security check before signing up.";
        return;
      }

      // Verify passwords match
      if (this.password != this.passwordConfirmed) {
        vm.errorMsg = "Your passwords did not match each other.";
        return;
      }

      // Setup post
      let post = {
        client_id: this.getClientId(),
        password: this.password,
        email: this.email,
        first: this.first,
        last: this.last,
        company: this.company,
        token: this.token,
        turnstile_token: this.turnstileToken,
      };

      // Clear error
      vm.errorMsg = "";

      // Start loader
      vm.isLoading = true;

      // Ajax request to  register user
      axios
        .post(this.getBaseUrl() + "/register", post)
        .then(function (response) {
          // Store access token in local storage.
          localStorage.setItem("user_id", response.data.user_id.toString());
          localStorage.setItem("user_email", vm.email);
          localStorage.setItem("access_token", response.data.access_token);
          localStorage.setItem(
            "account_id",
            response.data.account_id.toString()
          );

          // Mix panel track
          mixpanel.people.set({
            $first_name: vm.first,
            $last_name: vm.last,
            $email: vm.email,
          });
          mixpanel.identify(response.data.user_id);
          setTimeout(function () {
            mixpanel.track("register", {
              app: "web",
              accountId: response.data.account_id,
            });
          }, 1000);

          // Log events.
          if (window._paq) {
            // Track registration only when the optional analytics client is loaded.
            window._paq.push(["trackGoal", 2]);
            window._paq.push(["trackEvent", "Auth", "Register"]);
          }

          if ("ga" in window) {
            tracker = ga.getAll()[0];
            if (tracker) {
              tracker.send("event", "Auth", "Register");
            }
          }

          // Redirect to app give time for tracking to happen
          setTimeout(function () {
            window.location.href = "/";
          }, 2000);
        })
        .catch(function (error) {
          // Tokens are single-use even when later registration validation fails.
          vm.resetTurnstile();
          setTimeout(function () {
            // End loader.
            vm.isLoading = false;

            window.scrollTo(0, 0);

            if (!error.response || error.response.status >= 500 || error.response.status < 400) {
              alert(
                "An issue with our server happened. Please try again. If you have further issues please contact help@skyclerk.com."
              );
              return;
            }

            // Set error message
            vm.errorMsg = error.response.data.error;
          }, 2000);
        });
    },
  },

  // Start the security check after the widget container exists in the DOM.
  mounted: function () {
    // Load the runtime key and challenge script for ordinary and invited signups.
    this.loadTurnstile();
  },

  // Called on start up
  created() {
    let uri = window.location.href.split("?");

    if (uri.length == 2) {
      let vars = uri[1].split("&");
      let getVars = {};
      let tmp = "";
      vars.forEach(function (v) {
        tmp = v.split("=");
        if (tmp.length == 2) {
          getVars[tmp[0]] = tmp[1];
        }
      });

      // Set email if passed in.
      if (getVars.email) {
        this.email = getVars.email;
      }

      // Set first if passed in.
      if (getVars.first) {
        this.first = getVars.first;
      }

      // Set last if passed in.
      if (getVars.last) {
        this.last = getVars.last;
      }

      // Set token if passed in.
      if (getVars.token) {
        this.token = getVars.token;
      }
    }
  },
});
