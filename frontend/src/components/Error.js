import React from "react";

import Alert from "react-bootstrap/Alert";

export default function Error({ error, ignoreNotAuthenticated }) {
  if (!error) {
    return null;
  }

  if (error.graphQLErrors && error.graphQLErrors.length > 0) {
    let foundAuthenticated = false;
    const errors = error.graphQLErrors
      .map((err, key) => {
        if (
          ignoreNotAuthenticated &&
          err.extensions &&
          err.extensions.code === "NOT_AUTHENTICATED"
        ) {
          foundAuthenticated = true;
          return null;
        }

        return (
          <Alert key={key} variant="danger">
            {err.message || "Something went wrong :("}
          </Alert>
        );
      })
      .filter(e => !!e);

    if (errors.length > 0 || foundAuthenticated) {
      return errors;
    }
  }

  return (
    <Alert variant="danger">{error.message || "Something went wrong :("}</Alert>
  );
}
