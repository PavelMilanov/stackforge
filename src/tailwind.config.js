/** @type {import('tailwindcss').Config} */
export default {
  content: ['./views/**/*.templ', './views/**/*.go', './handlers/**/*.go'],
  theme: {
    extend: {
      colors: {
        stackforge: {
          bg: '#f6f7f8',
          sidebar: '#2e3238',
          green: '#609926',
          greenDark: '#4f7f20',
          greenSoft: '#eef7e8',
          border: '#d4d7dc',
          text: '#24292f',
          muted: '#6b7280'
        }
      }
    }
  }
};
