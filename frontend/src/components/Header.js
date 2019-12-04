import React from "react";
import { gql } from "apollo-boost";
import { Link } from "react-router-dom";
import Nav from "react-bootstrap/Nav";
import Navbar from "react-bootstrap/Navbar";
import NavDropdown from "react-bootstrap/NavDropdown";

import { app } from "../config";

export default function Header({
  cart,
  user,
  loggedIn,
  onLogout,
  onShowCart,
  onShowLoginSignup
}) {
  const totalInCart = React.useMemo(() =>
    cart && cart.variants
      ? cart.variants.reduce((p, c) => p + c.quantity, 0)
      : 0
  );

  return (
    <Navbar expand="sm" bg="light" sticky="top">
      <Navbar.Brand as={Link} to="/">
        {app.title}
      </Navbar.Brand>
      <Navbar.Toggle aria-controls="responsive-navbar-nav" />
      <Navbar.Collapse id="responsive-navbar-nav">
        <Nav className="mr-auto" />
        <Nav>
          <Nav.Link as={Link} to="/shop">
            Shop
          </Nav.Link>
          {totalInCart > 0 && (
            <Nav.Link
              as="span"
              role="button"
              style={{ cursor: "pointer" }}
              onClick={onShowCart}
            >
              <i className="oi oi-cart" aria-label="cart" aria-hidden="true" />{" "}
              Cart ({totalInCart})
            </Nav.Link>
          )}
          {loggedIn && user && user.role === "ADMIN" && (
            <NavDropdown title="Admin">
              <NavDropdown.Item as={Link} to="/admin/products">
                Products
              </NavDropdown.Item>
              <NavDropdown.Item as={Link} to="/admin/transactions">
                Transactions
              </NavDropdown.Item>
            </NavDropdown>
          )}
          {loggedIn ? (
            <NavDropdown
              title={
                <React.Fragment>
                  User{" "}
                  <i
                    className="oi oi-person"
                    aria-label="profile"
                    aria-hidden="true"
                  />
                </React.Fragment>
              }
              alignRight={true}
              aria-label="profile"
            >
              <NavDropdown.Item as={Link} to="/profile">
                Profile
              </NavDropdown.Item>
              <NavDropdown.Divider />
              <NavDropdown.Item
                as="span"
                role="button"
                style={{ cursor: "pointer" }}
                onClick={onLogout}
              >
                Logout
              </NavDropdown.Item>
            </NavDropdown>
          ) : (
            <Nav.Link
              as="span"
              role="button"
              style={{ cursor: "pointer" }}
              onClick={onShowLoginSignup}
            >
              Login
            </Nav.Link>
          )}
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  );
}

Header.fragments = {
  cart: gql`
    fragment HeaderCart on CartState {
      variants {
        quantity
      }
    }
  `,
  user: gql`
    fragment HeaderUser on LocalUser {
      id
      role
    }
  `
};
