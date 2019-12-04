import React from "react";
import { gql } from "apollo-boost";
import { useApolloClient, useMutation } from "@apollo/react-hooks";
import { useFormState } from "react-use-form-state";

import Modal from "react-bootstrap/Modal";

import AddressInput from "../components/AddressInput";
import Error from "../components/Error";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";

const CREATE_ADDRESS = gql`
  mutation NewAddressPopupCreateAddress($address: AddressInput!) {
    createAddress(address: $address) {
      id
    }
  }
`;

function NewAddressPopup({ show, onHide }) {
  const client = useApolloClient();
  const [createAddress, { error, loading }] = useMutation(CREATE_ADDRESS);

  const [validated, setValidated] = React.useState(false);
  const [form, { text }] = useFormState();

  const handleSubmit = React.useCallback(
    async event => {
      event.preventDefault();

      setValidated(true);
      if (event.currentTarget.checkValidity() === false) {
        event.stopPropagation();
        return;
      }

      if (loading) {
        return;
      }

      const { data } = await createAddress({
        variables: {
          address: {
            name: form.values.name,
            line1: form.values.line1,
            line2: form.values.line2,
            line3: form.values.line3,
            city: form.values.city,
            region: form.values.region,
            postalCode: form.values.postalCode,
            country: form.values.country
          }
        }
      });

      if (data && data.createAddress) {
        onHide && onHide();
        form.reset();
        await client.resetStore();
      }
    },
    [form, loading, createAddress]
  );

  return (
    <Modal show={show} onHide={onHide} animation={true} size="xl">
      <Modal.Header closeButton>
        <Modal.Title>New Address</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Form noValidate validated={validated} onSubmit={handleSubmit}>
          <Error error={error} />
          <AddressInput
            name={text("name")}
            line1={text("line1")}
            line2={text("line2")}
            line3={text("line3")}
            city={text("city")}
            region={text("region")}
            postalCode={text("postalCode")}
            country={text("country")}
          />

          <Button type="submit">Save</Button>
        </Form>
      </Modal.Body>
    </Modal>
  );
}

export default NewAddressPopup;
