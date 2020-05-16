// Because Angular does not play nice with the css on the register page
// We do this custom VUE app. Some day we should fix this and redo the CSS/HTML
var app = new Vue({
  el: '#app',

	// Data useed in this component
  data: {
    isLoading: false,
		first: "",
		last: "",
		email: "",
		password: "",
		passwordConfirmed: "",
		company: "",
		errorMsg: "",
    token: ""
  },

	// Method used in this component
	methods: {
		// Return the base URL
		getBaseUrl: function () {
			if(location.origin.indexOf("localhost") >= 0) {
				return "http://localhost:9090";
			}

			return "https://app.skyclerk.com";
		},

    // Return the client_id
		getClientId: function () {
			if(location.origin.indexOf("localhost") >= 0) {
				return "abc123";
			}

			return "p9KgPZ50YGrcTOb";
		},

		// Submit register form.
		submit: function () {
			const vm = this;

      // Verify passwords match
      if(this.password != this.passwordConfirmed) {
        vm.errorMsg = "Your passwords did not match each other.";
        return;
      }

      // Setup post
			let post = {
			  "client_id": this.getClientId(),
			  "password": this.password,
			  "email": this.email,
			  "first": this.first,
			  "last": this.last,
        "company": this.company,
        "token": this.token
			};

      // Clear error
      vm.errorMsg = "";

      // Start loader
      vm.isLoading = true;

			// Ajax request to  register user
			axios.post(this.getBaseUrl() + '/register', post)
			  .then(function (response) {
					// Store access token in local storage.
					localStorage.setItem('user_id', response.data.user_id.toString());
					localStorage.setItem('user_email', vm.email);
					localStorage.setItem('access_token', response.data.access_token);
					localStorage.setItem('account_id', response.data.account_id.toString());

          // Mix panel track
          mixpanel.people.set({ "$first_name": vm.first, "$last_name": vm.last, "$email": vm.email });
      		mixpanel.identify(response.data.user_id);
          setTimeout(function() {
            mixpanel.track('register', { app: "web", "accountId": response.data.account_id });
          }, 1000);

          // Log events.
					_paq.push(['trackGoal', 2]);
					_paq.push(['trackEvent', 'Auth', 'Register']);

          if ("ga" in window) {
            tracker = ga.getAll()[0];
            if (tracker) {
              tracker.send("event", "Auth", "Register");
            }
          }

					// Redirect to app give time for tracking to happen
          setTimeout(function() {
            window.location.href = "/";
          }, 2000);
			  })
			  .catch(function (error) {
          setTimeout(function() {
            // End loader.
            vm.isLoading = false;

            window.scrollTo(0, 0);

  					if((error.response.status >= 500) || (error.response.status < 400)) {
  						alert("An issue with our server happened. Please try again. If you have further issues please contact help@skyclerk.com.");
  						return;
  					}

  					// Set error message
  					vm.errorMsg = error.response.data.error;
          }, 2000);
			  });
		}
	},

  // Called on start up
  created()
  {
    let uri = window.location.href.split('?');

    if (uri.length == 2)  {
      let vars = uri[1].split('&');
      let getVars = {};
      let tmp = '';
      vars.forEach(function(v){
        tmp = v.split('=');
        if(tmp.length == 2) {
          getVars[tmp[0]] = tmp[1];
        }
     });

     // Set email if passed in.
     if(getVars.email) {
       this.email = getVars.email;
     }

     // Set first if passed in.
     if(getVars.first) {
       this.first = getVars.first;
     }

     // Set last if passed in.
     if(getVars.last) {
       this.last = getVars.last;
     }

     // Set token if passed in.
     if(getVars.token) {
       this.token = getVars.token;
     }
   }
 }
})
