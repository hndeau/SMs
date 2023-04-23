document.addEventListener("DOMContentLoaded", function () {
    const awsconfig = {

    };

    // Configure Amplify
    window.Amplify.configure(awsconfig);

    // Handle the response from Cognito's Hosted UI and get the username of the current user
    async function handleCallback() {
        try {
            const user = await window.Amplify.Auth.currentAuthenticatedUser();
            const username = user.username;
            console.log("Current username:", username);

            // Redirect to a desired route after successful login
            window.location.href = "/chat";
        } catch (error) {
            console.error("Error handling callback:", error);

            // Redirect to a login route if authentication fails
            window.location.href = "/login";
        }
    }

    handleCallback();
});
