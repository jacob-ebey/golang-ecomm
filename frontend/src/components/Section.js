import React from "react";
import { Link } from "react-router-dom";

import Button from "react-bootstrap/Button";
import ButtonToolbar from "react-bootstrap/ButtonToolbar";
import Col from "react-bootstrap/Col";
import Container from "react-bootstrap/Container";
import Row from "react-bootstrap/Row";
import Jumbotron from "react-bootstrap/Jumbotron";

import Markdown from "./Markdown";

export default function Section({ section }) {
  return (
    <React.Fragment>
      <Container>
        {section.rows &&
          section.rows.map((row, j) => (
            <Row key={j} className="mb-4 mt-4">
              {row.columns &&
                row.columns.map(
                  (col, k) =>
                    col && (
                      <Col
                        key={k}
                        className="d-flex flex-column mb-4 mt-4"
                        style={{
                          ...(col.align ? { textAlign: col.align } : {})
                        }}
                        {...col.breakpoints}
                      >
                        {col.content && <Markdown source={col.content} />}
                        {col.Component && <col.Component col={col} />}
                        {col.callToAction && (
                          <ButtonToolbar
                            style={{
                              marginTop: "auto",
                              justifyContent: col.callToAction.align
                            }}
                          >
                            <Button
                              as={Link}
                              to={col.callToAction.to}
                              variant="primary"
                              size={col.callToAction.size}
                            >
                              {col.callToAction.label}
                            </Button>
                          </ButtonToolbar>
                        )}
                      </Col>
                    )
                )}
            </Row>
          ))}
      </Container>
      {section.spacerImage && (
        <Jumbotron
          fluid={true}
          style={{
            minHeight: "40vh",
            backgroundAttachment: "fixed",
            backgroundPosition: "center",
            backgroundRepeat: "no-repeat",
            backgroundSize: "cover",
            backgroundImage: `url(${section.spacerImage})`,
            marginBottom: 0
          }}
        />
      )}
    </React.Fragment>
  );
}
