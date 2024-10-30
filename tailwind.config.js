/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/templates/**/*.tmpl.html"],
  theme: {
    extend: {
      placeholderOpacity: 0.2,
    },
    screens: {
      'xs': '480px',
      'sm': '640px',
      'md': '768px',
      'lg': '1024px',
      'xl': '1280px',
      '2xl': '1536px',
    }
  },
  plugins: [
    require('daisyui'),
  ]
}
