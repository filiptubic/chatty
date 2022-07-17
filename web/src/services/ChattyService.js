import axios from "axios";
import UserService from "./UserService";


const chattyClient = axios.create({
  baseURL: 'http://localhost:1234',
  timeout: 1000,
});

var wsConn = null

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

const createChat = (userId) => {
  return chattyClient.post("/v1/create-chat/" + userId)
}

const joinChat = (chatId) => {
  if (wsConn == null) {
    wsConn = new WebSocket('ws://localhost:1234/ws/' + chatId);
  }
  return wsConn
}

const ChattyService = {
  configure,
  joinChat,
  createChat
};

export default ChattyService;