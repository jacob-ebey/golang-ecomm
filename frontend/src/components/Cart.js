import React from "react";
import { gql } from "apollo-boost";
import currency from "currency.js";

import Button from "react-bootstrap/Button";
import Col from "react-bootstrap/Col";
import Image from "react-bootstrap/Image";
import Row from "react-bootstrap/Row";
import Table from "react-bootstrap/Table";
import Form from "react-bootstrap/Form";

export default function Cart({
  subtotal,
  cart,
  variants,
  onIncrementQuantity,
  onDecrementQuantity,
  onRemoveFromCart
}) {
  const infos = React.useMemo(
    () =>
      variants ? variants.reduce((p, c) => ({ ...p, [c.id]: c }), {}) : {},
    [variants]
  );

  const rows = React.useMemo(
    () =>
      cart
        ? cart.variants.map(variant => ({
            ...infos[variant.variantId],
            quantity: variant.quantity
          }))
        : [],
    [cart, infos]
  );

  return (
    <React.Fragment>
      <Table responsive>
        <thead className="thead-dark">
          <tr>
            <th style={{ width: 50 }}></th>
            <th style={{ width: 120 }}></th>
            <th>Product</th>
            <th>Price</th>
            <th className="text-center">Quantity</th>
            <th>Total</th>
          </tr>
        </thead>
        <tbody>
          {rows.map(row => (
            <tr key={row.id}>
              <td className="text-center align-middle" style={{ width: 50 }}>
                <Button
                  variant="outline-secondary"
                  onClick={() => onRemoveFromCart && onRemoveFromCart(row.id)}
                >
                  <i
                    className="oi oi-trash"
                    aria-label="trash"
                    aria-hidden="true"
                  />
                </Button>
              </td>
              <td className="text-center" style={{ minWidth: 120 }}>
                <Image
                  fluid
                  style={{
                    height: 100,
                    maxWidth: 100,
                    objectFit: "contain",
                    objectPosition: "top",
                    margin: "0 auto"
                  }}
                  src={`https://via.placeholder.com/${50 * row.id}`}
                />
              </td>
              <td>
                <p className="font-weight-bold">{row.name}</p>
                <p>
                  {row.selectedOptions &&
                    row.selectedOptions.map(option => option.value).join(", ")}
                </p>
              </td>
              <td>
                <p>${currency(row.price / 100).format()}</p>
              </td>
              <td>
                <Row className="flex-nowrap">
                  <Col>
                    <Button
                      disabled={row.quantity <= 1}
                      size="sm"
                      variant="outline-secondary"
                      onClick={() =>
                        onDecrementQuantity && onDecrementQuantity(row.id)
                      }
                    >
                      <i
                        className="oi oi-minus"
                        aria-label="minus"
                        aria-hidden="true"
                      />
                    </Button>
                  </Col>
                  <Col>
                    <Form.Control
                      className="text-center"
                      disabled
                      plaintext
                      type="number"
                      min="1"
                      value={row.quantity}
                    />
                    {/* <p className="text-center">{row.quantity}</p> */}
                  </Col>
                  <Col className="text-right">
                    <Button
                      size="sm"
                      variant="outline-secondary"
                      onClick={() =>
                        onIncrementQuantity && onIncrementQuantity(row.id)
                      }
                    >
                      <i
                        className="oi oi-plus"
                        aria-label="plus"
                        aria-hidden="true"
                      />
                    </Button>
                  </Col>
                </Row>
              </td>
              <td>
                <p>${currency((row.price * row.quantity) / 100).format()}</p>
              </td>
            </tr>
          ))}
          {rows.length > 0 && typeof subtotal === "number" && (
            <tr>
              <td style={{ width: 50 }} />
              <td style={{ width: 120 }} />
              <td />
              <td />
              <th>
                <p className="text-right">Subtotal</p>
              </th>
              <td>${currency(subtotal / 100).format()}</td>
            </tr>
          )}
        </tbody>
      </Table>

      {rows.length === 0 && <p>No items in cart</p>}
    </React.Fragment>
  );
}

Cart.fragments = {
  cart: gql`
    fragment CartCart on CartState {
      variants {
        quantity
        variantId
      }
    }
  `,
  variants: gql`
    fragment CartVariants on ProductVariant {
      id
      name
      price
      selectedOptions {
        value
      }
    }
  `
};
