import React from "react";
import { Link, useParams } from "react-router-dom";

import { useMutation, useQuery } from "@apollo/react-hooks";
import { gql } from "apollo-boost";

import Breadcrumb from "react-bootstrap/Breadcrumb";
import Button from "react-bootstrap/Button";
import Card from "react-bootstrap/Card";
import Col from "react-bootstrap/Col";
import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";

import ShippingEstimation from "../../components/ShippingEstimation";
import Error from "../../components/Error";
import Receipt from "../../components/Receipt";

const QUERY = gql`
  query Transaction($id: Int!) {
    transaction(id: $id) {
      id
      shippoRateId
      shippingEstimation {
        ...ShippingEstimationEstimation
      }
      shippingLabel {
        id
        labelUrl
      }
      ...ReceiptTransaction
    }
  }
  ${Receipt.fragments.transaction}
  ${ShippingEstimation.fragments.estimation}
`;

const PURCHASE_LABEL = gql`
  mutation TransactionPurchaseLabel(
    $transactionId: Int!
    $shippoRateId: String!
  ) {
    shippingLabel: purchaseShippoLabel(
      transactionId: $transactionId
      shippoRateId: $shippoRateId
    ) {
      id
      labelUrl
    }
  }
`;

export default function Transaction() {
  const { id } = useParams();

  const { data, error, loading } = useQuery(QUERY, {
    variables: { id }
  });

  const [
    purchaseLabel,
    { error: purchaseLabelError, loading: purchasingLabel }
  ] = useMutation(PURCHASE_LABEL, {
    variables: {
      transactionId: data && data.transaction && data.transaction.id,
      shippoRateId: data && data.transaction && data.transaction.shippoRateId
    },
    refetchQueries: [
      {
        query: QUERY,
        variables: { id }
      }
    ]
  });

  const onPurchaseLabel = React.useCallback(() => {
    if (purchasingLabel) {
      return;
    }

    purchaseLabel();
  }, [purchaseLabel, purchasingLabel]);

  return (
    <Container>
      <Row>
        <Col className="mt-4 mb-4" xs={12}>
          <Breadcrumb>
            <Breadcrumb.Item active>Admin</Breadcrumb.Item>
            <Link className="breadcrumb-item" to="/admin/transactions">
              Transactions
            </Link>
            <Breadcrumb.Item active>#{id}</Breadcrumb.Item>
          </Breadcrumb>
          <Error error={error} />
        </Col>
      </Row>

      {!loading && data && data.transaction && (
        <Row>
          <Col xs={12} lg={6}>
            <Card className="mb-4">
              <Card.Header>Shipping label</Card.Header>
              <Card.Body>
                <p>
                  <ShippingEstimation
                    estimation={data.transaction.shippingEstimation}
                  />
                </p>

                {data.transaction.shippingLabel ? (
                  <React.Fragment>
                    <p>
                      <a
                        rel="noopener noreferrer"
                        target="_blank"
                        href={data.transaction.shippingLabel.labelUrl}
                      >
                        <strong>View Label</strong>
                      </a>
                    </p>
                  </React.Fragment>
                ) : (
                  <React.Fragment>
                    <Error error={purchaseLabelError} />
                    <p>
                      <Button
                        disabled={purchasingLabel}
                        onClick={onPurchaseLabel}
                      >
                        Purchase label
                      </Button>
                    </p>
                  </React.Fragment>
                )}
              </Card.Body>
            </Card>
          </Col>
          <Col xs={12} lg={6}>
            <Receipt receipt={data.transaction} />
          </Col>
        </Row>
      )}
    </Container>
  );
}
