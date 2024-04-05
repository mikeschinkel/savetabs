/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './html/*.html',
    './html/scripts/*.js'
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require("daisyui")
  ],
}
