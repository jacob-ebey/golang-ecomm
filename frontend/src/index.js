import React from "react";
import ReactDOM from "react-dom";
import { Route } from "react-router-dom";
import { QueryParamProvider } from "use-query-params";

import { ApolloProvider } from "@apollo/react-hooks";
import { ApolloClient, gql, ApolloLink } from "apollo-boost";
import { BrowserRouter as Router } from "react-router-dom";
import decode from "jwt-decode";

import { InMemoryCache } from "apollo-cache-inmemory";
// import { ApolloClient } from "apollo-client";
import { createUploadLink } from "apollo-upload-client";
import { setContext } from "apollo-link-context";
import fetch from "isomorphic-fetch";

import "open-iconic/font/css/open-iconic-bootstrap.scss";
import "./styles/main.scss";

import { resolvers, typeDefs } from "./resolvers";
import App from "./containers/App";
import ScrollToTop from './components/ScrollToTop';

const refreshClient = new ApolloClient({
  link: ApolloLink.from([
    setContext(async (operation, previousContext) => {
      const {
        data: { auth }
      } = await client.query({
        query: gql`
          query RefreshClientAuth {
            auth @client {
              refreshToken
            }
          }
        `
      });

      const Authorization =
        auth && auth.refreshToken ? `Bearer ${auth.refreshToken}` : undefined;

      return {
        ...previousContext,
        headers: {
          ...previousContext.headers,
          Authorization
        }
      };
    }),
    createUploadLink({
      uri: process.env.REACT_APP_GRAPHQL_ENDPOINT || "/graphql",
      fetch
    })
  ]),
  cache: new InMemoryCache()
});

const LOGOUT = gql`
  mutation ClientLogout {
    logout @client
  }
`;

const client = (window.__APOLLO_CLIENT__ = new ApolloClient({
  resolvers,
  typeDefs,
  link: ApolloLink.from([
    setContext(async (operation, previousContext) => {
      // Get access token from local resolver / cache
      const {
        data: { auth }
      } = await client.query({
        query: gql`
          query ClientAuth {
            auth @client {
              accessToken
              refreshToken
            }
          }
        `
      });

      let token = auth && auth.accessToken;

      if (token) {
        const decoded = decode(token);

        if (decoded && Date.now() >= decoded.exp * 1000) {
          const decodedRefreshToken = decode(auth.refreshToken);

          if (
            decodedRefreshToken &&
            Date.now() < decodedRefreshToken.exp * 1000
          ) {
            const { data } = await refreshClient.mutate({
              mutation: gql`
                mutation RefreshToken {
                  refreshToken {
                    accessToken: token
                    refreshToken
                  }
                }
              `
            });

            if (data && data.refreshToken) {
              client.mutate({
                mutation: gql`
                  mutation ClientSetAuth(
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
                `,
                variables: {
                  accessToken: data.refreshToken.accessToken,
                  refreshToken: data.refreshToken.refreshToken,
                  stayLoggedIn:
                    typeof localStorage.getItem("auth") !== "undefined"
                }
              });
              token = data.refreshToken.accessToken;
            } else {
              client.mutate({ mutation: LOGOUT });
              token = undefined;
            }
          } else {
            client.mutate({ mutation: LOGOUT });
            token = undefined;
          }
        }
      }

      const Authorization = token ? `Bearer ${token}` : undefined;

      return {
        ...previousContext,
        headers: {
          ...previousContext.headers,
          Authorization
        }
      };
    }),
    createUploadLink({
      uri: process.env.REACT_APP_GRAPHQL_ENDPOINT || "/graphql",
      fetch
    })
  ]),
  cache: new InMemoryCache()
}));

//Apollo Client
ReactDOM.render(
  <ApolloProvider client={client}>
    <Router>
      <QueryParamProvider ReactRouterRoute={Route}>
        <ScrollToTop />
        <App />
      </QueryParamProvider>
    </Router>
  </ApolloProvider>,
  document.getElementById("root")
);
