import React from "react";
import { gql } from "apollo-boost";

import Button from "react-bootstrap/Button";
import Card from "react-bootstrap/Card";
import Form from "react-bootstrap/Form";

import MdEditor from "react-markdown-editor-lite";
import MarkdownIt from "markdown-it";

const mdParser = new MarkdownIt();

export default function GraphQLForm({
  as: htmlAs = Form,
  name,
  type,
  types = {},
  order,
  values: suppliedValues,
  onChange: suppliedOnChange,
  children,
  disabled,
  ...rest
}) {
  const Component = React.useMemo(() => htmlAs || Form, [htmlAs]);
  const [internalValues, setInternalValues] = React.useState(
    suppliedValues || {}
  );

  const values = React.useMemo(() => suppliedValues || internalValues, [
    suppliedValues,
    internalValues
  ]);
  const onChange = React.useMemo(() => suppliedOnChange || setInternalValues, [
    suppliedOnChange,
    setInternalValues
  ]);

  const handleOnChange = React.useCallback(
    name => newValue => {
      onChange({
        ...values,
        [name]: newValue
      });
    },
    [values, onChange]
  );

  const inputFields = React.useMemo(() => {
    if (!type.inputFields) {
      return [];
    }

    if (!order) {
      return type.inputFields.sort((a, b) => a.name.localeCompare(b.name));
    }

    return type.inputFields.sort((a, b) =>
      order.indexOf(a.name) > order.indexOf(b.name) ? 1 : -1
    );
  }, [type, order]);

  return (
    <Component disabled={disabled} {...rest}>
      {(typeof name === "undefined" || !!name) && <h2>{name || type.name}</h2>}
      {inputFields.map(field => (
        <GraphQLFormField
          key={field.name}
          types={types}
          field={field}
          disabled={disabled}
          value={values[field.name]}
          onChange={handleOnChange(field.name)}
        />
      ))}
      {children}
    </Component>
  );
}

function GraphQLFormField({ types, field, disabled, value, onChange }) {
  const label = React.useMemo(
    () => field.name.slice(0, 1).toUpperCase() + field.name.slice(1),
    [field]
  );

  return (
    <Form.Group controlId="formBasicEmail">
      <Form.Label>{label}</Form.Label>
      <GraphQLFormControl
        types={types}
        field={field}
        type={field.type}
        disabled={disabled}
        value={value}
        onChange={onChange}
      />
    </Form.Group>
  );
}

function getBaseType(type) {
  if (type.kind !== "LIST" && !!type.ofType) {
    return getBaseType(type.ofType);
  }

  return type;
}

function GraphQLFormControl({
  types,
  field,
  type: suppliedType,
  disabled,
  value,
  onChange
}) {
  const required = React.useMemo(() => suppliedType.kind === "NON_NULL", [
    suppliedType
  ]);
  const type = React.useMemo(() => getBaseType(suppliedType), [suppliedType]);

  return React.useMemo(() => {
    switch (type.kind) {
      case "SCALAR":
        switch (type.name) {
          case "Markdown":
            return (
              <GraphQLFormControlMarkdownScalar
                field={field}
                disabled={disabled}
                value={value}
                onChange={onChange}
                type={type}
                required={required}
              />
            );
          default:
            return (
              <GraphQLFormControlScalar
                field={field}
                disabled={disabled}
                value={value}
                onChange={onChange}
                type={type}
                required={required}
              />
            );
        }
      case "LIST":
        return (
          <GraphQLFormControlList
            types={types}
            field={field}
            disabled={disabled}
            value={value}
            onChange={onChange}
            type={type}
            required={required}
          />
        );
      case "INPUT_OBJECT":
        return (
          <GraphQLForm
            as="div"
            types={types}
            disabled={disabled}
            name={null}
            type={types[type.name]}
            values={value}
            onChange={onChange}
          />
        );
      default:
        return null;
    }
  }, [field, type, required, value, types, onChange]);
}

