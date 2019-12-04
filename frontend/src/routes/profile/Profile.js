import React from "react";
import { gql } from "apollo-boost";
import { useMutation, useQuery, useApolloClient } from "@apollo/react-hooks";

import { profile } from "../../config";

import Button from "react-bootstrap/Button";
import Jumbotron from "react-bootstrap/Jumbotron";
import Modal from "react-bootstrap/Modal";

import Address from "../../components/Address";
import Error from "../../components/Error";
import Receipt from "../../components/Receipt";
import Section from "../../components/Section";

import NewAddressPopup from "../../containers/NewAddressPopup";
import Card from "react-bootstrap/Card";

const QUERY = gql`
  query Profile {
    me {
      id
      email
      addresses {
        id
        ...AddressAddress
      }
      receipts {
        id
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
      }
    }
  }
  ${Address.fragments.address}
`;

const DELETE_ADDRESS = gql`
  mutation ProfileDeleteAddress($addressId: Int!) {
    deleteAddress(addressId: $addressId) {
      id
    }
  }
`;

function Profile() {
  const [showNewAddressPopup, setShowNewAddressPopup] = React.useState(false);
  const [addressToDelete, setAddressToDelete] = React.useState(null);

  const client = useApolloClient();
  const { data, error } = useQuery(QUERY);
  const [
    deleteAddress,
    { error: deleteAddressError, loading: deletingAddress }
  ] = useMutation(DELETE_ADDRESS);

  const stageAddressForDeletion = React.useCallback(
    address => () => setAddressToDelete(address),
    [setAddressToDelete]
  );
  const closeDeleteAddress = React.useCallback(() => setAddressToDelete(null), [
    setAddressToDelete
  ]);

  const openShowNewAddress = React.useCallback(
    () => setShowNewAddressPopup(true),
    [setShowNewAddressPopup]
  );
  const closeShowNewAddress = React.useCallback(
    () => setShowNewAddressPopup(false),
    [setShowNewAddressPopup]
  );

  const onDeleteAddress = React.useCallback(async () => {
    if (!addressToDelete || deletingAddress) {
      return;
    }

    await deleteAddress({
      variables: { addressId: addressToDelete.id }
    });
    client.resetStore();
    closeDeleteAddress();
  }, [deleteAddress, deletingAddress, addressToDelete]);

  return (
    <React.Fragment>
      <Jumbotron
        fluid={true}
        style={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          minHeight: profile.heroImage && "40vh",

          backgroundAttachment: "fixed",
          backgroundPosition: "center",
          backgroundRepeat: "no-repeat",
          backgroundSize: "cover",
          ...(profile.heroImage
            ? { backgroundImage: `url(${profile.heroImage})` }
            : {})
        }}
      >
        <div
          style={
            profile.heroImage && {
              backgroundColor: "rgba(0,0,0,0.4)",
              color: "white"
            }
          }
        >
          <Section
            section={{
              rows: [
                {
                  columns: [
                    {
                      align: "center",
                      content: `
# Profile
`
                    }
                  ]
                }
              ]
            }}
          />
        </div>
      </Jumbotron>
      <Error error={error} />

      {data && data.me && (
        <Section
          section={{
            rows: [
              {
                columns: [
                  {
                    breakpoints: {
                      xs: 12,
                      sm: 12,
                      md: 6
                    },
                    content: `
# Overview

**ID:** ${data.me.id}

**Email:** ${data.me.email}
`
                  },
                  {
                    breakpoints: {
                      xs: 12,
                      sm: 12,
                      md: 6
                    },
                    Component: () => (
                      <React.Fragment>
                        <h1>
                          Addresses
                          <Button
                            className="ml-2"
                            aria-label="add"
                            onClick={openShowNewAddress}
                          >
                            <i
                              className="oi oi-plus"
                              aria-label="add"
                              aria-hidden="true"
                            />
                          </Button>
                        </h1>
                        <Error error={deleteAddressError} />

                        {!data.me.addresses ||
                          (data.me.addresses.length === 0 && (
                            <p>No addresses found.</p>
                          ))}

                        {data.me.addresses &&
                          data.me.addresses.map(address => (
                            <div key={address.id} className="d-flex flex-row">
                              <div className="mr-2 mt-1">
                                <Button
                                  variant="danger"
                                  size="sm"
                                  onClick={stageAddressForDeletion(address)}
                                  aria-label="trash"
                                >
                                  <i
                                    className="oi oi-trash"
                                    aria-label="trash"
                                    aria-hidden="true"
                                  />
                                </Button>
                              </div>
                              <Address address={address} />
                            </div>
                          ))}
                      </React.Fragment>
                    )
                  }
                ]
              },
              {
                columns:
                  data.me.receipts &&
                  data.me.receipts.map(receipt => ({
                    breakpoints: {
                      xs: 12,
                      lg: 6
                    },
                    Component: () => (
                      <Card>
                        <Card.Header>Receipt #{receipt.id}</Card.Header>
                        <Card.Body>
                          <Receipt receipt={receipt} />
                        </Card.Body>
                      </Card>
                    )
                  }))
              }
            ]
          }}
        />
      )}

      <NewAddressPopup
        show={showNewAddressPopup}
        onHide={closeShowNewAddress}
      />

      <Modal
        show={!!addressToDelete}
        onHide={closeDeleteAddress}
        animation={true}
        size="sm"
      >
        <Modal.Header closeButton>
          <Modal.Title>Delete Address</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <p>Are you sure you wish to delete this address?</p>
          {addressToDelete && <Address address={addressToDelete} />}
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={closeDeleteAddress}>
            Cancel
          </Button>
          <Button variant="danger" onClick={onDeleteAddress}>
            Delete
          </Button>
        </Modal.Footer>
      </Modal>
    </React.Fragment>
  );
}

export default Profile;
