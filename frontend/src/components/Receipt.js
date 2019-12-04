import React from "react";
import { gql } from "apollo-boost";
import currency from "currency.js";

import Card from "react-bootstrap/Card";
import Table from "react-bootstrap/Table";

import Address from "./Address";

export default function Receipt({ receipt }) {
  return (
    <div>
      {receipt.shippingAddress && (
        <Card className="mb-4">
          <Card.Header>Shipping to</Card.Header>
          <Card.Body>
            <Address address={receipt.shippingAddress} />
          </Card.Body>
        </Card>
      )}

      <Table responsive>
        {receipt.lineItems &&
          receipt.lineItems.some(lineItem => lineItem.variant) && (
            <thead className="thead-dark">
              <tr>
                <th>Product</th>
                <th className="text-center">Price</th>
                <th className="text-center">Quantity</th>
                <th>Total</th>
              </tr>
            </thead>
          )}
        <tbody>
          {receipt.lineItems &&
            receipt.lineItems.map(lineItem =>
              lineItem.variant ? (
                <tr key={lineItem.variant.id}>
                  <td style={{ minWidth: 120 }}>
                    <p className="font-weight-bold">{lineItem.variant.name}</p>
                    <p>
                      {lineItem.variant.selectedOptions &&
                        lineItem.variant.selectedOptions
                          .map(option => option.value)
                          .join(", ")}
                    </p>
                  </td>
                  <td className="text-center">
                    <p>${currency(lineItem.price / 100).format()}</p>
                  </td>
                  <td className="text-center">
                    <p>{lineItem.quantity}</p>
                  </td>
                  <td>
                    <p>
                      $
                      {currency(
                        (lineItem.price * lineItem.quantity) / 100
                      ).format()}
                    </p>
                  </td>
                </tr>
              ) : null
            )}
          <tr>
            <th style={{ minWidth: 120 }} />
            <th />
            <th>
              <p className="text-right">Subtotal</p>
            </th>
            <td>${currency(receipt.subtotal / 100).format()}</td>
          </tr>
          <tr>
            <th style={{ minWidth: 120 }} />
            <th />
            <th>
              <p className="text-right">Taxes</p>
            </th>
            <td>${currency(receipt.taxes / 100).format()}</td>
          </tr>
          <tr>
            <th style={{ minWidth: 120 }} />
            <th />
            <th>
              <p className="text-right">Shipping</p>
            </th>
            <td>${currency(receipt.shipping / 100).format()}</td>
          </tr>
          <tr>
            <th style={{ minWidth: 120 }} />
            <th />
            <th>
              <p className="text-right font-weight-bold">Total</p>
            </th>
            <td>${currency(receipt.total / 100).format()}</td>
          </tr>
        </tbody>
      </Table>
    </div>
  );
}

Receipt.fragments = {
  receipt: gql`
    fragment ReceiptReceipt on Receipt {
      subtotal
      taxes
      shipping
      total
      lineItems {
        price
        quantity
        variant {
          id
          name
          selectedOptions {
            value
          }
        }
      }
      shippingAddress {
        ...AddressAddress
      }
    }
    ${Address.fragments.address}
  `,
  transaction: gql`
    fragment ReceiptTransaction on Transaction {
      subtotal
      taxes
      shipping
      total
      lineItems {
        price
        quantity
        variant {
          id
          name
          selectedOptions {
            value
          }
        }
      }
      shippingAddress {
        ...AddressAddress
      }
    }
    ${Address.fragments.address}
  `
};
