/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./views/**/*.{gotmpl,html}'],
  theme: {
    extend: {
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
};
