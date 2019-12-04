import React from "react";
import { gql } from "apollo-boost";

import Form from "react-bootstrap/Form";

export default function ProductOption({ option, ...rest }) {
  return (
    <Form.Group controlId={`product-option-${option.id}`}>
      <Form.Label>{option.label}</Form.Label>
      <Form.Control as="select" {...rest}>
        <option value="">Choose...</option>
        {option.values.map(value => (
          <option key={value.id} value={value.id}>
            {value.value}
          </option>
        ))}
      </Form.Control>
    </Form.Group>
  );
}

ProductOption.fragments = {
  option: gql`
    fragment ProductOptionOption on ProductOption {
      id
      label
      values {
        id
        value
      }
    }
  `
};
