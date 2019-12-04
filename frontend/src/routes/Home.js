import React from "react";

import Jumbotron from "react-bootstrap/Jumbotron";

import { home } from "../config";
import Section from "../components/Section";

function Home() {
  return (
    <React.Fragment>
      {home.hero && (
        <Jumbotron
          fluid={true}
          style={{
            display: "flex",
            flexDirection: "column",
            justifyContent: "center",
            minHeight: home.heroImage && "90vh",

            backgroundAttachment: "fixed",
            backgroundPosition: "center",
            backgroundRepeat: "no-repeat",
            backgroundSize: "cover",
            ...(home.heroImage
              ? { backgroundImage: `url(${home.heroImage})` }
              : {})
          }}
        >
          <div
            style={
              home.heroImage && {
                backgroundColor: "rgba(0,0,0,0.4)",
                color: "white"
              }
            }
          >
            <Section section={home.hero} />
          </div>
        </Jumbotron>
      )}

      {home.sections &&
        home.sections.map((section, i) => (
          <Section key={i} section={section} />
        ))}
    </React.Fragment>
  );
}

export default Home;
