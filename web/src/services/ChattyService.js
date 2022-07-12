import axios from "axios";
import UserService from "./UserService";

const HttpMethods = {
  GET: 'GET',
  POST: 'POST',
  DELETE: 'DELETE',
};

// const _axios = axios.create();
const chattyClient = axios.create({
  baseURL: 'http://localhost:1234',
  timeout: 1000,
});

const configure = () => {
  chattyClient.interceptors.request.use((config) => {
    if (UserService.isLoggedIn()) {
      const cb = () => {
        config.headers.Authorization = `Bearer ${UserService.getToken()}`;
        return Promise.resolve(config);
      };
      return UserService.updateToken(cb);
    }
  });
};


const getSession = () => {
  return chattyClient.get("/v1/session")
}

const getAxiosClient = () => chattyClient;

const HttpService = {
  HttpMethods,
  configure,
  getAxiosClient,
  getSession
};

export default HttpService;