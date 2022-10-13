export class User {
    constructor(cognitoUser) {
        this.username = cognitoUser.username;
        this.token = cognitoUser.signInUserSession.idToken.jwtToken;
    }
}