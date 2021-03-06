Authentication using Firebase

Firebase Authentication provides a mechanism for authenticating users to the application. 

Two authentication mechanism were evaluated. Authentication using Firebase and Auth0.

Auth0 offers simple authentication as a service. It provides a login screen via a lock widget
and a customizable UI. On logging into the user's application, the user is redirected
to the login screen provided by Auth0 to perform the authentication and login.
To avoid the redirection, Auth0 was eliminated.

Firebase authentication has been implemented for login. During login, the /user
backend Go API is invoked with the logic credentials. A custom token is minted
using the userId of the user fetched from the User database. This is sent
to the React application which performs the sign in using the custom token generated by Firebase.

The routes are then authenticated using by verifying the token sent along with the request header.
 

