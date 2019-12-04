import React from "react";
import { useQueryParam, NumberParam } from "use-query-params";
import { Link } from "react-router-dom";

import { useQuery } from "@apollo/react-hooks";
import { gql } from "apollo-boost";

import Breadcrumb from "react-bootstrap/Breadcrumb";
import Col from "react-bootstrap/Col";
import Container from "react-bootstrap/Container";
import Pagination from "react-bootstrap/Pagination";
import Row from "react-bootstrap/Row";
import Table from "react-bootstrap/Table";

import Error from "../../components/Error";
import PriceRange from "../../components/PriceRange";

const QUERY = gql`
  query AdminProducts($skip: Int, $limit: Int) {
    products(skip: $skip, limit: $limit) {
      id
      name
      description
      priceRange {
        ...PriceRangePriceRange
      }
    }
  }
  ${PriceRange.fragments.priceRange}
`;

function toLink(item) {
  return `/admin/products/${item.id}`;
}

export default function Products() {
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
            <Breadcrumb.Item active>Products</Breadcrumb.Item>
          </Breadcrumb>
          <Error error={error} />
        </Col>
      </Row>
      <Row>
        <Table responsive>
          <thead>
            <tr>
              <th>#</th>
              <th>Name</th>
              <th>Description</th>
              <th>Price</th>
            </tr>
          </thead>
          <tbody>
            {!loading &&
              data &&
              data.products &&
              data.products.map(product => (
                <tr key={product.id}>
                  <td>
                    <Link to={toLink(product)}>{product.id}</Link>
                  </td>
                  <td>
                    <Link to={toLink(product)}>{product.name}</Link>
                  </td>
                  <td>
                    <Link to={toLink(product)}>{product.description}</Link>
                  </td>
                  <td>
                    <PriceRange
                      as={Link}
                      to={toLink(product)}
                      priceRange={product.priceRange}
                    />
                  </td>
                </tr>
              ))}
          </tbody>
        </Table>
      </Row>
      <Row>
        {!loading && data && data.products && (
          <Pagination>
            {skip > 0 && (
              <Pagination.Prev as={"button"} onClick={perviousPage} />
            )}
            {data.products.length === limit && (
              <Pagination.Next as={"button"} onClick={nextPage} />
            )}
          </Pagination>
        )}
      </Row>
    </Container>
  );
}
