import React from "react";
import { gql } from "apollo-boost";
import currency from "currency.js";

export default function PriceRange({ as: component, priceRange, ...rest }) {
  const Component = component || "p";

  return priceRange ? (
    <Component {...rest}>
      ${currency(priceRange.min / 100).format()}
      {!!priceRange.max &&
        priceRange.min !== priceRange.max &&
        `- $${currency(priceRange.max / 100).format()}`}
    </Component>
  ) : <Component {...rest}>---</Component>;
}

PriceRange.fragments = {
  priceRange: gql`
    fragment PriceRangePriceRange on ProductPriceRange {
      min
      max
    }
  `
};
