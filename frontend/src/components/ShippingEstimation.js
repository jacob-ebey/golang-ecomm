import React from "react";
import { gql } from "apollo-boost";
import currency from "currency.js";

export default function ShippingEstimation({ estimation }) {
  return (
    <span>
      <strong>{estimation.carrier}</strong> - {estimation.service}
      <br />${currency(estimation.price / 100).format()}
      {estimation.durationTerms && (
        <React.Fragment>
          <br /> {estimation.durationTerms}
        </React.Fragment>
      )}
    </span>
  );
}

ShippingEstimation.fragments = {
  estimation: gql`
    fragment ShippingEstimationEstimation on ShippingRate {
      carrier
      service
      price
      durationTerms
    }
  `
};
