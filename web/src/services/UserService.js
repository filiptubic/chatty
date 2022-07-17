import Keycloak from "keycloak-js";
import axios from 'axios';

const keycloak = new Keycloak({
  url: 'http://localhost:8080/auth/',
  realm: 'chatty-realm',
  clientId: 'chatty-react'
});

const initKeycloak = (onAuthenticatedCallback) => {
  keycloak.init({
    onLoad: 'login-required',
  })
    .then((authenticated) => {
      if (!authenticated) {
        console.log("user is not authenticated..!");
      }
      onAuthenticatedCallback();
    })
    .catch(console.error);
};

const doLogin = keycloak.login;

const doLogout = keycloak.logout;

const getToken = () => keycloak.token;

const isLoggedIn = () => !!keycloak.token;

const updateToken = (successCallback) =>
  keycloak.updateToken(5)
    .then(successCallback)
    .catch(doLogin);

const getUsername = () => keycloak.tokenParsed.preferred_username;

const hasRole = (roles) => roles.some((role) => keycloak.hasRealmRole(role));

const getParsedToken = () => {return keycloak.tokenParsed}

const listUsers = (search) => {
  const token = getToken()
  // TODO move to chatty 
  return axios.get(
    "http://localhost:1234/v1/users",
    { 
      params: {"search": search},
      headers: {'Authorization': 'Bearer ' + token}
    }
  )
}

const UserService = {
  initKeycloak,
  doLogin,
  doLogout,
  isLoggedIn,
  getToken,
  updateToken,
  getUsername,
  hasRole,
  getParsedToken,
  listUsers,
  keycloak
};

export default UserService;