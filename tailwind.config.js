/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./views/**/*.{gotmpl,html}'],
  theme: {
    extend: {},
  },
  plugins: [require('daisyui')],
};
