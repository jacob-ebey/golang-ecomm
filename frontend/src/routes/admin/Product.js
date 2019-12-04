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
  query ProductProduct($id: Int!) {
    product(id: $id) {
      id
      name
      description
      details
      published
      images {
        id
        name
        raw
        thumbnail
        height600
      }
      options {
        id
        label
        values {
          id
          value
        }
      }
      variants {
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
      }
    }
    UpdateProductInput: __type(name: "UpdateProductInput") {
      ...GraphQLFormType
    }
    CreateProductOptionInput: __type(name: "CreateProductOptionInput") {
      ...GraphQLFormType
    }
    CreateProductOptionValueInput: __type(
      name: "CreateProductOptionValueInput"
    ) {
      ...GraphQLFormType
    }
    CreateProductPermutationsInput: __type(
      name: "CreateProductPermutationsInput"
    ) {
      ...GraphQLFormType
    }
    UpdateProductVariantInput: __type(name: "UpdateProductVariantInput") {
      ...GraphQLFormType
    }
  }
  ${GraphQLForm.fragments.type}
`;

const ADD_PRODUCT_IMAGE = gql`
  mutation ProductAddPrductImage($id: Int!, $image: Upload!) {
    addProductImage(id: $id, image: $image) {
      id
      name
      raw
      thumbnail
      height600
    }
  }
`;

const REMOVE_IMAGE = gql`
  mutation ProductRemoveImage($productId: Int!, $imageId: Int!) {
    removeProductImage(id: $productId, imageId: $imageId)
  }
`;

const UPDATE_PRODUCT = gql`
  mutation ProductUpdateProduct($id: Int!, $product: UpdateProductInput!) {
    updateProduct(id: $id, product: $product) {
      id
      name
      description
      details
    }
  }
`;

const PUBLISH_PRODUCT = gql`
  mutation ProductPublishProduct($id: Int!, $published: Boolean!) {
    publishProduct(published: $published, id: $id) {
      id
      published
    }
  }
`;

const CREATE_PRODUCT_OPTION = gql`
  mutation ProductCreateProductOption(
    $productId: Int!
    $option: CreateProductOptionInput!
  ) {
    createProductOption(productId: $productId, option: $option) {
      id
      label
      values {
        id
        value
      }
    }
  }
`;

const REMOVE_PRODUCT_OPTION = gql`
  mutation ProductRemoveProductOption(
    $productId: Int!
    $productOptionId: Int!
  ) {
    removeProductOption(
      productId: $productId
      productOptionId: $productOptionId
    ) {
      id
    }
  }
`;

const REMOVE_VARIANT = gql`
  mutation ProductRemoveVariant($id: Int!) {
    removeProductVariant(id: $id) {
      id
    }
  }
`;

const CREATE_PERMUTATIONS = gql`
  mutation ProductCreatePermutations(
    $productId: Int!
    $input: CreateProductPermutationsInput!
  ) {
    createProductPermutations(productId: $productId, input: $input) {
      id
    }
  }
