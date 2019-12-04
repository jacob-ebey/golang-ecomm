import React from "react";
import useInput from "@rooks/use-input";

import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";

import useInputCheckbox from "../hooks/use-input-checkbox";

export default function Signup({ disabled, onSignup }) {
  const email = useInput();
  const password = useInput();
  const confirmPassword = useInput();
  const stayLoggedIn = useInputCheckbox(true);
  const [validated, setValidated] = React.useState(false);

  const confirmPasswordInvalid = React.useMemo(
    () => password.value !== confirmPassword.value,
    [password, confirmPassword]
  );

  const handleSubmit = React.useCallback(
    event => {
      event.preventDefault();
      const form = event.currentTarget;
      setValidated(true);

      if (form.checkValidity() === false || confirmPasswordInvalid) {
        event.stopPropagation();
        return;
      }

      onSignup &&
        onSignup({
          email: email.value,
          password: password.value,
          confirmPassword: confirmPassword.value,
          stayLoggedIn: stayLoggedIn.checked
        });
    },
    [email, password, confirmPassword, stayLoggedIn, confirmPasswordInvalid]
  );

  return (
    <Form noValidate validated={validated} onSubmit={handleSubmit}>
      <Form.Group>
        <Form.Label>Email address</Form.Label>
        <Form.Control
          {...email}
          autoComplete="new-email"
          required
          type="email"
          placeholder="Enter email"
          disabled={disabled}
        />
        <Form.Control.Feedback type="invalid">
          Please enter a valid email address.
        </Form.Control.Feedback>
      </Form.Group>

      <Form.Group>
        <Form.Label>Password</Form.Label>
        <Form.Control
          {...password}
          autoComplete="new-password"
          required
          type="password"
          placeholder="Password"
          disabled={disabled}
        />
        <Form.Control.Feedback type="invalid">
          Please enter a password.
        </Form.Control.Feedback>
      </Form.Group>

      <Form.Group>
        <Form.Label>Confirm Password</Form.Label>
        <Form.Control
          {...confirmPassword}
          autoComplete="new-password"
          required
          type="password"
          placeholder="Confirm Password"
          pattern={password.value}
          disabled={disabled}
        />
        <Form.Control.Feedback type="invalid">
          Passwords do not match.
        </Form.Control.Feedback>
      </Form.Group>

      <Form.Group>
        <Form.Check
          {...stayLoggedIn}
          type="checkbox"
          label="Stay logged in"
          disabled={disabled}
        />
      </Form.Group>
      <Button variant="primary" type="submit" disabled={disabled}>
        Signup
      </Button>
    </Form>
  );
}
