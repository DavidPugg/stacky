/** @type {import('tailwindcss').Config} */

const colors = {
  primary: "#2698BB",
  secondary: "#90BE6D",
  accent: "#f471b5",
  neutral: "#111111",
  "base-100": "#ECECEC",
  "base-200": "#DCDCDC",
  "base-300": "#BDBDBD",
  "base-400": "#ACACAC",
  "base-500": "#9D9D9D",
  info: "#0ca6e9",
  success: "#2bd4bd",
  warning: "#f4c152",
  error: "#fb6f84",
};

const sizes = {
  xs: "500px",
};

module.exports = {
  content: ["./views/**/*.html"],
  theme: {
    extend: {
      screens: {
        ...sizes,
      },
      width: {
        ...sizes,
      },
      maxWidth: {
        ...sizes,
      },
      minWidth: {
        ...sizes,
      },
      colors: {
        ...colors,
      },
      transitionProperty: {
        bg: "background-color",
      },
      animation: {
        "show-alert": "show 4s ease-in-out",
      },
      keyframes: {
        show: {
          "0%": { opacity: 0 },
          "5%": { opacity: 1 },
          "95%": { opacity: 1 },
          "100%": { opacity: 0 },
        },
      },
    },
  },

  plugins: [require("daisyui")],

  daisyui: {
    themes: [
      {
        mytheme: {
          ...colors,
        },
      },
    ],
  },
};
