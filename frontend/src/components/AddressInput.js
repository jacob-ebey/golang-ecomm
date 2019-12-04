import React from "react";

import Col from "react-bootstrap/Col";
import Form from "react-bootstrap/Form";

export default function AddressInput({
  prefix,
  name,
  line1,
  line2,
  line3,
  city,
  region,
  postalCode,
  country
}) {
  return (
    <fieldset>
      <Form.Group>
        <Form.Label>Name</Form.Label>
        <Form.Control
          {...name}
          required
          id={`${prefix ? prefix + "-" : ""}input-address-name`}
          name={`${prefix ? prefix + "-" : ""}input-address-name`}
          autoComplete={`${prefix ? prefix + "-" : ""}input-address-name`}
          placeholder="Name"
        />
      </Form.Group>

      <Form.Group>
        <Form.Label>Address</Form.Label>
        <Form.Control
          {...line1}
          required
          id={`${prefix ? prefix + "-" : ""}input-address1`}
          name={`${prefix ? prefix + "-" : ""}address1`}
          autoComplete={`${prefix ? prefix + "-" : ""}address1`}
          placeholder="1234 Main St"
        />
      </Form.Group>

      <Form.Group>
        <Form.Label>Address 2</Form.Label>
        <Form.Control
          {...line2}
          id={`${prefix ? prefix + "-" : ""}input-address2`}
          name={`${prefix ? prefix + "-" : ""}address2`}
          autoComplete={`${prefix ? prefix + "-" : ""}address2`}
          placeholder="Apartment, studio, or floor"
        />
      </Form.Group>

      <Form.Group>
        <Form.Label>Address 3</Form.Label>
        <Form.Control
          {...line3}
          id={`${prefix ? prefix + "-" : ""}input-address3`}
          name={`${prefix ? prefix + "-" : ""}address3`}
          autoComplete={`${prefix ? prefix + "-" : ""}address3`}
          placeholder="Etc..."
        />
      </Form.Group>

      <Form.Row>
        <Form.Group as={Col} xs={12} md={4}>
          <Form.Label>City</Form.Label>
          <Form.Control
            {...city}
            required
            id={`${prefix ? prefix + "-" : ""}input-city`}
            name={`${prefix ? prefix + "-" : ""}city`}
            autoComplete={`${prefix ? prefix + "-" : ""}city`}
            placeholder="City"
          />
        </Form.Group>

        <Form.Group as={Col} xs={12} md={4}>
          <Form.Label>State</Form.Label>
          <Form.Control
            {...region}
            required
            id={`${prefix ? prefix + "-" : ""}input-region`}
            name={`${prefix ? prefix + "-" : ""}region`}
            autoComplete={`${prefix ? prefix + "-" : ""}region`}
            placeholder="State"
          />
        </Form.Group>

        <Form.Group as={Col} xs={12} md={4}>
          <Form.Label>Zip</Form.Label>
          <Form.Control
            {...postalCode}
            required
            id={`${prefix ? prefix + "-" : ""}input-postal-code`}
            name={`${prefix ? prefix + "-" : ""}postal-code`}
            autoComplete={`${prefix ? prefix + "-" : ""}postal-code`}
            placeholder="12345"
          />
        </Form.Group>
      </Form.Row>

      <Form.Group>
        <Form.Label>Country</Form.Label>
        <Form.Control
          {...country}
          required
          id={`${prefix ? prefix + "-" : ""}input-country`}
          name={`${prefix ? prefix + "-" : ""}country`}
          autoComplete={`${prefix ? prefix + "-" : ""}country`}
          placeholder="USA"
        />
      </Form.Group>
    </fieldset>
  );
}
