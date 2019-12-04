import React from "react";
import { gql } from "apollo-boost";

import Card from "react-bootstrap/Card";

export default function ProductTile({ product, ...rest }) {
  return (
    <Card style={{ marginBottom: "2rem" }} {...rest}>
      <Card.Img
        variant="top"
        alt={product.images && product.images[0] && product.images[0].name}
        src={product.images && product.images[0] && product.images[0].thumbnail}
      />
      <Card.Body>
        <Card.Title>{product.name}</Card.Title>
        <Card.Text>{product.description}</Card.Text>
      </Card.Body>
    </Card>
  );
}

ProductTile.fragments = {
  product: gql`
    fragment ProductTileProduct on Product {
      id
      slug
      name
      description
      images {
        id
        name
        thumbnail
      }
    }
  `
};
