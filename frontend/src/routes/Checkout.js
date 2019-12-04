import React from "react";
import { gql } from "apollo-boost";
import { useMutation, useQuery } from "@apollo/react-hooks";
import { useDebounce } from "use-debounce";
import { useFormState } from "react-use-form-state";
import currency from "currency.js";

import LoadingOverlay from "react-loading-overlay";
import DropIn from "braintree-web-drop-in-react";

import Accordion from "react-bootstrap/Accordion";
import Button from "react-bootstrap/Button";
import Card from "react-bootstrap/Card";
import Col from "react-bootstrap/Col";
import Container from "react-bootstrap/Container";
import Form from "react-bootstrap/Form";
import Jumbotron from "react-bootstrap/Jumbotron";
import Row from "react-bootstrap/Row";
import Spinner from "react-bootstrap/Spinner";

import { checkout } from "../config";

import Address from "../components/Address";
import AddressInput from "../components/AddressInput";
import Cart from "../components/Cart";
import Error from "../components/Error";
import Receipt from "../components/Receipt";
import Section from "../components/Section";
import ShippingEstimation from "../components/ShippingEstimation";

import useInputInt from "../hooks/use-input-int";

const QUERY = gql`
  query Checkout($variants: [CartInput!]!, $variantIds: [Int!]!) {
    braintreeClientToken
    subtotal(variants: $variants)
    variants: productVariantsByIds(variantIds: $variantIds) {
      ...CartVariants
    }
    me {
      id
      addresses {
        id
        ...AddressAddress
      }
    }
  }
  ${Address.fragments.address}
  ${Cart.fragments.variants}
`;

const TAXES_QUERY = gql`
  query CheckoutTaxes($address: AddressInput!) {
    taxes(address: $address) {
      totalRate
    }
  }
`;

const SHIPPING_QUERY = gql`
  query CheckoutShipping($toAddr: AddressInput!, $toEstimate: [CartInput!]!) {
    shippingEstimations(address: $toAddr, variants: $toEstimate) {
      id
      ...ShippingEstimationEstimation
    }
  }
  ${ShippingEstimation.fragments.estimation}
`;

const CLIENT_QUERY = gql`
  query CheckoutClient {
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
  mutation CheckoutChangeCartQuantity($variantId: Int!, $quantity: Int!) {
    changeCartQuantity(variantId: $variantId, quantity: $quantity) @client
  }
`;

const REMOVE_FROM_CART = gql`
  mutation CheckoutRemoveFromCart($variantId: Int!) {
    removeFromCart(variantId: $variantId) @client
  }
`;

const CHECKOUT = gql`
  mutation CheckoutCheckout(
    $braintreeNonce: String!
    $shippingAddressId: Int
    $shippingAddress: AddressInput
    $saveShippingAddress: Boolean
    $billingAddressId: Int
    $billingAddress: AddressInput
    $saveBillingAddress: Boolean
    $lineItems: [CartInput!]!
    $shippingRateId: String!
    $total: Int!
  ) {
    receipt: submitBraintreeTransaction(
      braintreeNonce: $braintreeNonce
      variants: $lineItems
      shippingRateId: $shippingRateId
      total: $total
      shippingAddressId: $shippingAddressId
      shippingAddress: $shippingAddress
      saveShippingAddress: $saveShippingAddress
      billingAddressId: $billingAddressId
      billingAddress: $billingAddress
      saveBillingAddress: $saveBillingAddress
    ) {
      id
      ...ReceiptReceipt
    }
  }
  ${Receipt.fragments.receipt}
`;

const CLEAR_CART = gql`
  mutation CheckoutClearCart {
    clearCart @client
  }
`;

