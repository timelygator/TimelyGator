import React from 'react';
import Security from './Security';
import { BrowserRouter as Router } from 'react-router-dom'; // Assuming you're using react-router

describe('<Security />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(
      <Router>
        <Security />
      </Router>
    );
  });

  it('clicks the Change Password button', () => {
    cy.mount(
      <Router>
        <Security />
      </Router>
    );

    // Select the Change Password button using the data-cy attribute
    cy.get('[data-cy=change-password-button]').as('changePasswordButton');

    // Verify the button exists and is clickable
    cy.get('@changePasswordButton').should('exist').and('be.visible').click();
  });

  it('toggles the Two-Factor Authentication switch', () => {
    cy.mount(
      <Router>
        <Security />
      </Router>
    );

    // Select the Two-Factor Authentication toggle switch using the data-cy attribute
    cy.get('[data-cy=two-factor-toggle]').as('twoFactorToggle');

    // Verify initial state (off)
    cy.get('@twoFactorToggle').should('have.class', 'bg-gray-600');

    // Click to toggle on
    cy.get('@twoFactorToggle').click();
    cy.get('@twoFactorToggle').should('have.class', 'bg-indigo-600');

    // Click to toggle off
    cy.get('@twoFactorToggle').click();
    cy.get('@twoFactorToggle').should('have.class', 'bg-gray-600');
  });
});