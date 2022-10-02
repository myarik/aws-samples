import {Button, Col, notification, Row} from "antd";
import {HomeOutlined, PoweroffOutlined} from "@ant-design/icons";
import {useNavigate} from "react-router-dom";
import {useAuth} from "../../context/auth";

export default function Dashboard() {
    let {signOut} = useAuth();
    let navigate = useNavigate();

    const handleSignOut = async () => {
        const result = await signOut();
        if (result.success) {
            navigate("/");
        } else {
            notification.error({
                message: "Error",
                description: result.message,
            });
        }
    };

    return (
        <>
            <Row>
                <Col flex={1} offset={22}>
                    <Button
                        icon={<HomeOutlined/>}
                        onClick={() => navigate("/")}
                    />
                </Col>
                <Col flex={18}>
                    <Button
                        icon={<PoweroffOutlined/>}
                        onClick={() => handleSignOut()}
                    />
                </Col>
            </Row>
            <main>
                <h2>Welcome to the dashboard!</h2>
                <p>You can do this, I believe in you.</p>
            </main>
        </>
    );
}