
let main = '#B8860B';
let darkmain = '#8B4513';
let black = '#1f1f1f';
let faintgray = '#f9f9f9';


/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'selector',
  content: ["./html/**/*.html", "./internal/**/*.go", "./static/js/**/*.js"],
  theme: {
    extend: {
      colors: {
        'main': main,
        'darkmain': darkmain,
        'black': black,
        'faintgray': faintgray,
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

