/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./pkg/web/**/*.html", "./pkg/web/**/*.templ"],
  theme: {
    extend: {},
  },
  daisyui: {
    themes: false,
  },
  plugins: [
      require('daisyui'),
  ],
};
