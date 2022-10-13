import {Link} from "react-router-dom";

export default function Home() {
    return (
        <>
            <main>
                <h2>Home page</h2>
            </main>
            <nav>
                <Link to="/dashboard">Dashboard</Link>
            </nav>
        </>
    );
}