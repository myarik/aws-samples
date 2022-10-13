import {Navigate, Route, Routes, useLocation} from "react-router-dom";
import Home from "./containers/home";
import Dashboard from "./containers/dashboard";
import Login from "./containers/login";
import {AuthProvider, useAuth} from "./context/auth";
/** Ant design stylesheet */
import 'antd/dist/antd.css';

// A wrapper for <Route> that redirects to the login
// screen if you're not yet authenticated.
function PrivateRoute({children}) {
    let {isAuthenticated} = useAuth();
    let location = useLocation();
    if (!isAuthenticated) {
        return <Navigate to="/login" state={{from: location}} replace/>;
    }
    return children;
}

// A wrapper for <Route> that redirects to the home
// screen if you're authenticated.
function AnonymousRoute({children}) {
    let {isAuthenticated} = useAuth();
    let location = useLocation();
    if (isAuthenticated) {
        return <Navigate to="/" state={{from: location}} replace/>;
    }
    return children;
}


function App() {
    return (
        <AuthProvider>
            <Routes>
                <Route path="/" element={<Home/>}/>
                <Route path="/dashboard"
                       element={<PrivateRoute><Dashboard/></PrivateRoute>}
                />
                <Route path="login"
                       element={<AnonymousRoute><Login/></AnonymousRoute>}/>
            </Routes>
        </AuthProvider>
    );
}

export default App;
