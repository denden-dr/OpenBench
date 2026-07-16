/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./ui/views/**/*.templ",
    "./ui/html/**/*.html",
  ],
  theme: {
    extend: {
      colors: {
        primary: '#324376',    // Dark Blue
        secondary: '#586ba4',  // Muted Blue
        accent: '#f5dd90',     // Soft Yellow
        tertiary: '#f68e5f',   // Coral Orange (Assumed from f68e5)
        danger: '#f76c5e',     // Coral Red
      },
      fontFamily: {
        sans: ['"Plus Jakarta Sans"', 'sans-serif'],
        display: ['"Outfit"', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'monospace'],
      }
    }
  },
  plugins: [],
}
