import * as React from 'react'

/** Amplify config */
import {AuthConfig} from '../config/auth';
import {Amplify, Auth} from "aws-amplify";
import {AUTH_USER_TOKEN_KEY} from "../utils/constants";
import {User} from "../models/user";

/** Configure amplify */
Amplify.configure({Auth: AuthConfig});


const AuthContext = React.createContext();

function AuthProvider({children}) {
    const [isAuthenticated, setIsAuthenticated] = React.useState(false);

    React.useEffect(() => {
        // Check if the user is logged in when the app loads
        Auth.currentAuthenticatedUser()
            .then((result) => {
                setIsAuthenticated(true);
                localStorage.setItem(AUTH_USER_TOKEN_KEY, JSON.stringify(new User(result)));
            })
            .catch(() => {
                setIsAuthenticated(false);
            });
    }, []);

    const signIn = async (username, password) => {
        try {
            const result = await Auth.signIn(username, password);
            localStorage.setItem(AUTH_USER_TOKEN_KEY, JSON.stringify(new User(result)));
            setIsAuthenticated(true);
            return {success: true, message: ""};
        } catch (error) {
            console.log('Error signing in user: ', error);
            return {
                success: false,
                message: "LOGIN FAIL",
            };
        }
    };

    const signOut = async () => {
        try {
            await Auth.signOut();
            localStorage.removeItem(AUTH_USER_TOKEN_KEY);
            setIsAuthenticated(false);
            return {success: true, message: ""};
        } catch (error) {
            console.log('Error signing out user: ', error);
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
