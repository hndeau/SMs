AWS.config.region = 'us-east-1'; // Replace with the region you are using
AWS.config.credentials = new AWS.CognitoIdentityCredentials({
    IdentityPoolId: 'us-east-1:175761ac-5e82-41e0-91ab-0d952714f634', // Replace with your Identity Pool ID
});

const poolData = {
    UserPoolId: 'us-east-1_WusNRaP2h', // Replace with your User Pool ID
    ClientId: '34mgjfocrlfp3c4ij35qoe8d4b', // Replace with your App Client ID
};

let cognitoUser;

function signIn() {
    var authenticationData = {
        Username: document.getElementById('username').value,
        Password: document.getElementById('password').value,
    };
    var authenticationDetails = new AmazonCognitoIdentity.AuthenticationDetails(
        authenticationData
    );
    var userPool = new AmazonCognitoIdentity.CognitoUserPool(poolData);
    var userData = {
        Username: authenticationData.Username,
        Pool: userPool,
    };
    cognitoUser = new AmazonCognitoIdentity.CognitoUser(userData);
    cognitoUser.authenticateUser(authenticationDetails, {
        onSuccess: function (result) {
            // Save tokens as cookies
            document.cookie = "access_token=" + result.getAccessToken().getJwtToken();
            document.cookie = "id_token=" + result.getIdToken().getJwtToken();
            document.cookie = "refresh_token=" + result.getRefreshToken().getToken();
            // Redirect to /chat
            window.location.href = "/chat";
        },
        onFailure: function (err) {
            console.log(err);
        },
    });
}

function signOut() {
    var userPool = new AmazonCognitoIdentity.CognitoUserPool(poolData);
    var userData = {
        Username: authenticationData.Username,
        Pool: userPool,
    };
    cognitoUser = new AmazonCognitoIdentity.CognitoUser(userData);
    // Get the current user from local storage, if available
    const currentUser = cognitoIdentityServiceProvider.getCurrentUser();

    if (currentUser !== null) {
        // If the current user is loaded in memory, sign them out
        currentUser.signOut(() => {
            console.log('User signed out successfully');
        });
    } else {
        console.log('No user is currently loaded in memory');
    }
}

function signUp() {
    var userPoolData = {
        UserPoolId: 'us-east-1_WusNRaP2h',
        ClientId: '34mgjfocrlfp3c4ij35qoe8d4b',
    };
    var userPool = new AmazonCognitoIdentity.CognitoUserPool(userPoolData);

    var username = document.getElementById('signup-username').value;
    var password = document.getElementById('signup-password').value;
    var email = document.getElementById('signup-email').value;

    var attributeList = [
        new AmazonCognitoIdentity.CognitoUserAttribute({
            Name: 'email',
            Value: email,
        }),
    ];

    userPool.signUp(username, password, attributeList, null, (err, result) => {
        if (err) {
            console.log(err);
            return;
        }
        cognitoUser = result.user;
        console.log('User signed up:', cognitoUser.getUsername());
    });
}

function getCurrentUser() {
    // Get the access token from cookies
    const accessToken = getCookie('access_token');

    // Set up the parameters for the getUser API call
    const params = {
        AccessToken: accessToken
    };

    // Create a new CognitoIdentityServiceProvider object
    const cognitoIdentityServiceProvider = new AWS.CognitoIdentityServiceProvider();

    // Call the getUser API with the session token
    cognitoIdentityServiceProvider.getUser(params, (err, result) => {
        if (err) {
            console.log(err);
            return;
        }

        // Get a reference to the output element
        const outputElement = document.getElementById("username");

        // Set the output element's text to the username
        outputElement.textContent = result.Username;
    });
}

async function sendMessage(conversationId, message) {
    const apiEndpoint = 'https://58z24w81cl.execute-api.us-east-1.amazonaws.com/test/sendMessage/';
    const payload = {
        conversation_id: conversationId,
        message: message
    };

    try {
        const response = await fetch(apiEndpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'id_token': getCookie("id_token")
            },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const jsonResponse = await response.json();
        return jsonResponse;
    } catch (error) {
        console.error('Error calling API Gateway endpoint:', error);
        throw error;
    }
}



// Helper function
function getCookie(name) {
    const value = "; " + document.cookie;
    const parts = value.split("; " + name + "=");
    if (parts.length === 2) {
        return parts.pop().split(";").shift();
    }
}