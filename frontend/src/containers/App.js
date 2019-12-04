import React from "react";
import { Switch, Route } from "react-router-dom";

import { gql } from "apollo-boost";
import { useApolloClient, useMutation, useQuery } from "@apollo/react-hooks";
import decode from "jwt-decode";

import Error from "../components/Error";
import Header from "../components/Header";

import AdminProductPage from "../routes/admin/Product";
import AdminProductsPage from "../routes/admin/Products";
import AdminTransactionPage from "../routes/admin/Transaction";
import AdminTransactionsPage from "../routes/admin/Transactions";
import AdminVariantPage from "../routes/admin/Variant";
import CheckoutPage from "../routes/Checkout";
import HomePage from "../routes/Home";
import ProductPage from "../routes/Product";
import ProfilePage from "../routes/profile/Profile";
import ShopPage from "../routes/Shop";

import CartPopup from "./CartPopup";
import LoginSignupPopup from "./LoginSignupPopup";

const APP_QUERY = gql`
  query App {
    auth @client {
      accessToken
      refreshToken
      user {
        ...HeaderUser
      }
    }
    cart @client {
      ...HeaderCart
    }
  }
  ${Header.fragments.cart}
  ${Header.fragments.user}
`;

const LOGOUT = gql`
  mutation AppLogout {
    logout @client
  }
`;

export default function App() {
  const client = useApolloClient();

  const [cartOpen, setCartOpen] = React.useState(false);
  const [loginSignupOpen, setLoginSignupOpen] = React.useState(false);

  const { data, error, loading } = useQuery(APP_QUERY);
  const [logout, { error: logoutError, loading: loggingOut }] = useMutation(
    LOGOUT
  );

  const user = React.useMemo(
    () =>
      data &&
      data.auth &&
      data.auth.accessToken &&
      decode(data.auth.accessToken),
    [data]
  );

  const loggedIn = React.useMemo(() => user && Date.now() < user.exp * 1000, [
    user
  ]);

  const closeCart = React.useCallback(() => setCartOpen(false), [setCartOpen]);
  const showCart = React.useCallback(() => setCartOpen(true), [setCartOpen]);

  const closeLoginSignup = React.useCallback(() => setLoginSignupOpen(false), [
    setLoginSignupOpen
  ]);
  const showLoginSignup = React.useCallback(() => setLoginSignupOpen(true), [
    setLoginSignupOpen
  ]);

  const onLogout = React.useCallback(async () => {
    if (loggingOut) {
      return;
    }

    await logout();
    await client.cache.reset();
  }, [loggingOut, logout]);

  return (
    <div id="goto-here">
      <Header
        cart={data && data.cart}
        user={data && data.auth && data.auth.user}
        loggedIn={loggedIn}
        onLogout={onLogout}
        onShowCart={showCart}
        onShowLoginSignup={showLoginSignup}
      />
      <Error error={!loading ? error : undefined} />
      <Error error={logoutError} />
      <Switch>
        <Route path="/admin/products/:id" component={AdminProductPage} />
        <Route path="/admin/products" component={AdminProductsPage} />
        <Route
          path="/admin/transactions/:id"
          component={AdminTransactionPage}
        />
        <Route path="/admin/transactions" component={AdminTransactionsPage} />
        <Route path="/admin/variants/:id" component={AdminVariantPage} />
        <Route path="/checkout" component={CheckoutPage} />
        <Route path="/profile" component={ProfilePage} />
        <Route path="/shop/:slug" component={ProductPage} />
        <Route path="/shop" component={ShopPage} />
        <Route path="/" component={HomePage} />
      </Switch>
      <CartPopup show={cartOpen} onHide={closeCart} />
      <LoginSignupPopup show={loginSignupOpen} onHide={closeLoginSignup} />
    </div>
  );
}
