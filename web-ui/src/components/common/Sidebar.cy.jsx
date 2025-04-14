import React from 'react';
import Sidebar from './Sidebar';
import { BrowserRouter as Router } from 'react-router-dom'; // Assuming you're using react-router

describe('<Sidebar />', () => {
  it('renders', () => {
    // see: https://on.cypress.io/mounting-react
    cy.mount(
      <Router>
        <Sidebar />
      </Router>
    );
  });

  it('collapses the sidebar on click', () => {
    cy.mount(
      <Router>
        <Sidebar />
      </Router>
    );

    // Select the collapse button using the data-cy attribute
    cy.get('[data-cy=collapse-button]').as('collapseButton');

    // Verify initial state (expanded)
    cy.get('@collapseButton').should('exist');
    cy.get('[data-cy=sidebar]').should('have.class', 'w-64');

    // Click to collapse
    cy.get('@collapseButton').click();
    cy.get('[data-cy=sidebar]').should('have.class', 'w-20');

    // Click to expand again
    cy.get('@collapseButton').click();
    cy.get('[data-cy=sidebar]').should('have.class', 'w-64');
  });
});