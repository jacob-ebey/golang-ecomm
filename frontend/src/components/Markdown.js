import React from "react";
import ReactMarkdown from "react-markdown";

const imageSizeRegex = /_33B2BF251EFD_([0-9]+x|x[0-9]+|[0-9]+x[0-9]+)$/;
const imagePreprocessor = source =>
  source.replace(
    /(!\[[^\]]*\]\([^)\s]+) =([0-9]+x|x[0-9]+|[0-9]+x[0-9]+)\)/g,
    "$1_33B2BF251EFD_$2)"
  );

function imageRenderer({ src, ...props }) {
  const match = imageSizeRegex.exec(src);

  if (!match) {
    return <img alt="Unknown" src={src} {...props} />;
  }

  const [width, height] = match[1]
    .split("x")
    .map(s => (s === "" ? undefined : Number(s)));
  return (
    <img
      alt="Unknown"
      src={src.replace(imageSizeRegex, "")}
      width={width}
      height={height}
      {...props}
    />
  );
}

/**
 * @param {ReactMarkdown.ReactMarkdownProps} props
 */
const Markdown = ({ source, ...props }) => {
  /** @type {ReactMarkdown.ReactMarkdownProps['renderers']} */
  const renderers = {};

  source = imagePreprocessor(source);
  renderers["image"] = imageRenderer;

  return <ReactMarkdown source={source} renderers={renderers} {...props} />;
};

export default Markdown;
