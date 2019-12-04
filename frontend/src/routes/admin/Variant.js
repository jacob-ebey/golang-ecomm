import React from "react";
import { Link, useParams } from "react-router-dom";

import { useMutation, useQuery } from "@apollo/react-hooks";
import { gql } from "apollo-boost";

import Accordion from "react-bootstrap/Accordion";
import Breadcrumb from "react-bootstrap/Breadcrumb";
import Button from "react-bootstrap/Button";
import ButtonToolbar from "react-bootstrap/ButtonToolbar";
import Card from "react-bootstrap/Card";
import Col from "react-bootstrap/Col";
import Container from "react-bootstrap/Container";
import Image from "react-bootstrap/Image";
import Modal from "react-bootstrap/Modal";
import Row from "react-bootstrap/Row";
import Spinner from "react-bootstrap/Spinner";

import Error from "../../components/Error";
import GraphQLForm from "../../components/GraphQLForm";

const QUERY = gql`
  query VariantAdmin($id: Int!) {
    variants: productVariantsByIds(variantIds: [$id]) {
      id
      name
      price
      length
      width
      height
      weight
      selectedOptions {
        id
        productOptionId
        value
      }
      images {
        id
        name
        raw
        thumbnail
        height600
      }
    }
    UpdateProductVariantInput: __type(name: "UpdateProductVariantInput") {
      ...GraphQLFormType
    }
  }
  ${GraphQLForm.fragments.type}
`;

const UPDATE_VARIANT = gql`
  mutation VariantUpdateVariant($id: Int!, $input: UpdateProductVariantInput!) {
    updateProductVariant(id: $id, input: $input) {
      id
    }
  }
`;

const ADD_VARIANT_IMAGE = gql`
  mutation VariantAddImage($id: Int!, $image: Upload!) {
    addProductVariantImage(id: $id, image: $image) {
      id
      name
      raw
      thumbnail
      height600
    }
  }
`;

export default function Variant() {
  const { id } = useParams();

  const { data, error, loading } = useQuery(QUERY, { variables: { id } });
  const variant = React.useMemo(
    () => data && data.variants && data.variants[0],
    [data]
  );

  const [variantToUpdate, setVariantToUpdate] = React.useState(null);
  const [
    updateVariant,
    { error: updateVariantError, loading: updatingVariant }
  ] = useMutation(UPDATE_VARIANT, {
    refetchQueries: [
      {
        query: QUERY,
        variables: { id }
      }
    ]
  });
  const handleUpdateVariant = React.useCallback(
    event => {
      event.preventDefault();

      const { name, price, width, length, height, weight } = variantToUpdate;

      updateVariant({
        variables: {
          id,
          input: {
            name,
            price,
            width,
            length,
            height,
            weight
          }
        }
      });
    },
    [updateVariant, id, variantToUpdate]
  );

  const [
    addImage,
    { error: addImageError, loading: addingImage }
  ] = useMutation(ADD_VARIANT_IMAGE, {
    refetchQueries: [
      {
        query: QUERY,
        variables: { id }
      }
    ]
  });

  const onImageChange = React.useCallback(
    ({
      target: {
        validity,
        files: [image]
      }
    }) => {
      validity.valid && addImage({ variables: { id, image } });
    },
    [addImage, id]
  );

  const [imageToDelete, setImageToDelete] = React.useState(false);
  const closeDeleteImage = () => setImageToDelete(false);
  const showDeleteImage = image => () => setImageToDelete(image);

  const [imageToShow, setImageToShow] = React.useState(false);
  const handleShowImage = React.useCallback(
    image => () => setImageToShow(image),
    [setImageToShow]
  );
  const closeShowImage = React.useCallback(() => setImageToShow(false), [
    setImageToShow
  ]);

  return (
    <Container>
      <Row>
        <Col className="mt-4 mb-4" xs={12}>
          <Breadcrumb>
            <Breadcrumb.Item active>Admin</Breadcrumb.Item>
            <Breadcrumb.Item active>Product Variants</Breadcrumb.Item>
            <Breadcrumb.Item active>#{id}</Breadcrumb.Item>
          </Breadcrumb>
          <Error error={error} />
        </Col>
      </Row>

      {!loading && variant && (
        <React.Fragment>
          <Row>
            <Col xs={12}>
              <h2>Overview</h2>
              <Error error={updateVariantError} />

              <p>
                <strong>Selected Options</strong> <br />
                {variant.selectedOptions &&
                  variant.selectedOptions
                    .map(option => option.value)
                    .join(", ")}
              </p>

              <GraphQLForm
                name={null}
                type={data.UpdateProductVariantInput}
                types={{}}
                order={["name", "price", "width", "length", "height", "weight"]}
                values={variantToUpdate || variant}
                onChange={setVariantToUpdate}
                onSubmit={handleUpdateVariant}
              >
                <Button
                  type="submit"
                  disabled={updatingVariant || !variantToUpdate}
                >
                  Update Overview
                </Button>
              </GraphQLForm>
            </Col>
          </Row>

          <Row>
            <Col className="mt-4 mb-4" xs={12}>
              <h2>Images</h2>

              <div className="custom-file">
                <input
                  type="file"
                  className="custom-file-input"
                  id="customFile"
                  onChange={onImageChange}
                  disabled={addingImage}
                />
                <label className="custom-file-label" htmlFor="customFile">
                  Upload Image
                </label>
              </div>
            </Col>

            {variant.images &&
              variant.images.map(image => (
                <Col
                  key={image.id}
                  xs={6}
                  md={4}
                  lg={3}
                  className="position-relative"
                >
                  <div className="position-relative">
                    <Image
                      src={image.thumbnail}
                      thumbnail
                      style={{ cursor: "pointer" }}
                      onClick={handleShowImage(image)}
                    />
                    <div
                      className="position-absolute"
                      style={{ width: "100%", top: 0 }}
                    >
                      <ButtonToolbar className="justify-content-end m-1">
                        <Button
                          variant="danger"
                          size="sm"
                          onClick={showDeleteImage(image)}
                        >
                          X
                        </Button>
                      </ButtonToolbar>
                    </div>
                  </div>
                  <p className="lead text-wrap text-break">{image.name}</p>
                </Col>
              ))}
            {false && (
              <Col xs={6} md={4} lg={3}>
                <div
                  className="position-relative w-100"
                  style={{ width: "100%" }}
                >
                  <div
                    className="img-thumbnail d-flex justify-content-center align-items-center"
                    style={{
                      position: "absolute",
                      width: "100%",
                      height: "100%"
                    }}
                  >
                    <Spinner animation="grow" />
                  </div>
                  <div style={{ paddingBottom: "100%" }} />
                </div>
              </Col>
            )}
          </Row>
        </React.Fragment>
      )}

      <Modal show={!!imageToShow} onHide={closeShowImage} size="xl">
        <Modal.Header closeButton>
          <Modal.Title className="text-wrap text-break">
            {imageToShow.name}
          </Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Image src={imageToShow.raw} fluid />
        </Modal.Body>
      </Modal>

      <Modal show={!!imageToDelete} onHide={closeDeleteImage} size="xl">
        <Modal.Header closeButton>
          <Modal.Title>Delete Image</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          {/* <Error error={removeProductImageError} /> */}
          <p>Do you really wish to delete this image?</p>
          <Image src={imageToDelete.height600} fluid />
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={closeDeleteImage}>
            Cancel
          </Button>
          <Button
            variant="danger"
            // disabled={removingProductImage}
            // onClick={onRemoveImage}
          >
            Delete
          </Button>
        </Modal.Footer>
      </Modal>
    </Container>
  );
}
