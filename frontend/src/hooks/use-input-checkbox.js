import React from "react";

export default function useInputCheckbox(defaultValue) {
  const [checked, setChecked] = React.useState(!!defaultValue);

  return React.useMemo(
    () => ({
      onChange: event => setChecked(!!event.target.checked),
      checked
    }),
    [checked, setChecked]
  );
}
