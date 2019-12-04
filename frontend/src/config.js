export const app = {
  title: "Go Store"
};

// Any images referenced here should be placed in the public folder or be an absolute URL (http or https).

export const shop = {
  hero: {
    rows: [
      {
        columns: [
          {
            content: `
# Shop Products
`,
            align: "center"
          }
        ]
      }
    ]
  },
  heroImage: "/images/hero-3.jpg"
};

export const product = {
  heroImage: "/images/hero-3.jpg"
};

export const checkout = {
  heroImage: "/images/hero-2.jpg"
};

export const profile = {
  heroImage: "/images/hero-2.jpg"
};

export const home = {
  hero: {
    rows: [
      {
        columns: [
          {
            content: `
# A Shop Built With GO, GraphQL and React
`,
            callToAction: {
              label: "Explore the shop",
              to: "/shop",
              size: "lg"
            }
          }
        ]
      }
    ]
  },
  heroImage: "/images/hero-1.jpg",
  sections: [
    {
      rows: [
        {
          columns: [
            {
              content: `
## Configurable Page Sections

Each of these sections, as well as their background images are configurable via JSON and Markdown.

It can lead to some powerful designs as Sections have the following properties:

\`\`\`json
{
  "rows": [{
    "columns": [{
      "content": "# A Markdown String :D",
      "align": "An text-align value (center|left|right|...)",
      "callToAction": {
        "label": "Explore the shop",
        "to": "A path to navigate to (/shop, /about, ...)",
        "size": "A size for the button (sm|lg)"
        "align": "A justify-content option (center|flex-start|...)"
      },
      "breakpoints": {
        "xs": 12, // numbers between 1 and 12
        "sm": 6,
        "md": 4,
        "lg": 3,
        "xl": 2
      },
    }]
  }],
  "spacerImage": "An image path to use after the section that has parallax scrolling."
}
\`\`\`
`
            }
          ]
        }
      ],
      spacerImage: "/images/hero-2.jpg"
    },
    {
      rows: [
        {
          columns: [
            {
              align: "center",
              breakpoints: {
                xs: 12,
                sm: 12,
                md: 4
              },
              content: `
![Markdown Icon](images/markdown.png)

### Markdown

Markdown makes it easy for anyone to create awesome content.
`,
              callToAction: {
                label: "Explore the shop",
                to: "/shop",
                align: "center"
              }
            },
            {
              align: "center",
              breakpoints: {
                xs: 12,
                sm: 12,
                md: 4
              },
              content: `
![Async Icon](images/async.svg =x64)

### API Configurable

If you don't want to have to make a code change, you have all the freedom in the world to load
your site config from any source you wish.
`,
              callToAction: {
                label: "Explore the shop",
                to: "/shop",
                align: "center"
              }
            },
            {
              align: "center",
              breakpoints: {
                xs: 12,
                sm: 12,
                md: 4
              },
              content: `
![GraphCMS Icon](images/graphcms.png =x64)

### CMS

Easily plug into a CMS such as GraphCMS to make your experience updating the frontend even better.
`,
              callToAction: {
                label: "Explore the shop",
                to: "/shop",
                align: "center"
              }
            }
          ]
        },
        {
          columns: [
            {
              content: `
#### Configuration Example

The above row looks like:

\`\`\`
{
  "columns": [
    {
      "align": "center",
      "breakpoints": {
        "xs": 12,
        "sm": 12,
        "md": 4
      },
      "content": "![Markdown Icon](images/markdown.png)\\n\\n### Markdown\\n\\nMarkdown makes it easy for anyone to create awesome content.",
      "callToAction": {
        "label": "Explore the shop",
        "to": "/shop",
        "align": "center"
      }
    },
    {
      "align": "center",
      "breakpoints": {
        "xs": 12,
        "sm": 12,
        "md": 4
      },
      "content": "![Async Icon](images/async.svg =x64)\\n\\n### API Configurable\\n\\nIf you don't want to have to make a code change, you have all the freedom in the world to load your site config from any source you wish.",
      "callToAction": {
        "label": "Explore the shop",
        "to": "/shop",
        "align": "center"
      }
    },
    {
      "align": "center",
      "breakpoints": {
        "xs": 12,
        "sm": 12,
        "md": 4
      },
      "content": "![GraphCMS Icon](images/graphcms.png =x64)\\n\\n### CMS\\n\\nEasily plug into a CMS such as GraphCMS to make your experience updating the frontend even better.",
      "callToAction": {
        "label": "Explore the shop",
        "to": "/shop",
        "align": "center"
      }
    }
  ]
}
\`\`\`
`
            }
          ]
        }
      ],
      spacerImage: "/images/hero-3.jpg"
    }
  ]
};
