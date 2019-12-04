import React from "react";
import { gql } from "apollo-boost";

export default function Address({ address }) {
  return (
    <address>
      <strong>{address.name}</strong> <br />
      {address.line1} <br />
      {address.line2 && (
        <React.Fragment>
          {address.line2} <br />
        </React.Fragment>
      )}
      {address.line3 && (
        <React.Fragment>
          {address.line3} <br />
        </React.Fragment>
      )}
      {address.city}, {address.region}, {address.postalCode} <br />
      {address.country}
    </address>
  );
}

Address.fragments = {
  address: gql`
    fragment AddressAddress on Address {
      name
      line1
      line2
      line3
      city
      region
      postalCode
      country
    }
  `
};
