import React from "react";
import useInput from "@rooks/use-input";

import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";

import useInputCheckbox from "../hooks/use-input-checkbox";

export default function Login({ disabled, onLogin }) {
  const email = useInput();
  const password = useInput();
  const stayLoggedIn = useInputCheckbox(true);
  const [validated, setValidated] = React.useState(false);

  const handleSubmit = event => {
    event.preventDefault();
    const form = event.currentTarget;
    setValidated(true);

    if (form.checkValidity() === false) {
      event.stopPropagation();
      return;
    }

    onLogin &&
      onLogin({
        email: email.value,
        password: password.value,
        stayLoggedIn: stayLoggedIn.checked
      });
  };

  return (
    <Form noValidate validated={validated} onSubmit={handleSubmit}>
      <Form.Group controlId="formBasicEmail">
        <Form.Label>Email address</Form.Label>
        <Form.Control
          {...email}
          autoComplete="current-email"
          required
          type="email"
          placeholder="Enter email"
          disabled={disabled}
        />
        <Form.Control.Feedback type="invalid">
          Please enter a valid email address.
        </Form.Control.Feedback>
      </Form.Group>

      <Form.Group controlId="formBasicPassword">
        <Form.Label>Password</Form.Label>
        <Form.Control
          {...password}
          autoComplete="current-password"
          required
          type="password"
          placeholder="Password"
          disabled={disabled}
        />
        <Form.Control.Feedback type="invalid">
          Please enter a password.
        </Form.Control.Feedback>
      </Form.Group>

      <Form.Group controlId="formBasicCheckbox">
        <Form.Check
          {...stayLoggedIn}
          type="checkbox"
          label="Stay logged in"
          disabled={disabled}
        />
      </Form.Group>
      <Button variant="primary" type="submit" disabled={disabled}>
        Login
      </Button>
    </Form>
  );
}
