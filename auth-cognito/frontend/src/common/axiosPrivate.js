import axios from "axios";
import {AUTH_USER_TOKEN_KEY} from "../utils/constants";

axios.defaults.baseURL = process.env.REACT_APP_API_URL;

axios.interceptors.request.use(
    async (config) => {
        const user = JSON.parse(localStorage.getItem(AUTH_USER_TOKEN_KEY));

        if (user?.token) {
            config.headers = {
                ...config.headers,
                authorization: `Bearer ${user?.token}`,
            };
        }

        return config;
    },
    (error) => Promise.reject(error)
);

axios.interceptors.response.use(
    (response) => response,
    async (error) => {
        console.log(error);
        return Promise.reject(error);
    }
);

export const axiosPrivate = axios;