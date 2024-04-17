
let mainColor = '#FF6363';
let black = '#1f1f1f';


/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./html/**/*.html", "./internal/**/*.go", "./static/js/**/*.js"],
  theme: {
    extend: {
      colors: {
        'primary': mainColor,
        'black': black,
      },
      fontFamily: {
        'body': ['Nunito']
      },
      screens: {
        'xs': '480px',   // Extra small devices (phones)
        'sm': '640px',   // Small devices (tablets)
        'md': '768px',   // Medium devices (small laptops)
        'lg': '1024px',  // Large devices (laptops/desktops)
        'xl': '1280px',  // Extra large devices (large desktops)
        '2xl': '1536px'  // Bigger than extra large devices
      }
    },
  },
  plugins: [],
}

