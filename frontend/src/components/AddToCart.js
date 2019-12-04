import React from "react";

import Button from "react-bootstrap/Button";
import Col from "react-bootstrap/Col";
import Form from "react-bootstrap/Form";

import useInputInt from "../hooks/use-input-int";

export default function AddToCart({ disabled, onAddToCart }) {
  const [quantityInput, quantity] = useInputInt("1");

  const increment = React.useCallback(
    () => quantityInput.onChange({ target: { value: quantity + 1 } }),
    [quantityInput, quantity]
  );

  const decrement = React.useCallback(
    () =>
      quantity > 1 &&
      quantityInput.onChange({ target: { value: quantity - 1 } }),
    [quantityInput, quantity]
  );

  const addToCart = React.useCallback(
    () => onAddToCart && onAddToCart(quantity)
  );

  return (
    <React.Fragment>
      <Form as="div">
        <Form.Row>
          <Col style={{ flex: 0 }}>
            <Button disabled={quantity <= 1} variant="outline-secondary" onClick={decrement}>
              <i
                className="oi oi-minus"
                aria-label="minus"
                aria-hidden="true"
              />
            </Button>
          </Col>
          <Col xs={4}>
            <Form.Control plaintext readOnly type="number" min={1} {...quantityInput} />
          </Col>
          <Col style={{ flex: 0 }}>
            <Button variant="outline-secondary" onClick={increment}>
              <i className="oi oi-plus" aria-label="plus" aria-hidden="true" />
            </Button>
          </Col>
        </Form.Row>
      </Form>

      <br />

      <Button disabled={disabled} size="lg" onClick={addToCart}>
        Add to Cart
      </Button>
    </React.Fragment>
  );
}
