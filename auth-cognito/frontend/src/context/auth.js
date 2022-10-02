import * as React from 'react'
import {AUTH_USER_TOKEN_KEY} from "../utils/constants";

// const jwtDecode = require('jwt-decode');

const AuthContext = React.createContext();

// function validateToken(token) {
//     if (!token) {
//         return false;
//     }
//     return true;
//     // try {
//     //     let decodedJwt = jwtDecode(token);
//     //     return decodedJwt.exp >= Date.now() / 1000;
//     // } catch (e) {
//     //     return false;
//     // }
// }

function AuthProvider({children}) {
    const [isAuthenticated, setIsAuthenticated] = React.useState(false);
    // const [token, setToken] = React.useState(null);

    // React.useEffect(() => {
    //     // Check if the user is logged in when the app loads
    //     // let token = localStorage.getItem(AUTH_USER_TOKEN_KEY);
    //     // if (token) {
    //     //     setIsAuthenticated(true);
    //     // }
    //     // if (validateToken(token)) {
    //     //     setIsAuthenticated(true);
    //     //     // setToken(token);
    //     // }
    // });

    const signIn = async (username, password) => {
        try {
            // const result = await Auth.signIn(username, password);
            // setUsername(result.username);
            // setIsAuthenticated(true);
            setIsAuthenticated(true);
            // setToken('fake-token');
            return {success: true, message: ""};
        } catch (error) {
            return {
                success: false,
                message: "LOGIN FAIL",
            };
        }
    };

    const signOut = async () => {
        try {
            // await Auth.signOut();
            // setToken(null);
            setIsAuthenticated(false);
            return {success: true, message: ""};
        } catch (error) {
            return {
                success: false,
                message: "LOGOUT FAIL",
            };
        }
    };
    const value = {isAuthenticated, signIn, signOut}
    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

function useAuth() {
    const context = React.useContext(AuthContext)
    if (context === undefined) {
        throw new Error('useAuth must be used within a AuthProvider')
    }
    return context
}


export {AuthProvider, useAuth};
