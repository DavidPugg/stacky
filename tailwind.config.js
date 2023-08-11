/** @type {import('tailwindcss').Config} */

const colors = {
  "primary": "#EA9010",
  "secondary": "#90BE6D",         
  "accent": "#f471b5",
  "neutral": "#FDFDFD",
  "base-100": "#17181C",
  "base-200": "#1E1F24",
  "base-300": "#27292F",
  "base-400": "#33353D",
  "base-500": "#454853",
  "info": "#0ca6e9",
  "success": "#2bd4bd",
  "warning": "#f4c152",
  "error": "#fb6f84"
}

module.exports = {
  content: ['./views/**/*.{gotmpl,html}'],
  theme: {
    extend: {
      colors :{
        ...colors
      },
      animation: {
        'show-alert': 'show 5s ease-in-out',
      },
      keyframes: {
        show: {
          '0%': { opacity: 0 },
          '5%': { opacity: 1 },
          '95%': { opacity: 1 },
          '100%': { opacity: 0 },
        },
      },
    },
  },

  plugins: [require('daisyui')],

  daisyui: {
    themes: [
      {
      mytheme: {
        ...colors
      },
    },
  ],
  },
};
