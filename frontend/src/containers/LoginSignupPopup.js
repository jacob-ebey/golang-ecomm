import React from "react";
import { gql } from "apollo-boost";
import { useApolloClient, useMutation } from "@apollo/react-hooks";

import Modal from "react-bootstrap/Modal";
import Tab from "react-bootstrap/Tab";
import Tabs from "react-bootstrap/Tabs";

import Error from "../components/Error";
import Login from "../components/Login";
import Signup from "../components/Signup";

const RESPONSE = gql`
  fragment LoginSignupPopupResponse on AuthResponse {
    accessToken: token
    refreshToken
  }
`;

const LOGIN = gql`
  mutation LoginSignupPopupLogin($email: String!, $password: String!) {
    login: signIn(email: $email, password: $password) {
      ...LoginSignupPopupResponse
    }
  }
  ${RESPONSE}
`;

const SET_AUTH = gql`
  mutation LoginSignupPopupSetAuth(
    $accessToken: String
    $refreshToken: String
    $stayLoggedIn: Boolean!
  ) {
    setAuth(
      accessToken: $accessToken
      refreshToken: $refreshToken
      stayLoggedIn: $stayLoggedIn
    ) @client
  }
`;

const SIGNUP = gql`
  mutation LoginSignupPopupSignup(
    $email: String!
    $password: String!
    $confirmPassword: String!
  ) {
    signup: signUp(
      email: $email
      password: $password
      confirmPassword: $confirmPassword
    ) {
      ...LoginSignupPopupResponse
    }
  }
  ${RESPONSE}
`;

function LoginSignupPopup({ show, onHide }) {
  const client = useApolloClient();

  const [login, { error: loginError, loading: loggingIn }] = useMutation(LOGIN);
  const [signup, { error: signupError, loading: signingUp }] = useMutation(
    SIGNUP
  );
  const [
    setAuthMutation,
    { error: setAuthError, loading: settingAuth }
  ] = useMutation(SET_AUTH);

  const setAuth = React.useCallback(
    async (creds, stayLoggedIn) => {
      await setAuthMutation({
        variables: {
          accessToken: creds.accessToken,
          refreshToken: creds.refreshToken,
          stayLoggedIn: stayLoggedIn
        }
      });

      client.resetStore();

      onHide && onHide();
    },
    [setAuthMutation, onHide, client]
  );

  const onLogin = React.useCallback(async creds => {
    if (loggingIn || settingAuth) {
      return;
    }

    const result = await login({
      variables: {
        email: creds.email,
        password: creds.password
      }
    });

    if (result.data && result.data.login) {
      setAuth(result.data.login, creds.stayLoggedIn);
    }
  });

  const onSignup = React.useCallback(async creds => {
    if (signingUp || settingAuth) {
      return;
    }

    const result = await signup({
      variables: {
        email: creds.email,
        password: creds.password,
        confirmPassword: creds.confirmPassword
      }
    });

    if (result.data && result.data.signup) {
      setAuth(result.data.signup, creds.stayLoggedIn);
    }
  });

  return (
    <Modal show={show} onHide={onHide} animation={true} size="sm">
      <Modal.Header closeButton>
        <Modal.Title>Login / Signup</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Tabs defaultActiveKey="login">
          <Tab eventKey="login" title="Login" disabled={loggingIn}>
            <Error error={loginError} />
            <Error error={setAuthError} />
            <Login disabled={loggingIn || settingAuth} onLogin={onLogin} />
          </Tab>
          <Tab eventKey="signup" title="Signup">
            <Error error={signupError} />
            <Error error={setAuthError} />
            <Signup disabled={signingUp || settingAuth} onSignup={onSignup} />
          </Tab>
        </Tabs>
      </Modal.Body>
    </Modal>
  );
}

export default LoginSignupPopup;
