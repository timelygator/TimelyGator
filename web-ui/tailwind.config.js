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
				cardTextHeading: "var(--color-card-text-heading)",
				cardText: "var(--color-card-text)",
				actionDefault: "var(--color-action-default)",
				actionHover: "var(--color-action-hover)",
				toggleSelected: "var(--color-action-toggle-selected)",
				toggleUnselected: "var(--color-action-toggle-unselected)",
				buttonDefaultAccent1: "var(--color-button-default-accent1)",
				buttonHoverAccent1: "var(--color-button-hover-accent1)"
				
			},
		},
	},
	plugins: [],
	darkMode: 'class',
};