function Profile() {
  const [validated, setValidated] = React.useState(false);
  const [billingAddressKey, setBillingAddressKey] = React.useState("0");
  const [
    newBillingAddress,
    { text: billingAddressText, checkbox: billingAddressCheckbox }
  ] = useFormState();
  const [existingBillingAddress, existingBillingAddressId] = useInputInt();

  const [shippingAddressKey, setShippingAddressKey] = React.useState("0");
  const [
    newShippingAddress,
    { text: shippingAddressText, checkbox: shippingAddressCheckbox }
  ] = useFormState();
  const [existingShippingAddress, existingShippingAddressId] = useInputInt();

  const [shippingEstimation, shippingEstimationIndex] = useInputInt();

  const [
    paymentMethodRequestable,
    setPaymentMethodRequestable
  ] = React.useState(false);
  const braintreeClient = React.useRef();

  const {
    data: clientData,
    error: clientError,
    loading: clientLoading
  } = useQuery(CLIENT_QUERY);

  const [setQuantity, { error: setQuantityError }] = useMutation(SET_QUANTITY);
  const [removeFromCart, { error: removeFromCartError }] = useMutation(
    REMOVE_FROM_CART
  );
  const [
    checkoutMutation,
    { data: checkoutData, error: checkoutError, loading: checkoutLoading }
  ] = useMutation(CHECKOUT);
  const [clearCartMutation] = useMutation(CLEAR_CART);

  const cart = React.useMemo(() => clientData && clientData.cart, [clientData]);

  const variantIds = React.useMemo(
    () =>
      !clientLoading && cart && cart.variants.map(variant => variant.variantId),
    [clientLoading, cart]
  );

  const queryVariables = React.useMemo(
    () => ({
      variants:
        cart &&
        cart.variants.map(variant => ({
          variantId: variant.variantId,
          quantity: variant.quantity
        })),
      variantIds
    }),
    [cart, variantIds]
  );

  const { data, error } = useQuery(QUERY, {
    skip: clientLoading || !variantIds || variantIds.length === 0,
    variables: queryVariables,
    errorPolicy: "all"
  });

  const [hasAddresses, subtotal, variants] = React.useMemo(
    () =>
      data
        ? [
            data &&
              data.me &&
              data.me.addresses &&
              data.me.addresses.length > 0,
            data.subtotal,
            data.variants
          ]
        : [],
    [data]
  );

  // TODO: Create a hook for dealing with addresses
  const [shippingAddress, saveShippingAddress] = React.useMemo(() => {
    if (hasAddresses && shippingAddressKey === "0") {
      if (data && data.me && data.me.addresses) {
        const address =
          data &&
          data.me &&
          data.me.addresses.find(
            address => address.id === existingShippingAddressId
          );

        return [
          address && {
            id: address.id,
            name: address.name,
            line1: address.line1,
            line2: address.line2,
            line3: address.line3,
            city: address.city,
            region: address.region,
            postalCode: address.postalCode,
            country: address.country
          },
          false
        ];
      }
    }

    if (
      !newShippingAddress.validity.name ||
      !newShippingAddress.validity.line1 ||
      !newShippingAddress.validity.city ||
      !newShippingAddress.validity.region ||
      !newShippingAddress.validity.postalCode ||
      !newShippingAddress.validity.country
    ) {
      return [undefined, newShippingAddress.values.save];
    }

    return [
      {
        name: newShippingAddress.values.name,
        line1: newShippingAddress.values.line1,
        line2: newShippingAddress.values.line2,
        line3: newShippingAddress.values.line3,
        city: newShippingAddress.values.city,
        region: newShippingAddress.values.region,
        postalCode: newShippingAddress.values.postalCode,
        country: newShippingAddress.values.country
      },
      newShippingAddress.values.save
    ];
  }, [
    hasAddresses,
    shippingAddressKey,
    existingShippingAddressId,
    newShippingAddress.values,
    newShippingAddress.validity,
    data
  ]);

  const [shippingAddressDebounced] = useDebounce(shippingAddress, 1000);

  const [billingAddress, saveBillingAddress] = React.useMemo(() => {
    if (hasAddresses ? billingAddressKey === "2" : billingAddressKey === "1") {
      return [shippingAddress, false];
    }

    if (hasAddresses && billingAddressKey === "0") {
      if (data && data.me && data.me.addresses) {
        const address =
          data &&
          data.me &&
          data.me.addresses.find(
            address => address.id === existingBillingAddressId
          );

        return [
          address && {
            id: address.id,
            name: address.name,
            line1: address.line1,
            line2: address.line2,
            line3: address.line3,
            city: address.city,
            region: address.region,
            postalCode: address.postalCode,
            country: address.country
          },
          false
        ];
      }
    }

    if (
      !newBillingAddress.validity.name ||
      !newBillingAddress.validity.line1 ||
      !newBillingAddress.validity.city ||
      !newBillingAddress.validity.region ||
      !newBillingAddress.validity.postalCode ||
      !newBillingAddress.validity.country
    ) {
      return [undefined, newBillingAddress.values.save];
    }

    return [
      {
        name: newBillingAddress.values.name,
        line1: newBillingAddress.values.line1,
        line2: newBillingAddress.values.line2,
        line3: newBillingAddress.values.line3,
        city: newBillingAddress.values.city,
        region: newBillingAddress.values.region,
        postalCode: newBillingAddress.values.postalCode,
        country: newBillingAddress.values.country
      },
      newBillingAddress.values.save
    ];
  }, [
    shippingAddress,
    hasAddresses,
    billingAddressKey,
    existingBillingAddressId,
    newBillingAddress.values,
    newBillingAddress.validity,
    data
  ]);

  const [billingAddressDebounced] = useDebounce(billingAddress, 1000);

  const {
    data: taxesData,
    error: taxesError,
    loading: taxesLoading
  } = useQuery(TAXES_QUERY, {
    skip: !billingAddressDebounced,
    variables: React.useMemo(() => {
      const { id, ...address } = billingAddressDebounced || {};

      return { address };
    }, [billingAddressDebounced])
  });

  const {
    data: shippingData,
    error: shippingError,
    loading: shippingLoading
  } = useQuery(SHIPPING_QUERY, {
    skip: !shippingAddressDebounced,
    variables: React.useMemo(() => {
      const { id, ...toAddr } = shippingAddressDebounced || {};

      return {
        toAddr,
        toEstimate: queryVariables.variants
      };
    }, [shippingAddressDebounced, queryVariables])
  });

  const shippingRate = React.useMemo(() => {
    return (
      !shippingError &&
      shippingData &&
      shippingData.shippingEstimations &&
      shippingData.shippingEstimations[shippingEstimationIndex - 1]
    );
  }, [shippingEstimationIndex, shippingData]);

  const taxRate = React.useMemo(
    () => taxesData && taxesData.taxes && taxesData.taxes,
    [taxesData]
  );

  const taxes = React.useMemo(
    () =>
      subtotal && taxRate
        ? Math.round(subtotal * taxRate.totalRate)
        : undefined,
    [taxRate, subtotal]
  );

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

  const onBraintreeInstance = React.useCallback(
    instance => {
      braintreeClient.current = instance;
      instance.on("paymentMethodRequestable", () => {
        setPaymentMethodRequestable(true);
      });
      instance.on("noPaymentMethodRequestable", () => {
        setPaymentMethodRequestable(false);
      });
    },
    [braintreeClient, setPaymentMethodRequestable]
  );

  const total = React.useMemo(() => {
    if (
      !shippingRate ||
      !shippingRate.price ||
      billingAddress !== billingAddressDebounced ||
      !taxes ||
      !subtotal
    ) {
      return null;
    }

    return subtotal + taxes + shippingRate.price;
  }, [shippingRate, billingAddress, billingAddressDebounced, taxes, subtotal]);

  const [submitting, setSubmitting] = React.useState(false);
  const readyForCheckout = React.useMemo(() => {
    if (
      submitting ||
      !shippingAddress ||
      !billingAddress ||
      !taxes ||
      !shippingRate ||
      !braintreeClient.current ||
      !paymentMethodRequestable ||
      !queryVariables ||
      !queryVariables.variants ||
      !queryVariables.variants.length === 0
    ) {
      return false;
    }

    return true;
  }, [
    shippingAddress,
    billingAddress,
    taxes,
    shippingRate,
    braintreeClient.current,
    paymentMethodRequestable,
    submitting,
    queryVariables
  ]);

  const [braintreeError, setBraintreeError] = React.useState(null);
  const handleSubmit = React.useCallback(
    async event => {
      event.preventDefault();
      const form = event.currentTarget;
      setValidated(true);

      if (form.checkValidity() === false) {
        event.stopPropagation();
        return;
      }

      if (!readyForCheckout || submitting) {
        return;
      }

      setSubmitting(true);
      setPaymentMethodRequestable(false);

      try {
        let braintreeNonce;
        try {
          const result = await braintreeClient.current.requestPaymentMethod();
          braintreeNonce = result.nonce;
        } catch (err) {
          setBraintreeError(err.message);
          braintreeNonce = undefined;
        }

        if (!braintreeNonce) {
          return;
        }

        const {
          data: { receipt }
        } = await checkoutMutation({
          variables: {
            braintreeNonce,
            shippingRateId: shippingRate.id,
            total,
            billingAddressId: billingAddress.id,
            billingAddress: billingAddress.id ? null : billingAddress,
            saveBillingAddress,
            shippingAddressId: shippingAddress.id,
            shippingAddress: shippingAddress.id ? null : shippingAddress,
            saveShippingAddress,
            lineItems: queryVariables.variants
          }
        });

        if (receipt) {
          clearCartMutation();
        }
      } finally {
        setSubmitting(false);
      }
    },
    [
      queryVariables,
      setValidated,
      setSubmitting,
      readyForCheckout,
      billingAddress,
      saveBillingAddress,
      shippingAddress,
      saveShippingAddress,
      braintreeClient,
      shippingRate,
      checkoutMutation,
      setBraintreeError,
      setPaymentMethodRequestable,
      clearCartMutation
    ]
  );

  return (
    <LoadingOverlay
      active={checkoutLoading}
      spinner
      text="Submitting your payment..."
    >
      <Jumbotron
        fluid={true}
        style={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          minHeight: checkout.heroImage && "40vh",

          backgroundAttachment: "fixed",
          backgroundPosition: "center",
          backgroundRepeat: "no-repeat",
          backgroundSize: "cover",
          ...(checkout.heroImage
            ? { backgroundImage: `url(${checkout.heroImage})` }
            : {})
        }}
      >
        <div
          style={
            checkout.heroImage && {
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
                      content:
                        checkoutData && checkoutData.receipt
                          ? `# Receipt #${checkoutData.receipt.id}`
                          : "# Checkout"
                    }
                  ]
                }
              ]
            }}
          />
        </div>
      </Jumbotron>

      <Container>
        <Row>
          <Col xs={12}>
            <Error error={error} ignoreNotAuthenticated={true} />
            <Error error={clientError} />
          </Col>
        </Row>
      </Container>

      {checkoutData && checkoutData.receipt ? (
        <Container>
          <Row className="mb-4">
            <Col className="mb-4">
              <Receipt receipt={checkoutData.receipt} />
            </Col>
          </Row>
        </Container>
      ) : (
        data && (
          <Form noValidate validated={validated} onSubmit={handleSubmit}>
            <Container>
              <Row className="mb-4">
                <Col className="mb-4">
                  <Error error={setQuantityError} />
                  <Error error={removeFromCartError} />
                  <Cart
                    cart={cart}
                    variants={variants}
                    onIncrementQuantity={increment}
                    onDecrementQuantity={decrement}
                    onRemoveFromCart={onRemoveFromCart}
                  />
                </Col>
              </Row>
              {!checkoutLoading && (
                <React.Fragment>
                  <Row className="mb-4">
                    <Col className="mb-4" xs={12}>
                      <h2 className="font-weight-normal">Shipping Address</h2>

                      <Accordion activeKey={shippingAddressKey}>
                        {hasAddresses && (
                          <Card>
                            <Card.Header>
                              <Accordion.Toggle
                                as={Button}
                                variant="link"
                                eventKey="0"
                                onClick={() => setShippingAddressKey("0")}
                              >
                                Existing Address
                              </Accordion.Toggle>
                            </Card.Header>
                            <Accordion.Collapse eventKey="0">
                              <Card.Body>
                                {shippingAddressKey === "0" && (
                                  <fieldset>
                                    <Form.Group>
                                      {data.me.addresses.map(address => (
                                        <Form.Check
                                          {...existingShippingAddress}
                                          required
                                          key={address.id}
                                          id={`existing-shipping-address-${address.id}`}
                                          type="radio"
                                          label={<Address address={address} />}
                                          name="existingShippingAddress"
                                          value={address.id}
                                          checked={
                                            existingShippingAddressId ===
                                            address.id
                                          }
                                        />
                                      ))}
                                    </Form.Group>
                                  </fieldset>
                                )}
                              </Card.Body>
                            </Accordion.Collapse>
                          </Card>
                        )}
                        <Card>
                          <Card.Header>
                            <Accordion.Toggle
                              as={Button}
                              variant="link"
                              eventKey={hasAddresses ? "1" : "0"}
                              onClick={() =>
                                setShippingAddressKey(hasAddresses ? "1" : "0")
                              }
                            >
                              New Address
                            </Accordion.Toggle>
                          </Card.Header>
                          <Accordion.Collapse
                            eventKey={hasAddresses ? "1" : "0"}
                          >
                            <Card.Body>
                              <Error error={shippingError} />
                              {(hasAddresses ? "1" : "0") ===
                                shippingAddressKey && (
                                <React.Fragment>
                                  <AddressInput
                                    prefix="shipping"
                                    name={shippingAddressText("name")}
                                    line1={shippingAddressText("line1")}
                                    line2={shippingAddressText("line2")}
                                    line3={shippingAddressText("line3")}
                                    city={shippingAddressText("city")}
                                    region={shippingAddressText("region")}
                                    postalCode={shippingAddressText(
                                      "postalCode"
                                    )}
                                    country={shippingAddressText("country")}
                                  />
                                  {data && data.me && (
                                    <Form.Group controlId="saveShippingAddress">
                                      <Form.Check
                                        {...shippingAddressCheckbox("save")}
                                        type="checkbox"
                                        label="Save address?"
                                      />
                                    </Form.Group>
                                  )}
                                </React.Fragment>
                              )}
                            </Card.Body>
                          </Accordion.Collapse>
                        </Card>
                      </Accordion>
                    </Col>
                  </Row>

                  <Row className="mb-4">
                    <Col className="mb-4" xs={12}>
                      <h2 className="font-weight-normal">Billing Address</h2>

                      <Error error={taxesError} />

                      <Accordion activeKey={billingAddressKey}>
                        {hasAddresses && (
                          <Card>
                            <Card.Header>
                              <Accordion.Toggle
                                as={Button}
                                variant="link"
                                eventKey="0"
                                onClick={() => setBillingAddressKey("0")}
                              >
                                Existing Address
                              </Accordion.Toggle>
                            </Card.Header>
                            <Accordion.Collapse eventKey="0">
                              <Card.Body>
                                {billingAddressKey === "0" && (
                                  <fieldset>
                                    <Form.Group>
                                      {data.me.addresses.map(address => (
                                        <Form.Check
                                          {...existingBillingAddress}
                                          required
                                          key={address.id}
                                          id={`existing-billing-address-${address.id}`}
                                          type="radio"
                                          label={<Address address={address} />}
                                          name="existingBillingAddress"
                                          value={address.id}
                                          checked={
                                            existingBillingAddressId ===
                                            address.id
                                          }
                                        />
                                      ))}
                                    </Form.Group>
                                  </fieldset>
                                )}
                              </Card.Body>
                            </Accordion.Collapse>
                          </Card>
                        )}
                        <Card>
                          <Card.Header>
                            <Accordion.Toggle
                              as={Button}
                              variant="link"
                              eventKey={hasAddresses ? "1" : "0"}
                              onClick={() =>
                                setBillingAddressKey(hasAddresses ? "1" : "0")
                              }
                            >
                              New Address
                            </Accordion.Toggle>
                          </Card.Header>
                          <Accordion.Collapse
                            eventKey={hasAddresses ? "1" : "0"}
                          >
                            <Card.Body>
                              {(hasAddresses ? "1" : "0") ===
                                billingAddressKey && (
                                <React.Fragment>
                                  <AddressInput
                                    prefix="billing"
                                    name={billingAddressText("name")}
                                    line1={billingAddressText("line1")}
                                    line2={billingAddressText("line2")}
                                    line3={billingAddressText("line3")}
                                    city={billingAddressText("city")}
                                    region={billingAddressText("region")}
                                    postalCode={billingAddressText(
                                      "postalCode"
                                    )}
                                    country={billingAddressText("country")}
                                  />
                                  {data && data.me && (
                                    <Form.Group controlId="saveBillingAddress">
                                      <Form.Check
                                        {...billingAddressCheckbox("save")}
                                        type="checkbox"
                                        label="Save address?"
                                      />
                                    </Form.Group>
                                  )}
                                </React.Fragment>
                              )}
                            </Card.Body>
                          </Accordion.Collapse>
                        </Card>
                        <Card>
                          <Card.Header>
                            <Accordion.Toggle
                              as={Button}
                              variant="link"
                              eventKey={hasAddresses ? "2" : "1"}
                              onClick={() =>
                                setBillingAddressKey(hasAddresses ? "2" : "1")
                              }
                            >
                              Same as shipping address
                            </Accordion.Toggle>
                          </Card.Header>
                          <Accordion.Collapse
                            eventKey={hasAddresses ? "2" : "1"}
                          >
                            <Card.Body>
                              {shippingAddress ? (
                                <Address address={shippingAddress} />
                              ) : (
                                <p className="text-danger">
                                  Please select a shipping address.
                                </p>
                              )}
                            </Card.Body>
                          </Accordion.Collapse>
                        </Card>
                      </Accordion>
                    </Col>
                  </Row>

                  <Row className="mb-4">
                    <Col className="mb-4" xs={12}>
                      <h2 className="font-weight-normal">Payment Method</h2>
                      {braintreeError && (
                        <p className="text-danger">{braintreeError}</p>
                      )}
                      {billingAddress ? (
                        <DropIn
                          options={{
                            authorization: data.braintreeClientToken,
                            paypal: {
                              flow: "vault"
                            },
                            card: {
                              ccv: {
                                required: true
                              },
                              overrides: {
                                fields: {
                                  postalCode: {
                                    prefill: billingAddress.postalCode
                                  }
                                }
                              }
                            }
                          }}
                          onInstance={onBraintreeInstance}
                        />
                      ) : (
                        <p className="text-danger">
                          Please select a billing address.
                        </p>
                      )}
                    </Col>
                  </Row>

                  <Row className="mb-4">
                    <Col className="mb-4" xs={12}>
                      <h2 className="font-weight-normal">Shipping Method</h2>
                      {shippingAddress ? (
                        shippingLoading ||
                        shippingAddress !== shippingAddressDebounced ? (
                          <Spinner animation="grow" />
                        ) : (
                          <Card>
                            <Card.Body>
                              <Error error={shippingError} />
                              {!shippingError &&
                                shippingData &&
                                shippingData.shippingEstimations && (
                                  <fieldset>
                                    <Form.Group>
                                      {shippingData.shippingEstimations.map(
                                        (estimation, index) => (
                                          <Form.Check
                                            {...shippingEstimation}
                                            required
                                            key={index}
                                            id={`shipping-estimation-${index}`}
                                            type="radio"
                                            label={
                                              <ShippingEstimation
                                                estimation={estimation}
                                              />
                                            }
                                            name="shppingEstimation"
                                            value={index + 1}
                                            checked={
                                              shippingEstimationIndex ===
                                              index + 1
                                            }
                                            className="mb-3"
                                          />
                                        )
                                      )}
                                    </Form.Group>
                                  </fieldset>
                                )}
                            </Card.Body>
                          </Card>
                        )
                      ) : (
                        <p className="text-danger">
                          Please select a shipping address.
                        </p>
                      )}
                    </Col>
                  </Row>

                  <Row className="mb-4">
                    <Col xs={12}>
                      <h2 className="font-weight-normal">Review</h2>
                      <Error error={checkoutError} />
                    </Col>
                    <Col className="mb-4" xs={12} md={5} lg={4}>
                      <Card>
                        <Card.Body>
                          <Row>
                            <Col xs={6}>
                              <h4 className="font-weight-light">Subtotal</h4>
                            </Col>
                            <Col xs={6}>
                              <h4 className="font-weight-light">
                                {subtotal
                                  ? `$${currency(subtotal / 100).format()}`
                                  : "--"}
                              </h4>
                            </Col>
                          </Row>
                          <Row>
                            <Col xs={6}>
                              <h4 className="font-weight-light">Taxes</h4>
                            </Col>
                            <Col xs={6}>
                              {taxesLoading ||
                              billingAddress !== billingAddressDebounced ? (
                                <Spinner animation="grow" size="sm" />
                              ) : (
                                <h4 className="font-weight-light">
                                  {taxes
                                    ? `$${currency(taxes / 100).format()}`
                                    : "--"}
                                </h4>
                              )}
                            </Col>
                          </Row>
                          <br />
                          <Row>
                            <Col xs={6}>
                              <h4 className="font-weight-light">Shipping</h4>
                            </Col>
                            <Col xs={6}>
                              <h4 className="font-weight-light">
                                {shippingRate && shippingRate.price
                                  ? `$${currency(
                                      shippingRate.price / 100
                                    ).format()}`
                                  : "--"}
                              </h4>
                            </Col>
                          </Row>
                          <br />
                          <Row>
                            <Col xs={6}>
                              <h3>Total</h3>
                            </Col>
                            <Col xs={6}>
                              <h3 className="font-weight-light">
                                {total
                                  ? `$${currency(total / 100).format()}`
                                  : "--"}
                              </h3>
                            </Col>
                          </Row>

                          <Button
                            disabled={!readyForCheckout}
                            type="submit"
                            size="lg"
                            className="mt-2"
                          >
                            Checkout
                          </Button>
                        </Card.Body>
                      </Card>
                    </Col>
                  </Row>
                </React.Fragment>
              )}
            </Container>
          </Form>
        )
      )}

      {(!checkoutData || !checkoutData.receipt) &&
        (!variantIds || variantIds.length === 0) && (
          <Container>
            <Row className="mb-4">
              <Col className="mb-4">
                <h1>No Items in your cart</h1>
              </Col>
            </Row>
          </Container>
        )}
    </LoadingOverlay>
  );
}

export default Profile;
