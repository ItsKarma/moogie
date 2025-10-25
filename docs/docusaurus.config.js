// @ts-check
// `@type` JSDoc annotations allow editor autocompletion and type checking
// (when paired with `@ts-check`).
// There are various equivalent ways to declare your Docusaurus config.
// See: https://docusaurus.io/docs/api/docusaurus-config

import { themes as prismThemes } from "prism-react-renderer";

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "Moogie Documentation",
  tagline: "Infrastructure monitoring made simple",
  favicon: "img/favicon.ico",

  // Set the production url of your site here
  url: "https://docs.moogie.dev",
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: "/",

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: "ItsKarma",
  projectName: "moogie",

  onBrokenLinks: "throw",

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  markdown: {
    mermaid: true,
    hooks: {
      onBrokenMarkdownLinks: "throw",
    },
  },

  themes: ["@docusaurus/theme-mermaid"],

  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: "./sidebars.js",
          routeBasePath: "/", // Serve docs at the root
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl: "https://github.com/ItsKarma/moogie/tree/main/docs/",
        },
        blog: false, // Disable blog
        theme: {
          customCss: "./src/css/custom.css",
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: "img/moogie-social-card.jpg",
      navbar: {
        title: "Moogie",
        logo: {
          alt: "Moogie Logo",
          src: "img/logo.svg",
        },
        items: [
          {
            type: "docSidebar",
            sidebarId: "tutorialSidebar",
            position: "left",
            label: "Documentation",
          },
          {
            href: "http://localhost:8080/api/v1",
            label: "API",
            position: "right",
          },
          {
            href: "http://localhost:3000",
            label: "Dashboard",
            position: "right",
          },
          {
            href: "https://github.com/ItsKarma/moogie",
            label: "GitHub",
            position: "right",
          },
        ],
      },
      footer: {
        style: "dark",
        links: [
          {
            title: "Documentation",
            items: [
              {
                label: "Getting Started",
                to: "/getting-started/docker-setup",
              },
              {
                label: "API Reference",
                to: "/api/overview",
              },
            ],
          },
          {
            title: "Community",
            items: [
              {
                label: "GitHub Issues",
                href: "https://github.com/ItsKarma/moogie/issues",
              },
              {
                label: "GitHub Discussions",
                href: "https://github.com/ItsKarma/moogie/discussions",
              },
            ],
          },
          {
            title: "More",
            items: [
              {
                label: "GitHub",
                href: "https://github.com/ItsKarma/moogie",
              },
              {
                label: "Dashboard",
                href: "http://localhost:3000",
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Moogie. Built with Docusaurus.`,
      },
      prism: {
        theme: prismThemes.github,
        darkTheme: prismThemes.dracula,
        additionalLanguages: ["bash", "yaml", "json", "go"],
      },
    }),
};

export default config;
