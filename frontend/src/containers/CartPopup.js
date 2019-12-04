import React from "react";
import { gql } from "apollo-boost";
import { useMutation, useQuery } from "@apollo/react-hooks";
import { Link } from "react-router-dom";

import Button from "react-bootstrap/Button";
import Modal from "react-bootstrap/Modal";

import Cart from "../components/Cart";
import Error from "../components/Error";

const QUERY = gql`
  query CartPopup($variants: [CartInput!]!, $variantIds: [Int!]!) {
    subtotal(variants: $variants)
    variants: productVariantsByIds(variantIds: $variantIds) {
      ...CartVariants
    }
  }
  ${Cart.fragments.variants}
`;

const CLIENT_QUERY = gql`
  query CartPopupClient {
    cart @client {
      ...CartCart
      variants {
        quantity
        variantId
      }
    }
  }
  ${Cart.fragments.cart}
`;

const SET_QUANTITY = gql`
  mutation CartPopupChangeCartQuantity($variantId: Int!, $quantity: Int!) {
    changeCartQuantity(variantId: $variantId, quantity: $quantity) @client
  }
`;

const REMOVE_FROM_CART = gql`
  mutation CartPopupRemoveFromCart($variantId: Int!) {
    removeFromCart(variantId: $variantId) @client
  }
`;

function CartPopup({ show, onHide }) {
  const {
    data: clientData,
    error: clientError,
    loading: clientLoading
  } = useQuery(CLIENT_QUERY);

  const [setQuantity, { error: setQuantityError }] = useMutation(SET_QUANTITY);
  const [removeFromCart, { error: removeFromCartError }] = useMutation(
    REMOVE_FROM_CART
  );

  const cart = clientData && clientData.cart;

  const variantIds = React.useMemo(
    () =>
      !clientLoading && cart && cart.variants.map(variant => variant.variantId),
    [clientLoading, cart]
  );

  const queryVariables = React.useMemo(() => ({
    variants:
      cart &&
      cart.variants.map(variant => ({
        variantId: variant.variantId,
        quantity: variant.quantity
      })),
    variantIds
  }));

  const { data, error } = useQuery(QUERY, {
    skip: clientLoading || !variantIds || variantIds.length === 0,
    variables: queryVariables,
    errorPolicy: "all"
  });

  const subtotal = data && data.subtotal;
  const variants = data && data.variants;

  const increment = React.useCallback(variantId => {
    setQuantity({
      variables: { variantId, quantity: 1 }
    });
  });

  const decrement = React.useCallback(variantId => {
    setQuantity({
      variables: { variantId, quantity: -1 }
    });
  });

  const onRemoveFromCart = React.useCallback(variantId => {
    removeFromCart({ variables: { variantId } });
  });

  return (
    <Modal show={show} onHide={onHide} animation={true} size="xl">
      <Modal.Header closeButton>
        <Modal.Title>Cart</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <Error error={error} />
        <Error error={clientError} />
        <Error error={setQuantityError} />
        <Error error={removeFromCartError} />
        <Cart
          subtotal={subtotal}
          cart={cart}
          variants={variants}
          onIncrementQuantity={increment}
          onDecrementQuantity={decrement}
          onRemoveFromCart={onRemoveFromCart}
        />
      </Modal.Body>
      <Modal.Footer>
        {!cart || cart.variants.length === 0 ? (
          <Button disabled={true}>Checkout</Button>
        ) : (
          <Button as={Link} to="/checkout" variant="primary" onClick={onHide}>
            Checkout
          </Button>
        )}
      </Modal.Footer>
    </Modal>
  );
}

export default CartPopup;
