import React from "react";
import { useQueryParam, NumberParam } from "use-query-params";
import { Link } from "react-router-dom";
import currency from "currency.js";

import { useQuery } from "@apollo/react-hooks";
import { gql } from "apollo-boost";

import Breadcrumb from "react-bootstrap/Breadcrumb";
import Col from "react-bootstrap/Col";
import Container from "react-bootstrap/Container";
import Pagination from "react-bootstrap/Pagination";
import Row from "react-bootstrap/Row";
import Table from "react-bootstrap/Table";

import Error from "../../components/Error";

const QUERY = gql`
  query AdminTransactions($skip: Int, $limit: Int) {
    transactions(skip: $skip, limit: $limit) {
      id
      subtotal
      taxes
      shipping
      total
    }
  }
`;

function toLink(item) {
  return `/admin/transactions/${item.id}`;
}

export default function Transactions() {
  const [skipParam, setSkip] = useQueryParam("skip", NumberParam);
  const [limitParam] = useQueryParam("limit", NumberParam);
  const skip = skipParam > 0 ? skipParam : 0;
  const limit = limitParam > 0 ? limitParam : 20;
  const perviousPage = React.useCallback(() => {
    const next = skip - limit;
    setSkip(next < 0 ? 0 : next, "pushIn");
  });
  const nextPage = React.useCallback(() => setSkip(skip + limit, "pushIn"));

  const { data, error, loading } = useQuery(QUERY, {
    variables: { skip, limit }
  });

  return (
    <Container>
      <Row>
        <Col className="mt-4 mb-4" xs={12}>
          <Breadcrumb>
            <Breadcrumb.Item active>Admin</Breadcrumb.Item>
            <Breadcrumb.Item active>Transactions</Breadcrumb.Item>
          </Breadcrumb>
          <Error error={error} />
        </Col>
      </Row>
      <Row>
        <Table responsive>
          <thead>
            <tr>
              <th>#</th>
              <th>Subtotal</th>
              <th>Taxes</th>
              <th>Shipping</th>
              <th>Total</th>
            </tr>
          </thead>
          <tbody>
            {!loading &&
              data &&
              data.transactions &&
              data.transactions.map(transaction => (
                <tr key={transaction.id}>
                  <td>
                    <Link to={toLink(transaction)}>{transaction.id}</Link>
                  </td>
                  <td>
                    <Link to={toLink(transaction)}>
                      ${currency(transaction.subtotal / 100).format()}
                    </Link>
                  </td>
                  <td>
                    <Link to={toLink(transaction)}>
                      ${currency(transaction.taxes / 100).format()}
                    </Link>
                  </td>
                  <td>
                    <Link to={toLink(transaction)}>
                      ${currency(transaction.shipping / 100).format()}
                    </Link>
                  </td>
                  <td>
                    <Link to={toLink(transaction)}>
                      ${currency(transaction.total / 100).format()}
                    </Link>
                  </td>
                </tr>
              ))}
          </tbody>
        </Table>
      </Row>
      <Row>
        {!loading && data && data.transactions && (
          <Pagination>
            {skip > 0 && (
              <Pagination.Prev as={"button"} onClick={perviousPage} />
            )}
            {data.transactions.length === limit && (
              <Pagination.Next as={"button"} onClick={nextPage} />
            )}
          </Pagination>
        )}
      </Row>
    </Container>
  );
}