function GraphQLFormControlList({
  types,
  field,
  disabled,
  value: suppliedValue,
  onChange,
  type,
  required
}) {
  const values = React.useMemo(() => {
    if (Array.isArray(suppliedValue)) {
      return suppliedValue;
    }

    return [null];
  }, [suppliedValue]);

  const handleAddNew = React.useCallback(() => {
    onChange([...values, null]);
  }, [onChange, values]);

  const handleRemove = React.useCallback(
    index => () => {
      onChange([...values.slice(0, index), ...values.slice(index + 1)]);
    },
    [onChange, values]
  );

  const handleOnChange = React.useCallback(
    index => newValue => {
      onChange(Object.assign([], values, { [index]: newValue }));
    },
    [onChange, values]
  );

  const children = React.useMemo(
    () =>
      values.map((value, index) => (
        <Card key={index} className="mb-4">
          <Card.Body>
            <div className="d-flex">
              <div className="flex-grow-1">
                <GraphQLFormControl
                  types={types}
                  disabled={disabled}
                  field={field}
                  type={type.ofType}
                  value={value}
                  onChange={handleOnChange(index)}
                />
              </div>
              {!disabled && (
                <div className="d-flex align-items-center ml-2">
                  <Button
                    size="sm"
                    variant="danger"
                    onClick={handleRemove(index)}
                  >
                    X
                  </Button>
                </div>
              )}
            </div>
          </Card.Body>
        </Card>
      )),
    [types, field, disabled, type, handleOnChange, handleRemove, values]
  );

  return React.useMemo(
    () => (
      <React.Fragment>
        {" "}
        {!disabled && (
          <Button size="sm" onClick={handleAddNew}>
            +
          </Button>
        )}
        {field.description && <p className="text-muted">{field.description}</p>}
        <div className="ml-2">{children}</div>
      </React.Fragment>
    ),
    [field, handleAddNew, disabled, children]
  );
}

function GraphQLFormControlMarkdownScalar({
  field,
  disabled,
  value: suppliedValue,
  onChange,
  type,
  required
}) {
  const value = React.useMemo(() => {
    if (typeof suppliedValue === "string") {
      return suppliedValue;
    }

    return "";
  }, [suppliedValue]);

  const handleEditorChange = React.useCallback(
    ({ text }) => {
      onChange(text);
    },
    [onChange]
  );

  return (
    <div style={{ height: "60vh" }}>
      <MdEditor
        disabled={disabled}
        value={value}
        renderHTML={text => mdParser.render(text)}
        onChange={handleEditorChange}
      />
    </div>
  );
}

function GraphQLFormControlScalar({
  field,
  disabled,
  value: suppliedValue,
  onChange,
  type,
  required
}) {
  const value = React.useMemo(() => {
    if (typeof suppliedValue !== "undefined") {
      return suppliedValue;
    }

    switch (type.name) {
      case "Int":
        return 0;
      case "Float":
        return 0.0;
      default:
        return "";
    }
  }, [suppliedValue]);

  const inputType = React.useMemo(() => {
    switch (type.name) {
      case "Int":
      case "Float":
        return "number";
      default:
        return "text";
    }
  }, [type]);

  const handleOnChange = React.useCallback(
    event => onChange(event.target.value),
    [onChange]
  );

  return React.useMemo(
    () => (
      <React.Fragment>
        <Form.Control
          disabled={disabled}
          type={inputType}
          placeholder={type.name}
          required={required}
          value={value}
          onChange={handleOnChange}
        />
        {field.description && (
          <Form.Text className="text-muted">{field.description}</Form.Text>
        )}
      </React.Fragment>
    ),
    [field, value, inputType, type, handleOnChange]
  );
}

GraphQLForm.fragments = {
  type: gql`
    fragment GraphQLFormType on __Type {
      name
      inputFields {
        ...GraphQLFormInputValue
      }
    }
    fragment GraphQLFormInputValue on __InputValue {
      name
      description
      type {
        ...GraphQLFormTypeRef
      }
      defaultValue
    }
    fragment GraphQLFormTypeRef on __Type {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
            }
          }
        }
      }
    }
  `
};
