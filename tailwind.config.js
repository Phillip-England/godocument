
let mainColor = '#FF6363';


/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./html/**/*.html", "./internal/**/*.go", "./static/js/**/*.js"],
  theme: {
    extend: {
      colors: {
        'primary': mainColor,
      },
      fontFamily: {
        'body': ['Nunito']
      }
    },
  },
  plugins: [],
}

