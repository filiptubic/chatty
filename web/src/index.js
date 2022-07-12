import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import UserService from './services/UserService';
import ChattyService from './services/ChattyService';

const root = ReactDOM.createRoot(document.getElementById('root'));
const render = () => root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

UserService.initKeycloak(render)
ChattyService.configure()
