// Because Angular does not play nice with the css on the register page
// We do this custom VUE app. Some day we should fix this and redo the CSS/HTML
var app = new Vue({
  el: '#app',

	// Data useed in this component
  data: {
		first: "",
		last: "",
		email: "",
		password: "",
		passwordConfirmed: "",
		company: "",
		errorMsg: ""
  },

	// Method used in this component
	methods: {
		// Return the base URL
		getBaseUrl: function () {
			if(location.origin.indexOf("localhost")) {
				return "http://localhost:9090";
			}

			return "https://app.skyclerk.com";
		},

    // Return the client_id
		getClientId: function () {
			if(location.origin.indexOf("localhost")) {
				return "abc123";
			}

			return "p9KgPZ50YGrcTOb";
		},

		// Submit register form.
		submit: function () {
			const vm = this;

			let post = {
			  "client_id": this.getClientId(),
			  "password": this.password,
			  "email": this.email,
			  "first": this.first,
			  "last": this.last,
        "company": this.company
			};

			// Ajax request to  register user
			axios.post(this.getBaseUrl() + '/register', post)
			  .then(function (response) {
					// Store access token in local storage.
					localStorage.setItem('user_id', response.data.user_id.toString());
					localStorage.setItem('user_email', vm.email);
					localStorage.setItem('access_token', response.data.access_token);
					localStorage.setItem('account_id', response.data.account_id.toString());

					// Redirect to app
					window.location.href = "/";
			  })
			  .catch(function (error) {
					if((error.response.status >= 500) || (error.response.status < 400)) {
						alert("An issue with our server happened. Please try again. If you have further issues please contact help@skyclerk.com.");
						return;
					}

					// Set error message
					vm.errorMsg = error.response.data.error;
			  });
		}
	}
})