`;

export default function Product() {
  const [imageToDelete, setImageToDelete] = React.useState(false);

  const closeDeleteImage = () => setImageToDelete(false);
  const showDeleteImage = image => () => setImageToDelete(image);

  const { id } = useParams();

  const { data, error, loading } = useQuery(QUERY, { variables: { id } });

  const [
    addImage,
    { error: addImageError, loading: addingImage }
  ] = useMutation(ADD_PRODUCT_IMAGE, {
    refetchQueries: [
      {
        query: QUERY,
        variables: { id }
      }
    ]
  });

  const [
    publishProduct,
    { error: publishProductError, loading: publishingProduct }
  ] = useMutation(PUBLISH_PRODUCT, {
    variables: {
      id,
      published: data && data.product && !data.product.published
    },
    refetchQueries: [
      {
        query: QUERY,
        variables: { id }
      }
    ]
  });

  const [productToUpdate, setProductToUpdate] = React.useState(null);
  const productToUpdateValues = React.useMemo(() => {
    if (!!productToUpdate) {
      return productToUpdate;
    }
    if (data && data.product) {
      const { name, description, details } = data.product;
      return {
        name,
        description,
        details
      };
    }
  }, [productToUpdate, data]);
  const [
    updateProduct,
    { error: updateProductError, loading: updatingProduct }
  ] = useMutation(UPDATE_PRODUCT, {
    variables: {
      id,
      product: productToUpdate
    }
  });
  const handleUpdateProduct = React.useCallback(
    event => {
      event.preventDefault();

      updateProduct();
    },
    [updateProduct]
  );

  const [
    removeProductImage,
    { error: removeProductImageError, loading: removingProductImage }
  ] = useMutation(REMOVE_IMAGE, {
    variables: {
      productId: id,
      imageId: imageToDelete && imageToDelete.id
    },
    refetchQueries: [
      {
        query: QUERY,
        variables: { id }
      }
    ]
  });

  const [variantToRemove, setVariantToRemove] = React.useState(false);
  const onRemoveProductVariant = React.useCallback(
    variant => () => setVariantToRemove(variant),
    [setVariantToRemove]
  );
  const closeRemoveVariant = React.useCallback(
    () => setVariantToRemove(false),
    [setVariantToRemove]
  );
  const [
    removeProductVariant,
    { error: removeProductVariantError, loading: removingProductVariant }
  ] = useMutation(REMOVE_VARIANT, {
    variables: {
      id: variantToRemove && variantToRemove.id
    },
    refetchQueries: [
      {
        query: QUERY,
        variables: { id }
      }
    ],
    onCompleted: data => {
      closeRemoveVariant();
    }
  });
  const handleRemoveProductVariant = React.useCallback(() => {
    removeProductVariant();
  }, [removeProductVariant]);

  const [
    productPermutationsInput,
    setProductPermutationsInput
  ] = React.useState({});
  const [
    createPermutations,
    { error: createPermutationsError, loading: creatingPermutations }
  ] = useMutation(CREATE_PERMUTATIONS, {
    variables: {
      productId: id,
      input: productPermutationsInput
    },
    refetchQueries: [
      {
        query: QUERY,
        variables: { id }
      }
    ]
  });

  const handleCreateProductPermutations = React.useCallback(() => {
    if (creatingPermutations) {
      return;
    }

    createPermutations();
  }, [createPermutations, creatingPermutations]);

  const [newOption, setNewOption] = React.useState({});
  const [
    createProductOption,
    { error: createProductOptionError, loading: creatingProductOption }
  ] = useMutation(CREATE_PRODUCT_OPTION, {
    variables: {
      productId: id,
      option: newOption
    },
    refetchQueries: [
      {
        query: QUERY,
        variables: { id }
      }
    ],
    onCompleted: data => {
      setNewOption({});
    }
  });
  const handleNewOption = React.useCallback(
    event => {
      event.preventDefault();
      createProductOption();
    },
    [newOption]
  );

  const [optionToRemove, setOptionToRemove] = React.useState(false);
  const closeRemoveOption = React.useCallback(() => setOptionToRemove(false), [
    setOptionToRemove
  ]);
  const [
    removeProductOption,
    { error: removeProductOptionError, loading: removingProductOption }
  ] = useMutation(REMOVE_PRODUCT_OPTION, {
    variables: {
      productId: id,
      productOptionId: optionToRemove && optionToRemove.id
    },
    refetchQueries: [
      {
        query: QUERY,
        variables: { id }
      }
    ],
    onCompleted: () => {
      closeRemoveOption();
    }
  });
  const onRemoveOption = React.useCallback(
    option => () => setOptionToRemove(option),
    [setOptionToRemove]
  );
  const handleRemoveOption = React.useCallback(() => {
    removeProductOption();
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

  const onPublish = React.useCallback(() => publishProduct(), [publishProduct]);

  const onRemoveImage = React.useCallback(
    () => removeProductImage().then(() => closeDeleteImage()),
    [removeProductImage, closeDeleteImage]
  );

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
            <Link className="breadcrumb-item" to="/admin/products">
              Products
            </Link>
            <Breadcrumb.Item active>#{id}</Breadcrumb.Item>
          </Breadcrumb>
          <Error error={error} />
          <Error error={addImageError} />
        </Col>
      </Row>

      {!loading && data && data.product && (
        <React.Fragment>
          <Row>
            <Col xs={12}>
              <ButtonToolbar className="justify-content-end">
                {data.product.published ? (
                  <Button
                    disabled={publishingProduct}
                    onClick={onPublish}
                    variant="danger"
                  >
                    Recall
                  </Button>
                ) : (
                  <Button
                    disabled={publishingProduct}
                    onClick={onPublish}
                    variant="primary"
                  >
                    Publish
                  </Button>
                )}
              </ButtonToolbar>
              <Error error={publishProductError} />

              <h2>Overview</h2>
              <Error error={updateProductError} />
              <GraphQLForm
                name={null}
                type={data.UpdateProductInput}
                types={{}}
                order={["name", "description", "details"]}
                values={productToUpdateValues}
                onChange={setProductToUpdate}
                onSubmit={handleUpdateProduct}
              >
                <Button
                  type="submit"
                  disabled={updatingProduct || !productToUpdate}
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

            {data.product.images &&
              data.product.images.map(image => (
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
            {addingImage && (
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

          <Row>
            <Col className="mt-4 mb-4" xs={12}>
              <h2>Options</h2>

              <Accordion>
                <Card>
                  <Card.Header>
                    <Accordion.Toggle as={Button} variant="link" eventKey="0">
                      New Option
                    </Accordion.Toggle>
                  </Card.Header>
                  <Accordion.Collapse eventKey="0">
                    <Card.Body>
                      <Error error={createProductOptionError} />
                      <GraphQLForm
                        name={null}
                        type={data.CreateProductOptionInput}
                        types={{
                          CreateProductOptionValueInput:
                            data.CreateProductOptionValueInput
                        }}
                        order={["label", "values"]}
                        disabled={creatingProductOption}
                        values={newOption}
                        onChange={setNewOption}
                        onSubmit={handleNewOption}
                      >
                        <Button type="submit" disabled={creatingProductOption}>
                          Create Option
                        </Button>
                      </GraphQLForm>
                    </Card.Body>
                  </Accordion.Collapse>
                </Card>
                <Card>
                  <Card.Header>
                    <Accordion.Toggle as={Button} variant="link" eventKey="1">
                      View Options
                    </Accordion.Toggle>
                  </Card.Header>
                  <Accordion.Collapse eventKey="1">
                    <Card.Body>
                      {data.product.options && data.product.options.length > 0
                        ? data.product.options.map(option => (
                            <Card className="mb-4" key={option.id}>
                              <Card.Body>
                                <div className="d-flex">
                                  <div className="flex-grow-1">
                                    <GraphQLForm
                                      name={null}
                                      disabled={true}
                                      type={data.CreateProductOptionInput}
                                      types={{
                                        CreateProductOptionValueInput:
                                          data.CreateProductOptionValueInput
                                      }}
                                      order={["label", "values"]}
                                      values={option}
                                    ></GraphQLForm>
                                  </div>
                                  <div className="d-flex align-items-center ml-2">
                                    <Button
                                      size="sm"
                                      variant="danger"
                                      onClick={onRemoveOption(option)}
                                    >
                                      X
                                    </Button>
                                  </div>
                                </div>
                              </Card.Body>
                            </Card>
                          ))
                        : "No options found."}
                    </Card.Body>
                  </Accordion.Collapse>
                </Card>
              </Accordion>
            </Col>
          </Row>

          <Row>
            <Col xs={12}>
              <h2>Variants</h2>

              <Error error={createPermutationsError} />
              {(!data.product.variants ||
                data.product.variants.length === 0) && (
                <React.Fragment>
                  <GraphQLForm
                    name={null}
                    type={data.CreateProductPermutationsInput}
                    types={{}}
                    order={["price", "width", "length", "height", "weight"]}
                    values={productPermutationsInput}
                    onChange={setProductPermutationsInput}
                    onSubmit={handleCreateProductPermutations}
                  >
                    <Button type="submit" disabled={creatingPermutations}>
                      Create all permutations.
                    </Button>
                  </GraphQLForm>
                </React.Fragment>
              )}
            </Col>
          </Row>
          <Row>
            {data.product.variants &&
              data.product.variants.map(variant => (
                <Col key={variant.id} className="mt-4" xs={12} md={6} lg={4}>
                  <Card as={Link} to={`/admin/variants/${variant.id}`}>
                    <Card.Body>
                      <ButtonToolbar className="justify-content-end">
                        <Button
                          size="sm"
                          variant="danger"
                          onClick={onRemoveProductVariant(variant)}
                        >
                          X
                        </Button>
                      </ButtonToolbar>

                      <p>
                        <strong>Selected Options</strong> <br />
                        {variant.selectedOptions &&
                          variant.selectedOptions
                            .map(option => option.value)
                            .join(", ")}
                      </p>

                      <GraphQLForm
                        disabled={true}
                        name={null}
                        type={data.UpdateProductVariantInput}
                        types={{}}
                        order={[
                          "name",
                          "price",
                          "width",
                          "length",
                          "height",
                          "weight"
                        ]}
                        values={variant}
                      />
                    </Card.Body>
                  </Card>
                </Col>
              ))}
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

      <Modal show={!!optionToRemove} onHide={closeRemoveOption} size="xl">
        <Modal.Header closeButton>
          <Modal.Title>Remove Option</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Error error={removeProductOptionError} />
          <p>Do you really wish to remove this option?</p>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={closeRemoveOption}>
            Cancel
          </Button>
          <Button
            variant="danger"
            disabled={removingProductOption}
            onClick={handleRemoveOption}
          >
            Remove
          </Button>
        </Modal.Footer>
      </Modal>

      <Modal show={!!variantToRemove} onHide={closeRemoveVariant} size="xl">
        <Modal.Header closeButton>
          <Modal.Title>Remove Variant</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Error error={removeProductVariantError} />
          <p>Do you really wish to remove this variant?</p>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={closeRemoveVariant}>
            Cancel
          </Button>
          <Button
            variant="danger"
            disabled={removingProductVariant}
            onClick={handleRemoveProductVariant}
          >
            Remove
          </Button>
        </Modal.Footer>
      </Modal>

      <Modal show={!!imageToDelete} onHide={closeDeleteImage} size="xl">
        <Modal.Header closeButton>
          <Modal.Title>Delete Image</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Error error={removeProductImageError} />
          <p>Do you really wish to delete this image?</p>
          <Image src={imageToDelete.height600} fluid />
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={closeDeleteImage}>
            Cancel
          </Button>
          <Button
            variant="danger"
            disabled={removingProductImage}
            onClick={onRemoveImage}
          >
            Delete
          </Button>
        </Modal.Footer>
      </Modal>
    </Container>
  );
}
