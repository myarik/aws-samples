import {useEffect, useState} from "react";

import {Button, Col, notification, Row, Space, Spin, Table, Tag} from "antd";
import {HomeOutlined, PoweroffOutlined} from "@ant-design/icons";
import ColumnGroup from "antd/es/table/ColumnGroup";
import Column from "antd/es/table/Column";
import {axiosPrivate} from "../../common/axiosPrivate";
import {useAuth} from "../../context/auth";
import {useNavigate} from "react-router-dom";

export default function Dashboard() {
    let {signOut} = useAuth();
    let navigate = useNavigate();

    const [data, setData] = useState([]);
    const [isLoading, setIsLoading] = useState(false);

    useEffect(() => {
        const fetchData = async () => {
            setIsLoading(true);

            try {
                const result = await axiosPrivate("/v1/users");
                console.log("result", result);
                setData(result.data);
            } catch (error) {
                console.log(error);
            }
            setIsLoading(false);
        };
        // Get data from API
        fetchData();
    }, []);


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
                {
                    isLoading ? <Spin size="middle"/>
                        : (
                            <Table dataSource={data}>
                                <ColumnGroup title="Name">
                                    <Column title="First Name" dataIndex="first_name" key="firstName"/>
                                    <Column title="Last Name" dataIndex="last_name" key="lastName"/>
                                </ColumnGroup>
                                <Column title="Age" dataIndex="age" key="age"/>
                                <Column title="Address" dataIndex="address" key="address"/>
                                <Column
                                    title="Tags"
                                    dataIndex="tags"
                                    key="tags"
                                    render={(tags) => (
                                        <>
                                            {tags.map((tag) => (
                                                <Tag color="blue" key={tag}>
                                                    {tag}
                                                </Tag>
                                            ))}
                                        </>
                                    )}
                                />
                                <Column
                                    title="Action"
                                    key="action"
                                    render={(_, record) => (
                                        <Space size="middle">
                                            <a>Delete</a>
                                        </Space>
                                    )}
                                />
                            </Table>)}
            </main>
        </>
    );
}