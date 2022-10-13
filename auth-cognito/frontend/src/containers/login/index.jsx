import * as React from 'react';
import {useLocation, useNavigate} from "react-router-dom";
import {Button, Form, Input, notification} from 'antd';
import {useAuth} from "../../context/auth";

export default function Login() {
    let {signIn} = useAuth();

    let location = useLocation();
    let navigate = useNavigate();

    let from = location.state?.from?.pathname || "/";

    const handleSubmit = async (values) => {
        const result = await signIn(values.username, values.password);
        if (result.success) {
            navigate(from);
        } else {
            console.log(result.error);
            notification.error({
                message: 'User confirmation failed',
                description: "Please check your username and password",
                placement: 'topRight',
                duration: 2
            });
        }
    };

    return (
        <>
            <main>
                <h2>Login page</h2>
            </main>
            <div>
                <p>You must log in to view the page at {from}</p>
            </div>
            <Form
                name="basic"
                onFinish={handleSubmit}
                labelCol={{span: 2}}
                wrapperCol={{span: 6}}
                autoComplete="off"
            >
                <Form.Item
                    label="Username"
                    name="username"
                    rules={[
                        {
                            required: true,
                            message: 'Please input your username!',
                        },
                    ]}
                >
                    <Input/>
                </Form.Item>

                <Form.Item
                    label="Password"
                    name="password"
                    rules={[
                        {
                            required: true,
                            message: 'Please input your password!',
                        },
                    ]}
                >
                    <Input.Password/>
                </Form.Item>

                <Form.Item
                    wrapperCol={{
                        offset: 2,
                    }}
                >
                    <Button type="primary" htmlType="submit">
                        Submit
                    </Button>
                </Form.Item>
            </Form>
        </>
    );
}