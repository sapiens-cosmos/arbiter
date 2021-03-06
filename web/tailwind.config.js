module.exports = {
  purge: ["./pages/**/*.{js,ts,jsx,tsx}", "./components/**/*.{js,ts,jsx,tsx}"],
  darkMode: false, // or 'media' or 'class'
  theme: {
    colors: {
      primary: "#DBE7DD",
      secondary: "#F3FAF1",
      border: "#BCE3B7",
      white: "#ffffff",
      black: "rgba(0, 0, 0, 0.85)",
      gray: "#A7BBC7",
    },
    backgroundColor: (themes) => ({
      ...themes("colors"),
    }),
    extend: {
      width: {
        button: "140px",
      },
      maxWidth: {
        default: "768px",
        modal: "550px",
      },
      maxHeight: {
        header: "80px",
      },
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
};
