import React from 'react';
import Header from './Header';

describe('<Header />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(<Header title="Test Header" />);
  });

  it('toggles dark mode on button click', () => {
    cy.mount(<Header title="Test Header" />);

    // Verify initial theme (dark mode by default)
    cy.document().then((doc) => {
      expect(doc.documentElement.classList.contains('dark')).to.be.true;
    });

    // Click the toggle button
    cy.get('button').click();

    // Verify theme switched to light mode
    cy.document().then((doc) => {
      expect(doc.documentElement.classList.contains('light')).to.be.true;
    });

    // Click the toggle button again
    cy.get('button').click();

    // Verify theme switched back to dark mode
    cy.document().then((doc) => {
      expect(doc.documentElement.classList.contains('dark')).to.be.true;
    });
  });
});