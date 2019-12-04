import React from "react";
import { useMap, useUpdateEffect } from "react-use";
import { gql } from "apollo-boost";

import ProductOption from "./ProductOption";

export default function ProductOptions({ options, onChange }) {
  const [
    selectedOptions,
    { set: setSelectedOption, remove: removeSelectedOption }
  ] = useMap();

  const productOptionChanged = React.useCallback(id => event => {
    const newValue = event.target.value;

    if (!newValue) {
      removeSelectedOption(id);
    } else {
      setSelectedOption(id, newValue);
    }
  });

  useUpdateEffect(
    () =>
      onChange &&
      onChange(
        Object.keys(selectedOptions).map(key =>
          Number.parseInt(selectedOptions[key], 10)
        )
      ),
    [selectedOptions]
  );

  return options
    ? options.map(option => (
        <ProductOption
          key={option.id}
          option={option}
          value={selectedOptions[option.id] || ""}
          onChange={productOptionChanged(option.id)}
        />
      ))
    : null;
}

ProductOptions.fragments = {
  options: gql`
    fragment ProductOptionsOptions on ProductOption {
      id
      ...ProductOptionOption
    }
    ${ProductOption.fragments.option}
  `
};
