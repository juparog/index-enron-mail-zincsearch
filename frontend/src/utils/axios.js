import axios from "axios";

const instance = axios.create();

instance.defaults.baseURL = "http://localhost:4080"; // "https://reqres.in/api",

instance.interceptors.request.use(
  (config) => {
    config.headers.Authorization = `Basic YWRtaW46YWRtaW4=`;
    return config;
  },
  (error) => Promise.reject(error)
);

export default instance;
