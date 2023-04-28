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
    cognitoUser = userPool.getCurrentUser();

    if (cognitoUser !== null) {
        cognitoUser.signOut();
        console.log('User signed out successfully');
    } else {
        console.log('No user is currently loaded in memory');
    }
    // Redirect to /login after successful sign up
    window.location.href = '/';
}

function signUp() {
    var userPoolData = {
        UserPoolId: 'us-east-1_WusNRaP2h',
        ClientId: '34mgjfocrlfp3c4ij35qoe8d4b',
    };
    var userPool = new AmazonCognitoIdentity.CognitoUserPool(userPoolData);

    var username = document.getElementById('signup-username').value;
    var password = document.getElementById('signup-password').value;
    var confirmPassword = document.getElementById('signup-confirm-password').value;
    var email = document.getElementById('signup-email').value;

    if (password !== confirmPassword) {
        document.getElementById('error-response').innerText = 'Passwords do not match.';
        return;
    }

    var attributeList = [
        new AmazonCognitoIdentity.CognitoUserAttribute({
            Name: 'email',
            Value: email,
        }),
    ];

    userPool.signUp(username, password, attributeList, null, (err, result) => {
        if (err) {
            document.getElementById('error-response').innerText = err.message || JSON.stringify(err);
            return;
        }
        cognitoUser = result.user;
        console.log('User signed up:', cognitoUser.getUsername());

        // Redirect to /login after successful sign up
        window.location.href = '/login';
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
    const proxyEndpoint = '/send'; // Replace with the actual URL where your Go proxy is running
    const payload = {
        conversation_id: conversationId,
        message: message
    };

    let result;
    try {
        const response = await fetch(proxyEndpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'id_token': getCookie("id_token")
            },
            body: JSON.stringify(payload)
        });
        result = await response;
        console.log(result)
        return result;
    } catch (error) {
        console.error('Error calling Go intermediary:', error);
        throw error;
    }
}

async function getMessages(conversationId, timestamp) {
    const proxyEndpoint = '/retrieve'; // Replace with the actual URL where your Go proxy is running
    const payload = {
        conversation_id: conversationId,
        timestamp: timestamp
    };

    try {
        const response = await fetch(proxyEndpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'id_token': getCookie("id_token")
            },
            body: JSON.stringify(payload)
        });
        result = await response;
        console.log(result)
        return result;
    } catch (error) {
        console.error('Error calling Go intermediary:', error);
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