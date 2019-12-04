import React from "react";
import { useParams } from "react-router-dom";
import { gql } from "apollo-boost";
import { useMutation, useQuery } from "@apollo/react-hooks";

import Col from "react-bootstrap/Col";
import Container from "react-bootstrap/Container";
import Jumbotron from "react-bootstrap/Jumbotron";
import Row from "react-bootstrap/Row";

import { product } from "../config";
import AddToCart from "../components/AddToCart";
import Error from "../components/Error";
import Images from "../components/Images";
import Markdown from "../components/Markdown";
import PriceRange from "../components/PriceRange";
import Section from "../components/Section";
import ProductOptions from "../components/ProductOptions";

const QUERY = gql`
  query Product($slug: String!) {
    product: productBySlug(slug: $slug) {
      id
      name
      description
      details
      images {
        id
        name
        height600
      }
      priceRange {
        ...PriceRangePriceRange
      }
      options {
        ...ProductOptionsOptions
      }
    }
  }
  ${PriceRange.fragments.priceRange}
  ${ProductOptions.fragments.options}
`;

const VARIANT_BY_OPTIONS = gql`
  query ShopVariantByOptions($productId: Int!, $selectedOptions: [Int!]!) {
    variant: productVariantBySelectedOptions(
      productId: $productId
      selectedOptions: $selectedOptions
    ) {
      id
      price
    }
  }
`;

const ADD_TO_CART = gql`
  mutation ShopAddToCart($variantId: Int!, $quantity: Int!) {
    addToCart: changeCartQuantity(variantId: $variantId, quantity: $quantity)
      @client
  }
`;

function Product() {
  const [selectedOptions, setSelectedOptions] = React.useState([]);

  const { slug } = useParams();

  const { data, error } = useQuery(QUERY, {
    variables: { slug }
  });

  const images = React.useMemo(
    () =>
      data &&
      data.product &&
      data.product.images &&
      data.product.images.map(image => ({
        url: image.height600,
        alt: image.name
      })),
    [data]
  );

  const {
    data: variantData,
    loading: variantLoading,
    error: variantError
  } = useQuery(VARIANT_BY_OPTIONS, {
    variables: {
      productId: data && data.product && data.product.id,
      selectedOptions
    },
    skip: !data || !data.product
  });

  const [addToCart, { error: addToCartError }] = useMutation(ADD_TO_CART);

  const priceRange = React.useMemo(() =>
    variantData && variantData.variant
      ? { min: variantData.variant.price, max: variantData.variant.price }
      : data && data.product && data.product.priceRange
  );

  const onAddToCart = React.useCallback(
    quantity => {
      if (
        variantError ||
        variantLoading ||
        !variantData ||
        !variantData.variant
      ) {
        return;
      }

      addToCart({ variables: { variantId: variantData.variant.id, quantity } });
    },
    [data, variantData, variantError, variantLoading, addToCart]
  );

  return (
    <React.Fragment>
      {product.heroImage && (
        <Jumbotron
          fluid={true}
          style={{
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            minHeight: product.heroImage && "40vh",

            backgroundAttachment: "fixed",
            backgroundPosition: "center",
            backgroundRepeat: "no-repeat",
            backgroundSize: "cover",
            ...(product.heroImage
              ? { backgroundImage: `url(${product.heroImage})` }
              : {})
          }}
        >
          {data && data.product && (
            <div
              style={
                product.heroImage && {
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
                          content: `
# ${data.product.name}

${data.product.description}
`
                        }
                      ]
                    }
                  ]
                }}
              />
            </div>
          )}
        </Jumbotron>
      )}

      <Error error={error} />
      <Error error={addToCartError} />
      <Error error={variantError} />

      {data && data.product && (
        <Container>
          <Row>
            <Col xs={12} sm={12} md={6} lg={6} style={{ height: "80vh" }}>
              <Images images={images} />
            </Col>
            <Col xs={12} sm={12} md={6} lg={6}>
              <h1>{data.product.name}</h1>
              <PriceRange as="h2" priceRange={priceRange} />

              {data.product.details && (
                <Markdown source={data.product.details} />
              )}

              <ProductOptions
                options={data.product.options}
                onChange={setSelectedOptions}
              />

              <AddToCart
                disabled={
                  variantLoading || !variantData || !variantData.variant
                }
                onAddToCart={onAddToCart}
              />
            </Col>
          </Row>
        </Container>
      )}
    </React.Fragment>
  );
}

export default Product;
