import React from "react";
import { Link } from "react-router-dom";
import { useQueryParam, NumberParam } from "use-query-params";

import { useQuery } from "@apollo/react-hooks";
import { gql } from "apollo-boost";

import Col from "react-bootstrap/Col";
import Container from "react-bootstrap/Container";
import Pagination from "react-bootstrap/Pagination";
import Row from "react-bootstrap/Row";
import Jumbotron from "react-bootstrap/Jumbotron";

import { shop } from "../config";
import Error from "../components/Error";
import ProductTile from "../components/ProductTile";
import Section from "../components/Section";

const QUERY = gql`
  query Shop($skip: Int, $limit: Int) {
    catalog(skip: $skip, limit: $limit) {
      id
      ...ProductTileProduct
    }
  }
  ${ProductTile.fragments.product}
`;

function Shop() {
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
    <React.Fragment>
      {shop.hero && (
        <Jumbotron
          fluid={true}
          style={{
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            minHeight: shop.heroImage && "40vh",

            backgroundAttachment: "fixed",
            backgroundPosition: "center",
            backgroundRepeat: "no-repeat",
            backgroundSize: "cover",
            ...(shop.heroImage
              ? { backgroundImage: `url(${shop.heroImage})` }
              : {})
          }}
        >
          <div
            style={
              shop.heroImage && {
                backgroundColor: "rgba(0,0,0,0.4)",
                color: "white"
              }
            }
          >
            <Section section={shop.hero} />
          </div>
        </Jumbotron>
      )}
      {!loading && data && data.catalog ? (
        <Container>
          <Row>
            {data.catalog.map(product => (
              <Col key={product.id} xs={12} sm={12} md={6} lg={4}>
                <ProductTile
                  as={Link}
                  to={`shop/${product.slug}`}
                  product={product}
                />
              </Col>
            ))}
          </Row>
          <Row>
          <Pagination>
            {skip > 0 && (
              <Pagination.Prev as={"button"} onClick={perviousPage} />
            )}
            {data.catalog.length === limit && (
              <Pagination.Next as={"button"} onClick={nextPage} />
            )}
          </Pagination>
      </Row>
        </Container>
      ) : (
        <Error error={error} />
      )}
    </React.Fragment>
  );
}

export default Shop;
