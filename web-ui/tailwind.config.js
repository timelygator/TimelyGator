/** @type {import('tailwindcss').Config} */
export default {
	content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
	theme: {
		extend: {
			colors:{
				background: "var(--color-primary-bg)",
				backgroundSecondary: "var(--color-secondary-bg)",
				text: "var(--color-primary-text)",
				border: "var(--color-border)",
				cardSecondaryText: "var(--color-card-secondary-text)",
			},
		},
	},
	plugins: [],
	darkMode: 'class',
};
