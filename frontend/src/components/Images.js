import React from "react";

import Carousel from "react-bootstrap/Carousel";
import Image from "react-bootstrap/Image";

export default function Images({ images }) {
  return images && images.length > 0 ? (
    <Carousel slide={true} touch={true}>
      {images.map((image, key) => (
        <Carousel.Item key={key}>
          <Image
            fluid
            style={{ height: "60vh", objectFit: "contain" }}
            alt={image.alt || "Unknown"}
            src={image.url}
          />
        </Carousel.Item>
      ))}
    </Carousel>
  ) : null;
}
